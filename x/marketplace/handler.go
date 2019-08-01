package marketplace

import (
	"fmt"
	"math/big"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	mptypes "github.com/dgamingfoundation/marketplace/x/marketplace/types"
	abci_types "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
)

// NewHandler returns a handler for "marketplace" type messages.
func NewHandler(keeper Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		switch msg := msg.(type) {
		case MsgMintNFT:
			return handleMsgMintNFT(ctx, keeper, msg)
		case MsgTransferNFT:
			return handleMsgTransferNFT(ctx, keeper, msg)
		case MsgPutNFTOnMarket:
			return handleMsgPutNFTOnMarket(ctx, keeper, msg)
		case MsgBuyNFT:
			return handleMsgBuyNFT(ctx, keeper, msg)
		case MsgCreateFungibleToken:
			return handleMsgCreateFungibleTokensCurrency(ctx, keeper, msg)
		case MsgTransferFungibleTokens:
			return handleMsgTransferFungibleTokens(ctx, keeper, msg)
		case MsgUpdateNFTParams:
			return handleMsgUpdateNFTParams(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized marketplace Msg type: %v", msg.Type())
			return sdk.ErrUnknownRequest(errMsg).Result()
		}
	}
}

func handleMsgCreateFungibleTokensCurrency(ctx sdk.Context, keeper Keeper, msg MsgCreateFungibleToken) sdk.Result {
	if err := keeper.CreateFungibleToken(ctx, msg.Creator, msg.Denom, msg.Amount); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to create currency: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgTransferFungibleTokens(ctx sdk.Context, keeper Keeper, msg MsgTransferFungibleTokens) sdk.Result {
	if err := keeper.TransferFungibleTokens(ctx, msg.Owner, msg.Recipient, msg.Denom, msg.Amount); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to transfer coins: %v", err)),
		}
	}
	return sdk.Result{}
}

func handleMsgMintNFT(ctx sdk.Context, keeper Keeper, msg MsgMintNFT) sdk.Result {
	nft := NewNFT(
		xnft.NewBaseNFT(
			msg.TokenID,
			msg.Owner,
			msg.Name,
			msg.Description,
			msg.Image,
			msg.TokenURI,
		),
		sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.NewInt(0))),
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

func handleMsgPutNFTOnMarket(ctx sdk.Context, k Keeper, msg MsgPutNFTOnMarket) sdk.Result {
	if !k.IsDenomExist(ctx, msg.Price) {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to PutNFTOnMarket: %v", "denom does not exist")),
		}

	}
	if err := k.PutNFTOnMarket(ctx, msg.TokenID, msg.Owner, msg.Beneficiary, msg.Price); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to PutNFTOnMarket: %v", err)),
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

	beneficiariesCommission := mptypes.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > k.config.MaximumBeneficiaryCommission {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to BuyNFT: beneficiary commission is too high")),
		}
	}

	priceAfterCommission, err := doNFTCommissions(
		ctx,
		k,
		msg.Buyer,
		nft.Owner,
		msg.Beneficiary,
		nft.SellerBeneficiary,
		nft.GetPrice(),
		beneficiariesCommission,
	)
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

	nft.BaseNFT = nft.SetOwner(msg.Buyer)
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

