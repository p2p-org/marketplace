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
	cdc.RegisterConcrete(MsgMintNFT{}, "marketplace/MintNFT", nil)
	cdc.RegisterConcrete(MsgTransferNFT{}, "marketplace/TransferNFT", nil)
	cdc.RegisterConcrete(MsgPutNFTOnMarket{}, "marketplace/PutNFTOnMarket", nil)
	cdc.RegisterConcrete(MsgBuyNFT{}, "marketplace/BuyNFT", nil)
	cdc.RegisterConcrete(MsgCreateFungibleToken{}, "marketplace/CreateFungibleToken", nil)
	cdc.RegisterConcrete(MsgTransferFungibleTokens{}, "marketplace/TransferFungibleTokens", nil)
	cdc.RegisterConcrete(FungibleToken{}, "marketplace/FungibleToken", nil)
	cdc.RegisterConcrete(MsgUpdateNFTParams{}, "marketplace/UpdateNFTParams", nil)
}
