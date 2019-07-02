package marketplace

import (
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey
)

var (
	NewNFT        = types.NewNFT
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgMintNFT     = types.MsgMintNFT
	MsgTransferNFT = types.MsgTransferNFT
	MsgSellNFT     = types.MsgSellNFT
	MsgBuyNFT      = types.MsgBuyNFT
	NFT            = types.NFT
)
