package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/corestario/marketplace/common"
	"github.com/corestario/marketplace/x/marketplace/types"
)

// NFTModuleMarketplace overrides the NFT module for custom handlers
type NFTModuleMarketplace struct {
	nft.AppModule
	nftKeeper *nft.Keeper
	mpKeeper  *Keeper
}

// NewHandler module handler for the NFTModuleMarketplace
func (m NFTModuleMarketplace) NewHandler() sdk.Handler {
	return CustomNFTHandler(m.nftKeeper, m.mpKeeper)
}

// NewNFTModuleMarketplace generates a new NFT Module
func NewNFTModuleMarketplace(appModule nft.AppModule, nftKeeper *nft.Keeper, mpKeeper *Keeper) *NFTModuleMarketplace {
	return &NFTModuleMarketplace{
		AppModule: appModule,
		nftKeeper: nftKeeper,
		mpKeeper:  mpKeeper,
	}
}

// CustomNFTHandler routes the messages to the handlers
func CustomNFTHandler(nftKeeper *nft.Keeper, mpKeeper *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case nft.MsgTransferNFT:
			return HandleMsgTransferNFTMarketplace(ctx, msg, nftKeeper, mpKeeper)
		case nft.MsgEditNFTMetadata:
			return nft.HandleMsgEditNFTMetadata(ctx, msg, *nftKeeper)
		case nft.MsgMintNFT:
			return HandleMsgMintNFTMarketplace(ctx, msg, nftKeeper, mpKeeper)
		case nft.MsgBurnNFT:
			return HandleMsgBurnNFTMarketplace(ctx, msg, nftKeeper, mpKeeper)
		default:
			errMsg := fmt.Sprintf("unrecognized nft message type: %T", msg)
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgTransferNFTMarketplace(ctx sdk.Context, msg nft.MsgTransferNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) sdk.Result {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgTransferNFT)
	res := nft.HandleMsgTransferNFT(ctx, msg, *nftKeeper)
	if !res.IsOK() {
		return res
	}

	// Create an account for recipient.
	if acc := mpKeeper.accKeeper.GetAccount(ctx, msg.Recipient); acc == nil {
		mpKeeper.accKeeper.SetAccount(ctx, mpKeeper.accKeeper.NewAccountWithAddress(ctx, msg.Recipient))
	}

	err := mpKeeper.TransferNFT(ctx, msg.ID, msg.Sender, msg.Recipient)
	if err != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to TransferNFT: %v", err)).Result()
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgTransferNFT)
	return res
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgMintNFTMarketplace(ctx sdk.Context, msg nft.MsgMintNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) sdk.Result {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMintNFT)

	// Create an account for the recipient of the minted NFTs.
	if acc := mpKeeper.accKeeper.GetAccount(ctx, msg.Recipient); acc == nil {
		mpKeeper.accKeeper.SetAccount(ctx, mpKeeper.accKeeper.NewAccountWithAddress(ctx, msg.Recipient))
	}

	res := nft.HandleMsgMintNFT(ctx, msg, *nftKeeper)
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

func HandleMsgBurnNFTMarketplace(ctx sdk.Context, msg nft.MsgBurnNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) sdk.Result {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBurnNFT)
	token, err := mpKeeper.GetNFT(ctx, msg.ID)
	if err != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to BurnNFT: no token with ID %s", msg.ID)).Result()
	}
	for _, offer := range token.Offers {
		if _, err := mpKeeper.coinKeeper.AddCoins(ctx, offer.Buyer, offer.Price); err != nil {
			return wrapError("failed to BurnNFT", err)
		}
	}
	res := nft.HandleMsgBurnNFT(ctx, msg, *nftKeeper)
	if !res.IsOK() {
		return res
	}
	if err := mpKeeper.BurnNFT(ctx, msg.ID); err != nil {
		return sdk.ErrUnknownRequest(fmt.Sprintf("failed to BurnNFT: %v", err)).Result()
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgBurnNFT)
	return sdk.Result{}
}