func doNFTCommissions(
	ctx sdk.Context,
	k Keeper,
	buyer,
	seller,
	sellerBeneficiary,
	buyerBeneficiary sdk.AccAddress,
	price sdk.Coins,
	beneficiariesCommission float64,
) (priceAfterCommission sdk.Coins, err error) {
	logger := ctx.Logger()

	// Check that buyer has enough funds (for both the commission and the asset itself).
	if !k.coinKeeper.HasCoins(ctx, buyer, price) {
		return nil, fmt.Errorf("user %s does not have enough funds", buyer.String())
	}
	logger.Info("user has enough funds, o.k.")

	votes := ctx.VoteInfos()
	var vals []abci_types.Validator
	for _, vote := range votes {
		if vote.SignedLastBlock {
			vals = append(vals, vote.Validator)
		}
	}

	// first calculate all commissions and total commission as sum of them
	singleValCommission := GetCommission(price, mptypes.DefaultValidatorsCommission/float64(len(vals)))
	totalValsCommission := sdk.NewCoins()
	for i := 0; i < len(vals); i++ {
		totalValsCommission = totalValsCommission.Add(singleValCommission)
	}

	totalCommission := sdk.NewCoins()
	beneficiaryCommission := GetCommission(price, beneficiariesCommission/2)
	logger.Info("calculated beneficiary commission", "beneficiary_commission", beneficiaryCommission.String())

	totalCommission = totalCommission.Add(beneficiaryCommission)
	totalCommission = totalCommission.Add(beneficiaryCommission)
	totalCommission = totalCommission.Add(totalValsCommission)

	priceAfterCommission = price.Sub(totalCommission)
	logger.Info("calculated total commission", "total_commission", totalCommission.String(),
		"price_after_commission", priceAfterCommission.String())

	var initialBalances = GetBalances(ctx, k, buyer, seller, buyerBeneficiary, sellerBeneficiary)
	// Pay commission to the beneficiaries.
	if err := k.coinKeeper.SendCoins(ctx, buyer, sellerBeneficiary, beneficiaryCommission); err != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return nil, fmt.Errorf("failed to pay commission to beneficiary: %v", err)
	}
	logger.Info("payed seller beneficiary commission", "seller_beneficiary", sellerBeneficiary.String())
	if err := k.coinKeeper.SendCoins(ctx, buyer, buyerBeneficiary, beneficiaryCommission); err != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return nil, fmt.Errorf("failed to pay commission to beneficiary: %v", err)
	}
	logger.Info("payed buyer beneficiary commission", "buyer_beneficiary", buyerBeneficiary.String())

	// First we take tokens from the buyer, then we allocate tokens to validators via distribution module.
	if _, err := k.coinKeeper.SubtractCoins(ctx, buyer, totalValsCommission); err != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return nil, fmt.Errorf("failed to take validators commission from buyer: %v", err)
	}
	logger.Info("wrote off validators commission")

	logger.Info("paying validators", "validator_commission", singleValCommission.String(),
		"num_validators", len(vals))
	for _, val := range vals {
		consVal := k.stakingKeeper.ValidatorByConsAddr(ctx, sdk.ConsAddress(val.Address))
		k.distrKeeper.AllocateTokensToValidator(ctx, consVal, sdk.NewDecCoins(singleValCommission))
	}

	return priceAfterCommission, nil
}

type balance struct {
	addr   sdk.AccAddress
	amount sdk.Coins
}

func GetBalances(ctx sdk.Context, k Keeper, addrs ...sdk.AccAddress) []*balance {
	var out []*balance
	for _, addr := range addrs {
		out = append(out, &balance{
			addr:   addr,
			amount: k.coinKeeper.GetCoins(ctx, addr),
		})
	}

	return out
}

func RollbackCommissions(ctx sdk.Context, k Keeper, logger log.Logger, initialBalances []*balance) {
	for _, balance := range initialBalances {
		if err := k.coinKeeper.SetCoins(ctx, balance.addr, balance.amount); err != nil {
			logger.Error("failed to rollback commissions", "addr", balance.addr.String(), "error", err)
		}
	}
}

func GetCommission(price sdk.Coins, rat64 float64) sdk.Coins {
	// TODO: maybe we can do it somehow easier.
	var rat = new(big.Rat)
	rat = rat.SetFloat64(rat64)
	num, denom := sdk.NewDecFromBigInt(rat.Num()), sdk.NewDecFromBigInt(rat.Denom())
	priceDec := sdk.NewDecCoins(price)
	totalCommission, _ := priceDec.MulDec(num).QuoDec(denom).TruncateDecimal()
	return totalCommission
}

func handleMsgUpdateNFTParams(ctx sdk.Context, k Keeper, msg MsgUpdateNFTParams) sdk.Result {
	nft, err := k.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to get nft in UpdateNFTParams: %v", err)),
		}
	}
	if !nft.Owner.Equals(msg.Owner) {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("user is not an owner: %v", msg.Owner.String())),
		}
	}

	for _, v := range msg.Params {
		v := v
		switch v.Key {
		case types.FlagParamPrice:
			price, err := sdk.ParseCoins(v.Value)
			if err != nil {
				return sdk.Result{
					Code:      sdk.CodeUnknownRequest,
					Codespace: "marketplace",
					Data:      []byte(fmt.Sprintf("failed to UpdateNFTParams.Price: %v", err)),
				}
			}
			if !k.IsDenomExist(ctx, price) {
				return sdk.Result{
					Code:      sdk.CodeUnknownRequest,
					Codespace: "marketplace",
					Data:      []byte(fmt.Sprintf("failed to UpdateNFTParams.Price: %v", "denom is not registered")),
				}
			}
			nft.Price = price
		case types.FlagParamDescription:
			nft.Description = v.Value
		case types.FlagParamTokenName:
			nft.Name = v.Value
		case types.FlagParamTokenURI:
			nft.TokenURI = v.Value
		case types.FlagParamImage:
			nft.Image = v.Value
		}
	}

	if err := k.UpdateNFT(ctx, nft); err != nil {
		return sdk.Result{
			Code:      sdk.CodeUnknownRequest,
			Codespace: "marketplace",
			Data:      []byte(fmt.Sprintf("failed to UpdateNFTParams: %v", err)),
		}
	}

	return sdk.Result{}
}
