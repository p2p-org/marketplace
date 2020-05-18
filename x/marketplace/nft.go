package marketplace

import (
	"fmt"

	"github.com/p2p-org/marketplace/common"
	"github.com/p2p-org/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
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
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
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
			return nil, fmt.Errorf(errMsg)
		}
	}
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgTransferNFTMarketplace(ctx sdk.Context, msg nft.MsgTransferNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgTransferNFT)
	res, err := nft.HandleMsgTransferNFT(ctx, msg, *nftKeeper)
	if err != nil {
		return nil, err //TODO: return error
	}

	// Create an account for recipient.
	if acc := mpKeeper.accKeeper.GetAccount(ctx, msg.Recipient); acc == nil {
		mpKeeper.accKeeper.SetAccount(ctx, mpKeeper.accKeeper.NewAccountWithAddress(ctx, msg.Recipient))
	}

	err = mpKeeper.TransferNFT(ctx, msg.ID, msg.Sender, msg.Recipient)
	if err != nil {
		return nil, fmt.Errorf("failed to TransferNFT: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgTransferNFT)
	return res, nil
}

// HandleMsgMintNFTMarketplace handles MsgMintNFT
func HandleMsgMintNFTMarketplace(ctx sdk.Context, msg nft.MsgMintNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMintNFT)

	deletedStore := ctx.KVStore(mpKeeper.deletedStoreKey)
	if deletedStore.Has([]byte(msg.ID)) {
		return nil, fmt.Errorf("NFT #%s has been deleted", msg.ID)
	}

	// Create an account for the recipient of the minted NFTs.
	if acc := mpKeeper.accKeeper.GetAccount(ctx, msg.Recipient); acc == nil {
		mpKeeper.accKeeper.SetAccount(ctx, mpKeeper.accKeeper.NewAccountWithAddress(ctx, msg.Recipient))
	}

	res, err := nft.HandleMsgMintNFT(ctx, msg, *nftKeeper)
	if err != nil {
		return nil, err
	}

	mpNFToken := NewNFT(msg.ID, msg.Denom, msg.Recipient, sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.NewInt(0))))
	if err := mpKeeper.MintNFT(ctx, mpNFToken); err != nil {
		return nil, err
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMintNFT)
	return res, nil
}

func HandleMsgBurnNFTMarketplace(ctx sdk.Context, msg nft.MsgBurnNFT, nftKeeper *nft.Keeper, mpKeeper *Keeper) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBurnNFT)
	token, err := mpKeeper.GetNFT(ctx, msg.ID)
	if err != nil {
		return nil, fmt.Errorf("failed to BurnNFT: no token with ID %s", msg.ID)
	}
	for _, offer := range token.Offers {
		if _, err := mpKeeper.coinKeeper.AddCoins(ctx, offer.Buyer, offer.Price); err != nil {
			return nil, err
		}
	}
	res, err := nft.HandleMsgBurnNFT(ctx, msg, *nftKeeper)
	if err != nil {
		return nil, err
	}
	if err := mpKeeper.BurnNFT(ctx, msg.ID); err != nil {
		return nil, fmt.Errorf("failed to BurnNFT: %v", err)
	}

	deletedStore := ctx.KVStore(mpKeeper.deletedStoreKey)
	deletedStore.Set([]byte(msg.ID), []byte{})

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgBurnNFT)
	return res, nil
}
