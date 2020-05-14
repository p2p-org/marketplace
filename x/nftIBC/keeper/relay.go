package keeper

import (
	"github.com/corestario/marketplace/x/marketplace"
	"github.com/cosmos/modules/incubator/nft"
	"strings"

	"github.com/corestario/marketplace/x/nftIBC/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	channeltypes "github.com/cosmos/cosmos-sdk/x/ibc/04-channel/types"
	ibctypes "github.com/cosmos/cosmos-sdk/x/ibc/types"
)

// SendTransfer handles transfer sending logic. There are 2 possible cases:
//
// 1. Sender chain is the source chain of the coins (i.e where they were minted): the coins
// are transferred to an escrow address (i.e locked) on the sender chain and then
// transferred to the destination chain (i.e not the source chain) via a packet
// with the corresponding fungible token data.
//
// 2. Coins are not native from the sender chain (i.e tokens sent where transferred over
// through IBC already): the coins are burned and then a packet is sent to the
// source chain of the tokens.
func (k Keeper) SendTransfer(
	ctx sdk.Context,
	sourcePort,
	sourceChannel string,
	destHeight uint64,
	id string,
	denom string,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
) error {
	sourceChannelEnd, found := k.channelKeeper.GetChannel(ctx, sourcePort, sourceChannel)
	if !found {
		return sdkerrors.Wrap(channeltypes.ErrChannelNotFound, sourceChannel)
	}

	destinationPort := sourceChannelEnd.GetCounterparty().GetPortID()
	destinationChannel := sourceChannelEnd.GetCounterparty().GetChannelID()

	// get the next sequence
	sequence, found := k.channelKeeper.GetNextSequenceSend(ctx, sourcePort, sourceChannel)
	if !found {
		return channeltypes.ErrSequenceSendNotFound
	}

	return k.createOutgoingPacket(ctx, sequence, sourcePort, sourceChannel, destinationPort, destinationChannel, destHeight, id, denom, sender, receiver)
}

// See spec for this function: https://github.com/cosmos/ics/tree/master/spec/ics-020-fungible-token-transfer#packet-relay
func (k Keeper) createOutgoingPacket(
	ctx sdk.Context,
	seq uint64,
	sourcePort, sourceChannel,
	destinationPort, destinationChannel string,
	destHeight uint64,
	id string,
	denom string,
	sender sdk.AccAddress,
	receiver sdk.AccAddress,
) error {
	channelCap, ok := k.scopedKeeper.GetCapability(ctx, ibctypes.ChannelCapabilityPath(sourcePort, sourceChannel))
	if !ok {
		return sdkerrors.Wrap(channeltypes.ErrChannelCapabilityNotFound, "module does not own channel capability")
	}
	// NOTE:
	// - Coins transferred from the destination chain should have their denomination
	// prefixed with source port and channel IDs.
	// - Coins transferred from the source chain can have their denomination
	// clear from prefixes when transferred to the escrow account (i.e when they are
	// locked) BUT MUST have the destination port and channel ID when constructing
	// the packet data.

	prefix := types.GetDenomPrefix(destinationPort, destinationChannel)
	source := strings.HasPrefix(denom, prefix)

	tokenURI := ""

	if source {

		tempDenom := denom

		if strings.HasPrefix(denom, prefix) {
			tempDenom = denom[len(prefix):]
		}

		token, err := k.nftKeeper.GetNFT(ctx, tempDenom, id)
		if err != nil {
			return err
		}
		tokenURI = token.GetTokenURI()

		// escrow tokens if the destination chain is the same as the sender's
		escrowAddress := types.GetEscrowAddress(sourcePort, sourceChannel)

		msgTransferNFT := nft.NewMsgTransferNFT(sender, escrowAddress, tempDenom, id)

		if _, err = marketplace.HandleMsgTransferNFTMarketplace(ctx, msgTransferNFT, k.nftKeeper, k.mpKeeper); err != nil {
			return err
		}

	} else {
		// build the receiving denomination prefix if it's not present
		prefix = types.GetDenomPrefix(sourcePort, sourceChannel)

		if !strings.HasPrefix(denom, prefix) {
			return sdkerrors.Wrapf(types.ErrInvalidDenomForTransfer, "denom was: %s", denom)
		}

		token, err := k.nftKeeper.GetNFT(ctx, denom, id)
		if err != nil {
			return err
		}
		tokenURI = token.GetTokenURI()

		msgBurnNFT := nft.NewMsgBurnNFT(sender, id, denom)

		if _, err := marketplace.HandleMsgBurnNFTMarketplace(ctx, msgBurnNFT, k.nftKeeper, k.mpKeeper); err != nil {
			return err
		}
	}

	packetData := types.NewNFTPacketData(
		id, denom, sender, receiver, tokenURI,
	)

	packet := channeltypes.NewPacket(
		packetData.GetBytes(),
		seq,
		sourcePort,
		sourceChannel,
		destinationPort,
		destinationChannel,
		destHeight+DefaultPacketTimeoutHeight,
		DefaultPacketTimeoutTimestamp,
	)

	return k.channelKeeper.SendPacket(ctx, channelCap, packet)
}

