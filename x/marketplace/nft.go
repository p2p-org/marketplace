package marketplace

import (
	"fmt"

	sdk "github.com/dgamingfoundation/cosmos-sdk/types"
	"github.com/dgamingfoundation/cosmos-sdk/x/nft"
	"github.com/dgamingfoundation/marketplace/common"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

// NFTModuleMarketplace overrides the NFT module for custom handlers
type NFTModuleMarketplace struct {
	nft.AppModule
	nftKeeper nft.Keeper
	mpKeeper  *Keeper
}

// NewHandler module handler for the NFTModuleMarketplace
func (m NFTModuleMarketplace) NewHandler() sdk.Handler {
	return CustomNFTHandler(m.nftKeeper, m.mpKeeper)
}

// NewNFTModuleMarketplace generates a new NFT Module
func NewNFTModuleMarketplace(appModule nft.AppModule, nftKeeper nft.Keeper, mpKeeper *Keeper) NFTModuleMarketplace {
	return NFTModuleMarketplace{
		AppModule: appModule,
		nftKeeper: nftKeeper,
		mpKeeper:  mpKeeper,
	}
}

// CustomNFTHandler routes the messages to the handlers
func CustomNFTHandler(nftKeeper nft.Keeper, mpKeeper *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case nft.MsgTransferNFT:
			return HandleMsgTransferNFTMarketplace(ctx, msg, nftKeeper, mpKeeper)
		case nft.MsgEditNFTMetadata:
			return nft.HandleMsgEditNFTMetadata(ctx, msg, nftKeeper)
		case nft.MsgMintNFT:
			return HandleMsgMintNFTMarketplace(ctx, msg, nftKeeper, mpKeeper)
		case nft.MsgBurnNFT:
			return nft.HandleMsgBurnNFT(ctx, msg, nftKeeper)
		default:
			errMsg := fmt.Sprintf("unrecognized nft message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgTransferNFTMarketplace(ctx sdk.Context, msg nft.MsgTransferNFT, nftKeeper nft.Keeper, mpKeeper *Keeper) sdk.Result {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMintNFT)
	res := nft.HandleMsgTransferNFT(ctx, msg, nftKeeper)
	if !res.IsOK() {
		return res
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMintNFT)
	return res
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgMintNFTMarketplace(ctx sdk.Context, msg nft.MsgMintNFT, nftKeeper nft.Keeper, mpKeeper *Keeper) sdk.Result {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMintNFT)
	res := nft.HandleMsgMintNFT(ctx, msg, nftKeeper)
	if !res.IsOK() {
		return res
	}

	mpNFToken := NewNFT(msg.ID, msg.Denom, msg.Recipient, sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.NewInt(0))))
	if err := mpKeeper.MintNFT(ctx, mpNFToken); err != nil {
		sdk.ErrUnknownRequest(err.Error()).Result()
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMintNFT)
	return sdk.Result{}
}

func HandleMsgBurnNFTMarketplace(ctx sdk.Context, msg nft.MsgBurnNFT, nftKeeper nft.Keeper, mpKeeper *Keeper) sdk.Result {
	res := nft.HandleMsgBurnNFT(ctx, msg, nftKeeper)
	if !res.IsOK() {
		return res
	}
	// TODO @rybnov: implement BurnNFT for marketplace keeper.
	// if err := mpKeeper.BurnNFT(msg.ID, msg.Sender); err != nil {
	// 		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to BurnNFT: %v", err)).Result()
	// }
	return sdk.Result{}
}
