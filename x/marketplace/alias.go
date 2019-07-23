package marketplace

import (
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

const (
	ModuleName = types.ModuleName
	RouterKey  = types.RouterKey
	StoreKey   = types.StoreKey

	FungibleTokenCreationPrice = types.FungibleTokenCreationPrice
	FungibleCommissionAddress  = types.FungibleCommissionAddress

	MaxBeneficiaryCommission = types.FlagMaxCommission
)

var (
	NewNFT        = types.NewNFT
	ModuleCdc     = types.ModuleCdc
	RegisterCodec = types.RegisterCodec
)

type (
	MsgMintNFT         = types.MsgMintNFT
	MsgTransferNFT     = types.MsgTransferNFT
	MsgPutNFTOnMarket  = types.MsgPutNFTOnMarket
	MsgBuyNFT          = types.MsgBuyNFT
	NFT                = types.NFT
	MsgUpdateNFTParams = types.MsgUpdateNFTParams
	FungibleToken      = types.FungibleToken

	MsgCreateFungibleToken    = types.MsgCreateFungibleToken
	MsgTransferFungibleTokens = types.MsgTransferFungibleTokens
)