func (k Keeper) OnRecvPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NFTPacketData) error {
	// NOTE: packet data type already checked in handler.go

	prefix := types.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	source := strings.HasPrefix(data.Denom, prefix)

	if source {
		mintNFTMsg := nft.NewMsgMintNFT(data.Owner, data.Receiver, data.Id, data.Denom, data.TokenURI)
		if _, err := marketplace.HandleMsgMintNFTMarketplace(ctx, mintNFTMsg, k.nftKeeper, k.mpKeeper); err != nil {
			return err
		}
		return nil
	}

	// check the denom prefix
	prefix = types.GetDenomPrefix(packet.GetSourcePort(), packet.GetSourceChannel())

	if !strings.HasPrefix(data.Denom, prefix) {
		return sdkerrors.Wrapf(
			sdkerrors.ErrInvalidCoins,
			"%s doesn't contain the prefix '%s'", data.Denom, prefix,
		)
	}

	data.Denom = data.Denom[len(prefix):]

	// unescrow tokens
	escrowAddress := types.GetEscrowAddress(packet.GetDestPort(), packet.GetDestChannel())

	msgTransferNFT := nft.NewMsgTransferNFT(escrowAddress, data.Receiver, data.Denom, data.Id)
	_, err := marketplace.HandleMsgTransferNFTMarketplace(ctx, msgTransferNFT, k.nftKeeper, k.mpKeeper)
	return err
}

func (k Keeper) OnAcknowledgementPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NFTPacketData, ack types.NFTPacketAcknowledgement) error {
	if !ack.Success {
		return k.refundPacketAmount(ctx, packet, data)
	}
	return nil
}

func (k Keeper) OnTimeoutPacket(ctx sdk.Context, packet channeltypes.Packet, data types.NFTPacketData) error {
	return k.refundPacketAmount(ctx, packet, data)
}

func (k Keeper) refundPacketAmount(ctx sdk.Context, packet channeltypes.Packet, data types.NFTPacketData) error {
	// NOTE: packet data type already checked in handler.go

	// check the denom prefix
	prefix := types.GetDenomPrefix(packet.GetDestPort(), packet.GetDestChannel())
	source := strings.HasPrefix(data.Denom, prefix)

	if source {

		if !strings.HasPrefix(data.Denom, prefix) {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidCoins, "%s doesn't contain the prefix '%s'", data.Denom, prefix)
		}
		data.Denom = data.Denom[len(prefix):]

		// unescrow tokens back to sender
		escrowAddress := types.GetEscrowAddress(packet.GetSourcePort(), packet.GetSourceChannel())

		msgTransferNFT := nft.NewMsgTransferNFT(escrowAddress, data.Owner, data.Denom, data.Id)
		_, err := marketplace.HandleMsgTransferNFTMarketplace(ctx, msgTransferNFT, k.nftKeeper, k.mpKeeper)
		return err
	}

	// mint vouchers back to sender
	mintNFTMsg := nft.NewMsgMintNFT(data.Owner, data.Owner, data.Id, data.Id, data.TokenURI)
	if _, err := marketplace.HandleMsgMintNFTMarketplace(ctx, mintNFTMsg, k.nftKeeper, k.mpKeeper); err != nil {
		return err
	}
	return nil
}
