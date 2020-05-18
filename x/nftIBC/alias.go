package nftIBC

import (
	"github.com/p2p-org/marketplace/x/nftIBC/keeper"
	"github.com/p2p-org/marketplace/x/nftIBC/types"
	ibcTypes "github.com/cosmos/cosmos-sdk/x/ibc/20-transfer/types"
)

const (
	DefaultPacketTimeoutHeight    = keeper.DefaultPacketTimeoutHeight
	DefaultPacketTimeoutTimestamp = keeper.DefaultPacketTimeoutTimestamp
	EventTypeTimeout              = types.EventTypeTimeout
	EventTypePacket               = types.EventTypePacket
	EventTypeChannelClose         = types.EventTypeChannelClose
	AttributeKeyReceiver          = types.AttributeKeyReceiver
	AttributeKeyValue             = types.AttributeKeyValue
	AttributeKeyRefundReceiver    = types.AttributeKeyRefundReceiver
	AttributeKeyRefundValue       = types.AttributeKeyRefundValue
	AttributeKeyAckSuccess        = types.AttributeKeyAckSuccess
	AttributeKeyAckError          = types.AttributeKeyAckError
	ModuleName                    = types.ModuleName
	StoreKey                      = types.StoreKey
	RouterKey                     = types.RouterKey
	QuerierRoute                  = types.QuerierRoute
)

var (
	// functions aliases
	NewKeeper            = keeper.NewKeeper
	RegisterCodec        = types.RegisterCodec
	GetEscrowAddress     = types.GetEscrowAddress
	GetDenomPrefix       = types.GetDenomPrefix
	GetModuleAccountName = types.GetModuleAccountName
	NewMsgTransferNFT    = types.NewMsgTransferNFT
	RegisterInterfaces   = types.RegisterInterfaces

	// variable aliases
	ModuleCdc              = types.ModuleCdc
	AttributeValueCategory = types.AttributeValueCategory
)

type (
	Keeper                   = keeper.Keeper
	BankKeeper               = ibcTypes.BankKeeper
	ChannelKeeper            = ibcTypes.ChannelKeeper
	ClientKeeper             = ibcTypes.ClientKeeper
	ConnectionKeeper         = ibcTypes.ConnectionKeeper
	NFTPacketData            = types.NFTPacketData
	NFTPacketAcknowledgement = types.NFTPacketAcknowledgement
	MsgTransferNFT           = types.MsgTransferNFT
)
