package marketplace

import (
	"encoding/json"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/ibc"
	client "github.com/cosmos/cosmos-sdk/x/ibc/02-client"
	connection "github.com/cosmos/cosmos-sdk/x/ibc/03-connection"
	channel "github.com/cosmos/cosmos-sdk/x/ibc/04-channel"
	transfer "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer"
	"github.com/corestario/marketplace/x/marketplace/types"
)

type IBCModuleMarketplace struct {
	ibc.AppModule
	ibcKeeper *ibc.Keeper
	mpKeeper  *Keeper
}

func (m IBCModuleMarketplace) NewHandler() sdk.Handler {
	return CustomIBCHandler(m.ibcKeeper, m.mpKeeper)
}

func NewIBCModuleMarketplace(appModule ibc.AppModule, ibcKeeper *ibc.Keeper, mpKeeper *Keeper) *IBCModuleMarketplace {
	return &IBCModuleMarketplace{
		AppModule: appModule,
		ibcKeeper: ibcKeeper,
		mpKeeper:  mpKeeper,
	}
}

func CustomIBCHandler(ibcKeeper *ibc.Keeper, mpKeeper *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		ctx = ctx.WithEventManager(sdk.NewEventManager())

		switch msg := msg.(type) {
		// IBC client msgs
		case client.MsgCreateClient:
			return client.HandleMsgCreateClient(ctx, ibcKeeper.ClientKeeper, msg)

		case client.MsgUpdateClient:
			return client.HandleMsgUpdateClient(ctx, ibcKeeper.ClientKeeper, msg)

		case client.MsgSubmitMisbehaviour:
			return client.HandleMsgSubmitMisbehaviour(ctx, ibcKeeper.ClientKeeper, msg)

		// IBC connection  msgs
		case connection.MsgConnectionOpenInit:
			return connection.HandleMsgConnectionOpenInit(ctx, ibcKeeper.ConnectionKeeper, msg)

		case connection.MsgConnectionOpenTry:
			return connection.HandleMsgConnectionOpenTry(ctx, ibcKeeper.ConnectionKeeper, msg)

		case connection.MsgConnectionOpenAck:
			return connection.HandleMsgConnectionOpenAck(ctx, ibcKeeper.ConnectionKeeper, msg)

		case connection.MsgConnectionOpenConfirm:
			return connection.HandleMsgConnectionOpenConfirm(ctx, ibcKeeper.ConnectionKeeper, msg)

		// IBC channel msgs
		case channel.MsgChannelOpenInit:
			return channel.HandleMsgChannelOpenInit(ctx, ibcKeeper.ChannelKeeper, msg)

		case channel.MsgChannelOpenTry:
			return channel.HandleMsgChannelOpenTry(ctx, ibcKeeper.ChannelKeeper, msg)

		case channel.MsgChannelOpenAck:
			return channel.HandleMsgChannelOpenAck(ctx, ibcKeeper.ChannelKeeper, msg)

		case channel.MsgChannelOpenConfirm:
			return channel.HandleMsgChannelOpenConfirm(ctx, ibcKeeper.ChannelKeeper, msg)

		case channel.MsgChannelCloseInit:
			return channel.HandleMsgChannelCloseInit(ctx, ibcKeeper.ChannelKeeper, msg)

		case channel.MsgChannelCloseConfirm:
			return channel.HandleMsgChannelCloseConfirm(ctx, ibcKeeper.ChannelKeeper, msg)

		case transfer.MsgTransfer:
			return transfer.HandleMsgTransfer(ctx, ibcKeeper.TransferKeeper, msg)

		case transfer.MsgRecvPacket:
			return HandleMsgRecvPacket(ctx, mpKeeper, ibcKeeper, msg)

		default:
			errMsg := fmt.Sprintf("unrecognized IBC message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func HandleMsgRecvPacket(ctx sdk.Context, mpKeeper *Keeper, k *ibc.Keeper, msg transfer.MsgRecvPacket) (res sdk.Result) {
	if _, err := k.ChannelKeeper.RecvPacket(ctx, msg.Packet, msg.Proofs[0], msg.Height, nil, sdk.NewKVStoreKey(ibc.StoreKey)); err != nil {
		return sdk.ResultFromError(err)
	}

	switch msg.Packet.GetDestPort() {
	case types.IBCNFTPort:
		var data types.NFTPacketData
		if err := json.Unmarshal(msg.Packet.GetData(), &data); err != nil {
			return sdk.ResultFromError(err)
		}
		if err := mpKeeper.ReceiveNFTByIBCTransferTx(ctx, data, msg.Packet); err != nil {
			return sdk.ResultFromError(err)
		}
	case "bank":
		var data transfer.PacketData
		err := data.UnmarshalJSON(msg.Packet.GetData())
		if err != nil {
			return sdk.ResultFromError(err)
		}

		return sdk.ResultFromError(k.TransferKeeper.ReceiveTransfer(ctx, msg.Packet.GetSourcePort(), msg.Packet.GetSourceChannel(), msg.Packet.GetDestPort(), msg.Packet.GetDestChannel(), data))
	}
	return sdk.Result{Events: ctx.EventManager().Events()}
}

func HandleMsgTransferNFTByIBC(ctx sdk.Context, k *Keeper, msg MsgTransferNFTByIBC) sdk.Result {
	nft, err := k.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return sdk.ResultFromError(err)
	}
	fullNFT, err := k.nftKeeper.GetNFT(ctx, nft.Denom, nft.ID)
	if err != nil {
		return sdk.ResultFromError(err)
	}
	err = k.SendNFTByIBCTransferTx(ctx, fullNFT.GetID(), nft.Denom, fullNFT.GetTokenURI(), msg.SourcePort, msg.SourceChannel, msg.Sender, msg.Receiver, msg.Source)
	if err != nil {
		return sdk.ResultFromError(err)
	}

	return sdk.Result{Events: ctx.EventManager().Events()}
}
