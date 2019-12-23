package marketplace

import (
	"github.com/corestario/marketplace/x/marketplace/types"
)

const (
	ModuleName                 = types.ModuleName
	RouterKey                  = types.RouterKey
	StoreKey                   = types.StoreKey
	RegisterCurrencyKey        = types.RegisterCurrency
	AuctionKey                 = types.AuctionKey
	DeletedNFTKey              = types.DeletedNFTKey
	FungibleTokenCreationPrice = types.FungibleTokenCreationPrice
	FungibleCommissionAddress  = types.FungibleCommissionAddress

	MaxBeneficiaryCommission = types.FlagMaxCommission
)

var (
	NewNFT          = types.NewNFT
	ModuleCdc       = types.ModuleCdc
	RegisterCodec   = types.RegisterCodec
	EventKeyOfferID = types.AttributeKeyOfferID
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
	MsgBatchTransfer          = types.MsgBatchTransfer
	MsgBatchPutOnMarket       = types.MsgBatchPutOnMarket
	MsgBatchRemoveFromMarket  = types.MsgBatchRemoveFromMarket
	MsgBatchBuyOnMarket       = types.MsgBatchBuyOnMarket
	MsgMakeOffer              = types.MsgMakeOffer
	MsgAcceptOffer            = types.MsgAcceptOffer
	MsgRemoveOffer            = types.MsgRemoveOffer
	MsgTransferNFTByIBC       = types.MsgTransferNFTByIBC
)
