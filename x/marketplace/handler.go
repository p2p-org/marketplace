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
	// TODO: make commissions configurable.
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
	if err := k.SellNFT(ctx, msg.TokenID, msg.Owner, msg.Beneficiary, msg.Price); err != nil {
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

	priceAfterCommission, err := doCommissions(ctx, k, msg.Buyer, msg.Beneficiary, nft.SellerBeneficiary, nft.GetPrice())
	if err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: failed to pay commissions: %v", err)),
		}
	}

	err = k.coinKeeper.SendCoins(ctx, msg.Buyer, nft.GetOwner(), priceAfterCommission)
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

func doCommissions(
	ctx sdk.Context,
	k Keeper,
	payer,
	sellerBeneficiary,
	buyerBeneficiary sdk.AccAddress,
	price sdk.Coins,
) (priceAfterCommission sdk.Coins, err error) {
	logger := ctx.Logger()

	// Check that payer has enough funds (for both the commission and the asset itself).
	totalCommission := getCommission(price, ValidatorsCommission+BeneficiariesCommission)
	logger.Info("calculated total commission", "total_commission", totalCommission.String())

	priceAfterCommission = price.Sub(totalCommission)
	if !k.coinKeeper.HasCoins(ctx, payer, totalCommission) {
		return nil, fmt.Errorf("user %s does not have enough funds", payer.String())
	}

	// Pay commission to the beneficiaries.
	beneficiaryCommission := getCommission(price, BeneficiariesCommission/2)
	logger.Info("calculated beneficiary commission", "beneficiary_commission", beneficiaryCommission.String())
	if err := k.coinKeeper.SendCoins(ctx, payer, sellerBeneficiary, beneficiaryCommission); err != nil {
		return nil, fmt.Errorf("failed to pay commission to beneficiary: %v", err)
	}
	logger.Info("payed seller beneficiary commission", "seller_beneficiary", sellerBeneficiary.String())
	if err := k.coinKeeper.SendCoins(ctx, payer, buyerBeneficiary, beneficiaryCommission); err != nil {
		return nil, fmt.Errorf("failed to pay commission to beneficiary: %v", err)
	}
	logger.Info("payed buyer beneficiary commission", "buyer_beneficiary", buyerBeneficiary.String())

	votes := ctx.VoteInfos()
	var vals []abci_types.Validator
	for _, vote := range votes {
		if vote.SignedLastBlock {
			vals = append(vals, vote.Validator)
		}
	}

	// First we take tokens from the payer, then we allocate tokens to validators via distribution module.
	totalValsCommission := getCommission(price, BeneficiariesCommission)
	if _, err := k.coinKeeper.SubtractCoins(ctx, payer, totalValsCommission); err != nil {
		return nil, fmt.Errorf("failed to take validators commission from payer: %v", err)
	}
	logger.Info("wrote off validators commission")
	singleValCommission := getCommission(price, BeneficiariesCommission/float64(len(vals)))
	logger.Info("paying validators", "validator_commission", singleValCommission.String(),
		"num_validators", len(vals))
	for _, val := range vals {
		consVal := k.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(val.Address))
		k.distrKeeper.AllocateTokensToValidator(ctx, consVal, sdk.NewDecCoins(singleValCommission))
	}

	return priceAfterCommission, nil
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
