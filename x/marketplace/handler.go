package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/google/uuid"
	abci_types "github.com/tendermint/tendermint/abci/types"
)

const (
	ValidatorsCommission = 0.01
)

// NewHandler returns a handler for "marketplace" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMintNFT:
			return handleMsgMintNFT(ctx, keeper, msg)
		case MsgTransferNFT:
			return handleMsgTransferNFT(ctx, keeper, msg)
		case MsgSellNFT:
			return handleMsgSellNFT(ctx, keeper, msg)
		case MsgBuyNFT:
			return handleMsgBuyNFT(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized marketplace Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgMintNFT(ctx sdk.Context, keeper Keeper, msg MsgMintNFT) sdk.Result {
	nft := NewNFT(
		xnft.NewBaseNFT(
			uuid.New().String(),
			msg.Owner,
			msg.Name,
			msg.Description,
			msg.Image,
			msg.TokenURI,
		),
		sdk.NewCoin("token", sdk.NewInt(0)),
	)
	if err := keeper.MintNFT(ctx, nft); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to MintNFT: %v", err)),
		}
	}

	return sdk.Result{}
}

func handleMsgTransferNFT(ctx sdk.Context, keeper Keeper, msg MsgTransferNFT) sdk.Result {
	if err := keeper.TransferNFT(ctx, msg.TokenID, msg.Sender, msg.Recipient); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to TransferNFT: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgSellNFT(ctx sdk.Context, keeper Keeper, msg MsgSellNFT) sdk.Result {
	if err := keeper.SellNFT(ctx, msg.TokenID, msg.Owner, msg.Price); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to SellNFT: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgBuyNFT(ctx sdk.Context, keeper Keeper, msg MsgBuyNFT) sdk.Result {
	nft, err := keeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: %v", err)),
		}
	}

	if !nft.IsOnSale() {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: token %s is not for sale", nft.GetID())),
		}
	}

	err = keeper.coinKeeper.SendCoins(ctx, msg.Buyer, nft.GetOwner(), sdk.NewCoins(nft.GetPrice()))
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	nft.SetOwner(msg.Buyer)
	nft.SetOnSale(false)

	if err := keeper.UpdateNFT(ctx, nft); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: %v", err)),
		}
	}

	return sdk.Result{}
}

func doCommissions(ctx sdk.Context, keeper Keeper, payer sdk.AccAddress, price sdk.Coin) error {
	votes := ctx.VoteInfos()
	var vals []abci_types.Validator
	for _, vote := range votes {
		if vote.SignedLastBlock {
			vals = append(vals, vote.Validator)
		}
	}

	priceFloat64 := float64(price.Amount.Int64())
	totalValRewardAmount := sdk.NewCoin(
		price.Denom,
		sdk.NewInt(int64(priceFloat64*ValidatorsCommission)),
	)

	if !keeper.coinKeeper.HasCoins(ctx, payer, sdk.NewCoins(totalValRewardAmount)) {
		return fmt.Errorf("user %s does not have enough funds", payer.String())
	}

	singleValRewardAmount := sdk.NewCoin(
		price.Denom,
		sdk.NewInt((int64(priceFloat64*ValidatorsCommission))/int64(len(vals))),
	)

	for idx, val := range vals {
		if err := keeper.coinKeeper.SendCoins(ctx, payer, val.Address, sdk.NewCoins(singleValRewardAmount)); err != nil {
			// TODO: rollback payments.
			fmt.Println("Rollback after: ", idx)
			return fmt.Errorf("failed to pay commission: %v", err)
		}
	}

	return nil
}
