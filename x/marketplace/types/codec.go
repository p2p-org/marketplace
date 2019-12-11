package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
)

var ModuleCdc = codec.New()

func init() {
	RegisterCodec(ModuleCdc)
}

// RegisterCodec registers concrete types on the Amino codec
func RegisterCodec(cdc *codec.Codec) {
	cdc.RegisterConcrete(NFT{}, "marketplace/NFT", nil)
	cdc.RegisterConcrete(MsgPutNFTOnMarket{}, "marketplace/PutNFTOnMarket", nil)
	cdc.RegisterConcrete(MsgBuyNFT{}, "marketplace/BuyNFT", nil)
	cdc.RegisterConcrete(MsgCreateFungibleToken{}, "marketplace/CreateFungibleToken", nil)
	cdc.RegisterConcrete(MsgTransferFungibleTokens{}, "marketplace/TransferFungibleTokens", nil)
	cdc.RegisterConcrete(FungibleToken{}, "marketplace/FungibleToken", nil)
	cdc.RegisterConcrete(AuctionLot{}, "marketplace/AuctionLot", nil)
	cdc.RegisterConcrete(AuctionBid{}, "marketplace/AuctionBid", nil)
	cdc.RegisterConcrete(MsgUpdateNFTParams{}, "marketplace/UpdateNFTParams", nil)
	cdc.RegisterConcrete(MsgRemoveNFTFromMarket{}, "marketplace/RemoveNFTFromMarket", nil)
	cdc.RegisterConcrete(MsgPutNFTOnAuction{}, "marketplace/MsgPutNFTOnAuction", nil)
	cdc.RegisterConcrete(MsgRemoveNFTFromAuction{}, "marketplace/MsgRemoveNFTFromAuction", nil)
	cdc.RegisterConcrete(MsgMakeBidOnAuction{}, "marketplace/MsgMakeBidOnAuction", nil)
	cdc.RegisterConcrete(MsgBuyoutOnAuction{}, "marketplace/MsgBuyoutOnAuction", nil)
	cdc.RegisterConcrete(MsgFinishAuction{}, "marketplace/MsgFinishAuction", nil)
	cdc.RegisterConcrete(MsgBurnFungibleTokens{}, "marketplace/BurnFungibleTokens", nil)
	cdc.RegisterConcrete(NFTMetaData{}, "marketplace/NFTMetaData", nil)
	cdc.RegisterConcrete(NFTInfo{}, "marketplace/NFTInfo", nil)
	cdc.RegisterConcrete(MsgBatchTransfer{}, "marketplace/MsgBatchTransfer", nil)
	cdc.RegisterConcrete(MsgBatchPutOnMarket{}, "marketplace/MsgBatchPutOnMarket", nil)
	cdc.RegisterConcrete(MsgBatchRemoveFromMarket{}, "marketplace/MsgBatchRemoveFromMarket", nil)
	cdc.RegisterConcrete(MsgBatchBuyOnMarket{}, "marketplace/MsgBatchBuyOnMarket", nil)
	cdc.RegisterConcrete(MsgMakeOffer{}, "marketplace/MakeOffer", nil)
	cdc.RegisterConcrete(MsgAcceptOffer{}, "marketplace/AcceptOffer", nil)
	cdc.RegisterConcrete(MsgRemoveOffer{}, "marketplace/RemoveOffer", nil)
	cdc.RegisterConcrete(MsgTransferNFTByIBC{}, "marketplace/MsgTransferNFT", nil)
}
