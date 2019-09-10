package marketplace

import (
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

const (
	ModuleName                 = types.ModuleName
	RouterKey                  = types.RouterKey
	StoreKey                   = types.StoreKey
	RegisterCurrencyKey        = types.RegisterCurrency
	AuctionKey                 = types.AuctionKey
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
	MsgPutNFTOnMarket      = types.MsgPutNFTOnMarket
	MsgRemoveNFTFromMarket = types.MsgRemoveNFTFromMarket
	MsgBuyNFT              = types.MsgBuyNFT
	NFT                    = types.NFT
	NFTInfo                = types.NFTInfo
	MsgUpdateNFTParams     = types.MsgUpdateNFTParams
	FungibleToken          = types.FungibleToken

	MsgPutNFTOnAuction      = types.MsgPutNFTOnAuction
	MsgRemoveNFTFromAuction = types.MsgRemoveNFTFromAuction
	MsgMakeBidOnAuction     = types.MsgMakeBidOnAuction
	MsgFinishAuction        = types.MsgFinishAuction
	MsgBuyoutOnAuction      = types.MsgBuyoutOnAuction

	MsgCreateFungibleToken    = types.MsgCreateFungibleToken
	MsgTransferFungibleTokens = types.MsgTransferFungibleTokens
	MsgBurnFungibleToken      = types.MsgBurnFungibleTokens
	MsgMakeOffer              = types.MsgMakeOffer
	MsgAcceptOffer            = types.MsgAcceptOffer
)
