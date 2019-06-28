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
	NewMsgMintNFT = types.NewMsgMintNFT
	NewNFT        = types.NewNFT
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgMintNFT      = types.MsgMintNFT
	QueryResResolve = types.QueryResResolve
	QueryResNames   = types.QueryResNames
	NFT             = types.NFT
)
