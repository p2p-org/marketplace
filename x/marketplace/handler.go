package marketplace

import (
	"fmt"
	"math/big"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/google/uuid"
	abci_types "github.com/tendermint/tendermint/abci/types"
)

var (
	ValidatorsCommission    = 0.01
	BeneficiariesCommission = 0.015
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
		sdk.NewCoins(sdk.NewCoin("token", sdk.NewInt(0))),
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

func handleMsgTransferNFT(ctx sdk.Context, k Keeper, msg MsgTransferNFT) sdk.Result {
	if err := k.TransferNFT(ctx, msg.TokenID, msg.Sender, msg.Recipient); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to TransferNFT: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgSellNFT(ctx sdk.Context, k Keeper, msg MsgSellNFT) sdk.Result {
	nft, err := k.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: %v", err)),
		}
	}

	if err := doCommissions(ctx, k, msg.Owner, msg.Beneficiary, nft.Price); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: failed to pay commissions: %v", err)),
		}
	}

	if err := k.SellNFT(ctx, msg.TokenID, msg.Owner, msg.Price); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to SellNFT: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgBuyNFT(ctx sdk.Context, k Keeper, msg MsgBuyNFT) sdk.Result {
	nft, err := k.GetNFT(ctx, msg.TokenID)
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

	if err := doCommissions(ctx, k, msg.Buyer, msg.Beneficiary, nft.Price); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: failed to pay commissions: %v", err)),
		}
	}

	err = k.coinKeeper.SendCoins(ctx, msg.Buyer, nft.GetOwner(), nft.GetPrice())
	if err != nil {
		return sdk.ErrInsufficientCoins("Buyer does not have enough coins").Result()
	}

	nft.SetOwner(msg.Buyer)
	nft.SetOnSale(false)

	if err := k.UpdateNFT(ctx, nft); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: %v", err)),
		}
	}

	return sdk.Result{}
}

func doCommissions(ctx sdk.Context, k Keeper, payer, beneficiary sdk.AccAddress, price sdk.Coins) error {
	// Check that payer has enough funds (for both the commission and the asset itself).
	totalCommission := getCommission(price, ValidatorsCommission+BeneficiariesCommission)
	if !k.coinKeeper.HasCoins(ctx, payer, totalCommission) {
		return fmt.Errorf("user %s does not have enough funds", payer.String())
	}

	// Pay commission to the beneficiary.
	beneficiaryCommission := getCommission(price, BeneficiariesCommission)
	if err := k.coinKeeper.SendCoins(ctx, payer, beneficiary, beneficiaryCommission); err != nil {
		return fmt.Errorf("failed to pay commission to beneficiary: %v", err)
	}

	votes := ctx.VoteInfos()
	var vals []abci_types.Validator
	for _, vote := range votes {
		if vote.SignedLastBlock {
			vals = append(vals, vote.Validator)
		}
	}

	singleValRewardAmount := getCommission(price, BeneficiariesCommission/float64(len(vals)))
	for valIdx, val := range vals {
		if err := k.coinKeeper.SendCoins(ctx, payer, val.Address, singleValRewardAmount); err != nil {
			fmt.Printf("Failed to pay commission to validator %s, rolling back transactions", val.Address)

			for rollbackIdx := 0; rollbackIdx < valIdx; rollbackIdx++ {
				if err := k.coinKeeper.SendCoins(ctx, val.Address, payer, singleValRewardAmount); err != nil {
					panic(fmt.Sprintf("failed to rollback commission to validator %s: %v", val.Address, err))
				}
			}

			return fmt.Errorf("failed to pay commission to validator %s: %v", val.Address, err)
		}
	}

	return nil
}

func getCommission(price sdk.Coins, rat64 float64) sdk.Coins {
	// TODO: maybe we can do it somehow easier.
	var rat = new(big.Rat)
	rat = rat.SetFloat64(rat64)
	num, denom := sdk.NewDecFromBigInt(rat.Num()), sdk.NewDecFromBigInt(rat.Denom())
	priceDec := sdk.NewDecCoins(price)
	totalCommission, _ := priceDec.MulDec(num).QuoDec(denom).TruncateDecimal()
	return totalCommission
}
