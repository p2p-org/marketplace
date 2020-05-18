package marketplace

import (
	"fmt"
	"strconv"

	"github.com/p2p-org/marketplace/common"
	"github.com/p2p-org/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/google/uuid"
	abci_types "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/types/time"
)

// NewHandler returns a handler for "marketplace" type messages.
func NewHandler(keeper *Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgPutNFTOnMarket:
			return handleMsgPutNFTOnMarket(ctx, keeper, msg)
		case MsgRemoveNFTFromMarket:
			return handleMsgRemoveNFTFromMarket(ctx, keeper, msg)
		case MsgBuyNFT:
			return handleMsgBuyNFT(ctx, keeper, msg)
		case MsgPutNFTOnAuction:
			return handleMsgPutNFTOnAuction(ctx, keeper, msg)
		case MsgRemoveNFTFromAuction:
			return handleMsgRemoveNFTFromAuction(ctx, keeper, msg)
		case MsgMakeBidOnAuction:
			return handleMsgMakeBidOnAuction(ctx, keeper, msg)
		case MsgFinishAuction:
			return handleMsgFinishAuction(ctx, keeper, msg)
		case MsgBuyoutOnAuction:
			return handleMsgBuyoutOnAuction(ctx, keeper, msg)
		case MsgBatchTransfer:
			return handleMsgBatchTransfer(ctx, keeper, msg)
		case MsgBatchPutOnMarket:
			return handleMsgBatchPutOnMarket(ctx, keeper, msg)
		case MsgBatchRemoveFromMarket:
			return handleMsgBatchRemoveFromMarket(ctx, keeper, msg)
		case MsgBatchBuyOnMarket:
			return handleMsgBatchBuyOnMarket(ctx, keeper, msg)
		case MsgMakeOffer:
			return handleMsgMakeOffer(ctx, keeper, msg)
		case MsgAcceptOffer:
			return handleMsgAcceptOffer(ctx, keeper, msg)
		case MsgRemoveOffer:
			return handleMsgRemoveOffer(ctx, keeper, msg)
		case MsgUpdateNFTParams:
			return handleMsgUpdateNFTParams(ctx, keeper, msg)
		case MsgCreateFungibleToken:
			return handleMsgCreateFungibleTokensCurrency(ctx, keeper, msg)
		case MsgTransferFungibleTokens:
			return handleMsgTransferFungibleTokens(ctx, keeper, msg)
		case MsgBurnFungibleToken:
			return handleMsgBurnFungibleToken(ctx, keeper, msg)
		//case MsgTransferNFTByIBC:
		//	return HandleMsgTransferNFTByIBC(ctx, keeper, msg)
		default:
			errMsg := fmt.Sprintf("Unrecognized marketplace Msg type: %v", msg.Type())
			return nil, fmt.Errorf(errMsg)
		}
	}
}

func handleMsgPutNFTOnMarket(ctx sdk.Context, mpKeeper *Keeper, msg MsgPutNFTOnMarket) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgPutNFTOnMarket)

	if !mpKeeper.IsDenomExist(ctx, msg.Price) {
		return nil, fmt.Errorf("failed to PutNFTOnMarket: denom does not exist")
	}

	if err := mpKeeper.PutNFTOnMarket(ctx, msg.TokenID, msg.Owner, msg.Beneficiary, msg.Price); err != nil {
		return nil, fmt.Errorf("failed to PutNFTOnMarket: %v", err)
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgPutNFTOnMarket)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.Beneficiary.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRemoveNFTFromMarket(ctx sdk.Context, mpKeeper *Keeper, msg MsgRemoveNFTFromMarket) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgRemoveNFTFromMarket)
	if err := mpKeeper.RemoveNFTFromMarket(ctx, msg.TokenID, msg.Owner); err != nil {
		return nil, fmt.Errorf("failed to RemoveNFTFromMarket: %v", err)
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgRemoveNFTFromMarket)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBuyNFT(ctx sdk.Context, mpKeeper *Keeper, msg MsgBuyNFT) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBuyNFT)
	token, err := mpKeeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to BuyNFT: %v", err)
	}

	if !token.IsOnMarket() {
		return nil, fmt.Errorf("failed to BuyNFT: token is not for sale")
	}

	beneficiariesCommission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > mpKeeper.config.MaximumBeneficiaryCommission {
		return nil, fmt.Errorf("failed to BuyNFT: beneficiary commission is too high")
	}

	priceAfterCommission, err := doNFTCommissions(
		ctx,
		mpKeeper,
		msg.Buyer,
		token.Owner,
		msg.Beneficiary,
		token.SellerBeneficiary,
		token.GetPrice(),
		beneficiariesCommission,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to BuyNFT: failed to pay commissions: %v", err)
	}

	err = mpKeeper.coinKeeper.SendCoins(ctx, msg.Buyer, token.Owner, priceAfterCommission)
	if err != nil {
		return nil, fmt.Errorf("Buyer does not have enough coins")
	}

	token.Owner = msg.Buyer
	token.SetSellerBeneficiary(sdk.AccAddress{})
	token.SetStatus(types.NFTStatusDefault)

	if err := mpKeeper.UpdateNFT(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to BuyNFT: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgBuyNFT)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyBuyer, msg.Buyer.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.Beneficiary.String()),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.BeneficiaryCommission),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgMakeOffer(ctx sdk.Context, mpKeeper *Keeper, msg MsgMakeOffer) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMakeOffer)

	if !mpKeeper.coinKeeper.HasBalance(ctx, msg.Buyer, msg.Price[0]) {
		return nil, fmt.Errorf("buyer does not have the offered funds")
	}

	token, err := mpKeeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to MakeOffer: %v", err)
	}

	token.AddOffer(&types.Offer{
		ID:                    uuid.New().String(),
		Price:                 msg.Price,
		Buyer:                 msg.Buyer,
		BuyerBeneficiary:      msg.BuyerBeneficiary,
		BeneficiaryCommission: msg.BeneficiaryCommission,
	})

	if _, err := mpKeeper.coinKeeper.SubtractCoins(ctx, msg.Buyer, msg.Price); err != nil {
		return nil, err
	}

	if err := mpKeeper.UpdateNFT(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to MakeOffer: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMakeOffer)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOfferID, token.Offers[len(token.Offers)-1].ID),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyPrice, msg.Price.String()),
			sdk.NewAttribute(types.AttributeKeyBuyer, msg.Buyer.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.BuyerBeneficiary.String()),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.BeneficiaryCommission),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgAcceptOffer(ctx sdk.Context, mpKeeper *Keeper, msg MsgAcceptOffer) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgAcceptOffer)
	token, err := mpKeeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to AcceptOffer: %v", err)
	}

	offer, ok := token.GetOffer(msg.OfferID)
	if !ok {
		return nil, fmt.Errorf("failed to AcceptOffer: no ofer with ID %s", msg.OfferID)
	}

	beneficiariesCommission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > mpKeeper.config.MaximumBeneficiaryCommission {
		return nil, fmt.Errorf("failed to AcceptOffer: beneficiary commission is too high")
	}

	if token.IsOnMarket() {
		err := mpKeeper.RemoveNFTFromMarket(ctx, msg.TokenID, msg.Seller)
		if err != nil {
			return nil, fmt.Errorf("failed to AcceptOffer: could not remove token from market")
		}
	}

	if token.IsOnAuction() {
		lot, err := mpKeeper.GetAuctionLot(ctx, msg.TokenID)
		if err != nil {
			return nil, fmt.Errorf("failed to AcceptOffer: could not get auction lot")
		}

		if lot.ExpirationTime.Before(time.Now().UTC()) {
			return nil, fmt.Errorf("auction is already finished")
		}

		// return bid to last bidder if exists
		if lot.LastBid != nil {
			_, err := mpKeeper.coinKeeper.AddCoins(ctx, lot.LastBid.Bidder, lot.LastBid.Bid)
			if err != nil {
				return nil, fmt.Errorf("failed to AcceptOffer: could not get return coins to bidder")
			}
		}

		err = mpKeeper.RemoveNFTFromAuction(ctx, msg.TokenID, msg.Seller)
		if err != nil {
			return nil, fmt.Errorf("failed to AcceptOffer: could not remove token from market")
		}
	}

	// Return frozen funds to the buyer so that doNFTCommissions works correctly
	if _, err = mpKeeper.coinKeeper.AddCoins(ctx, offer.Buyer, offer.Price); err != nil {
		return nil, err
	}

	priceAfterCommission, err := doNFTCommissions(
		ctx,
		mpKeeper,
		offer.Buyer,
		token.Owner,
		offer.BuyerBeneficiary,
		msg.SellerBeneficiary,
		offer.Price,
		beneficiariesCommission,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to AcceptOffer: failed to pay commissions: %v", err)
	}

	if err = mpKeeper.coinKeeper.SendCoins(ctx, offer.Buyer, token.Owner, priceAfterCommission); err != nil {
		return nil, fmt.Errorf("failed to AcceptOffer: buyer does not have enough coins")
	}

	if ok := token.RemoveOffer(offer.ID, offer.Buyer); !ok {
		return nil, fmt.Errorf("failed to AcceptOffer: no ofer with ID %s", msg.OfferID)
	}

	token.Owner = offer.Buyer
	token.SetSellerBeneficiary(sdk.AccAddress{})
	token.SetStatus(types.NFTStatusDefault)

	if err := mpKeeper.UpdateNFT(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to AcceptOffer: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgAcceptOffer)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.BeneficiaryCommission),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, offer.Buyer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRemoveOffer(ctx sdk.Context, mpKeeper *Keeper, msg MsgRemoveOffer) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgRemoveOffer)
	token, err := mpKeeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to RemoveOffer: %v", err)
	}

	offer, ok := token.GetOffer(msg.OfferID)
	if !ok {
		return nil, fmt.Errorf("failed to RemoveOffer: no ofer with ID %s", msg.OfferID)
	}

	ok = token.RemoveOffer(msg.OfferID, msg.Buyer)
	if !ok {
		return nil, fmt.Errorf("failed to RemoveOffer: no ofer with ID %s", msg.OfferID)
	}

	if _, err := mpKeeper.coinKeeper.AddCoins(ctx, msg.Buyer, offer.Price); err != nil {
		return nil, err
	}

	if err := mpKeeper.UpdateNFT(ctx, token); err != nil {
		return nil, fmt.Errorf("failed to RemoveOffer: %v", err)
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgRemoveOffer)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyOfferID, msg.OfferID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgUpdateNFTParams(ctx sdk.Context, mpKeeper *Keeper, msg MsgUpdateNFTParams) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgUpdateNFTParams)
	nft, err := mpKeeper.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get nft in UpdateNFTParams: %v", err)
	}
	if !nft.Owner.Equals(msg.Owner) {
		return nil, fmt.Errorf("user is not an owner: %v", msg.Owner.String())
	}

	for _, v := range msg.Params {
		v := v
		switch v.Key {
		case types.FlagParamPrice:
			price, err := sdk.ParseCoins(v.Value)
			if err != nil {
				return nil, fmt.Errorf("failed to UpdateNFTParams.Price: %v", err)
			}
			if !mpKeeper.IsDenomExist(ctx, price) {
				return nil, fmt.Errorf("failed to UpdateNFTParams.Price: denom is not registered")

			}
			nft.Price = price
		}
	}

	if err := mpKeeper.UpdateNFT(ctx, nft); err != nil {
		return nil, fmt.Errorf("failed to UpdateNFTParams: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgUpdateNFTParams)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.Params.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgCreateFungibleTokensCurrency(ctx sdk.Context, mpKeeper *Keeper, msg MsgCreateFungibleToken) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgCreateFungibleToken)
	if err := mpKeeper.CreateFungibleToken(ctx, msg.Creator, msg.Denom, msg.Amount); err != nil {
		return nil, fmt.Errorf("failed to create currency: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgCreateFungibleToken)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Creator.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatInt(msg.Amount, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Creator.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgTransferFungibleTokens(ctx sdk.Context, mpKeeper *Keeper, msg MsgTransferFungibleTokens) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgTransferFungibleTokens)
	if err := mpKeeper.TransferFungibleTokens(ctx, msg.Owner, msg.Recipient, msg.Denom, msg.Amount); err != nil {
		return nil, fmt.Errorf("failed to transfer coins: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgTransferFungibleTokens)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyRecipient, msg.Recipient.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatInt(msg.Amount, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBurnFungibleToken(ctx sdk.Context, mpKeeper *Keeper, msg MsgBurnFungibleToken) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBurnFT)
	if err := mpKeeper.BurnFungibleTokens(ctx, msg.Owner, msg.Denom, msg.Amount); err != nil {
		return nil, fmt.Errorf("failed to burn coins: %v", err)
	}
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBurnFT)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyDenom, msg.Denom),
			sdk.NewAttribute(types.AttributeKeyAmount, strconv.FormatInt(msg.Amount, 10)),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func doNFTCommissions(
	ctx sdk.Context,
	k *Keeper,
	buyer,
	seller,
	sellerBeneficiary,
	buyerBeneficiary sdk.AccAddress,
	price sdk.Coins,
	beneficiariesCommission float64,
) (priceAfterCommission sdk.Coins, err error) {
	logger := ctx.Logger()

	// Check that buyer has enough funds (for both the commission and the asset itself).
	if !k.coinKeeper.HasBalance(ctx, buyer, price[0]) {
		return nil, fmt.Errorf("user %s does not have enough funds", buyer.String())
	}
	logger.Info("user has enough funds, o.nftKeeper.")

	votes := ctx.VoteInfos()
	var vals []abci_types.Validator
	for _, vote := range votes {
		if vote.SignedLastBlock {
			vals = append(vals, vote.Validator)
		}
	}
	lenVals := float64(len(vals))
	if len(vals) == 0 {
		lenVals = 1.0
	}
	// first calculate all commissions and total commission as sum of them
	singleValCommission := GetCommission(price, types.DefaultValidatorsCommission/lenVals)
	totalValsCommission := sdk.NewCoins()
	for i := 0; i < int(lenVals); i++ {
		totalValsCommission = totalValsCommission.Add(singleValCommission...)
	}

	totalCommission := sdk.NewCoins()
	beneficiaryCommission := GetCommission(price, beneficiariesCommission/2)
	logger.Info("calculated beneficiary commission", "beneficiary_commission", beneficiaryCommission.String())

	totalCommission = totalCommission.Add(beneficiaryCommission...)
	totalCommission = totalCommission.Add(beneficiaryCommission...)
	totalCommission = totalCommission.Add(totalValsCommission...)

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
		k.distrKeeper.AllocateTokensToValidator(ctx, consVal, sdk.NewDecCoinsFromCoins(singleValCommission...))
	}

	return priceAfterCommission, nil
}

type balance struct {
	addr   sdk.AccAddress
	amount sdk.Coins
}

func GetBalances(ctx sdk.Context, mpKeeper *Keeper, addrs ...sdk.AccAddress) []*balance {
	var out []*balance
	for _, addr := range addrs {
		out = append(out, &balance{
			addr:   addr,
			amount: mpKeeper.coinKeeper.GetAllBalances(ctx, addr),
		})
	}

	return out
}

func RollbackCommissions(ctx sdk.Context, mpKeeper *Keeper, logger log.Logger, initialBalances []*balance) {
	for _, balance := range initialBalances {
		if err := mpKeeper.coinKeeper.SetBalances(ctx, balance.addr, balance.amount); err != nil {
			logger.Error("failed to rollback commissions", "addr", balance.addr.String(), "error", err)
		}
	}
}

func calculateNumAndDenom(p float64) (sdk.Dec, sdk.Dec) {
	if p == 0 {
		return sdk.NewDec(0), sdk.NewDec(1)
	}
	/*
		//	Considering a float64 less than 1.0 (e.g. 0.0015)
		//	as a quotient (fraction) of two numbers p/q
		//	(p (numerator) divided by q (denominator)), where
		//	p is an integer number (e.g. 15) and
		//	q is an integer number, which is product of 10 (e.g. 10000),
		//	we can express input float64 value of commission
		//	this way (as p/q)
		//	This is supposed to simplify commission calculation
	*/
	//	init q
	q := int64(1)
	//	for loop till the precision limit
	for i := 0; i < sdk.Precision; i++ {
		//	multiply input float number
		p *= 10
		//	multiply q as well
		q *= 10
		//	check if input number became integer
		//	this if faster than
		//	math.Trunc(p) == p
		if float64(int64(p)) == p {
			break
		}
	}
	return sdk.NewDec(int64(p)), sdk.NewDec(q)
}

func GetCommission(price sdk.Coins, rat64 float64) sdk.Coins {
	if rat64 >= 1 {
		return price
	}
	num, denom := calculateNumAndDenom(rat64)
	priceDec := sdk.NewDecCoinsFromCoins(price...)
	totalCommission, _ := priceDec.MulDec(num).QuoDec(denom).TruncateDecimal()
	return totalCommission
}

func handleMsgBatchTransfer(ctx sdk.Context, mpKeeper *Keeper, msg MsgBatchTransfer) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBatchTransfer)

	for _, tokenID := range msg.TokenIDs {
		token, err := mpKeeper.GetNFT(ctx, tokenID)
		if err != nil {
			return nil, fmt.Errorf("failed to find token %s: %v", tokenID, err)
		}
		_, err = HandleMsgTransferNFTMarketplace(ctx, nft.MsgTransferNFT{
			Sender:    msg.Sender,
			Recipient: msg.Recipient,
			Denom:     token.Denom,
			ID:        tokenID,
		}, mpKeeper.nftKeeper, mpKeeper)
		if err != nil {
			ctx.Logger().Info("batch transfer error, tokenID:", tokenID)
			continue
		}
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgBatchTransfer)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBatchPutOnMarket(ctx sdk.Context, mpKeeper *Keeper, msg MsgBatchPutOnMarket) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMsgBatchPutOnMarket)

	for k, price := range msg.TokenPrices {
		k := k
		tokenID := msg.TokenIDs[k]
		price := price

		_, err := handleMsgPutNFTOnMarket(ctx, mpKeeper, MsgPutNFTOnMarket{
			Owner:       msg.Owner,
			Beneficiary: msg.Beneficiary,
			TokenID:     tokenID,
			Price:       price,
		})
		if err != nil {
			ctx.Logger().Info("batch put on market error, tokenID:", tokenID)
			return nil, err
		}
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMsgBatchPutOnMarket)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBatchRemoveFromMarket(ctx sdk.Context, mpKeeper *Keeper, msg MsgBatchRemoveFromMarket) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMsgBatchRemoveFromMarket)

	for _, tokenID := range msg.TokenIDs {
		_, err := handleMsgRemoveNFTFromMarket(ctx, mpKeeper, MsgRemoveNFTFromMarket{
			Owner:   msg.Owner,
			TokenID: tokenID,
		})
		if err != nil {
			ctx.Logger().Info("batch remove from market error, tokenID:", tokenID, "result:")
			continue
		}
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMsgBatchRemoveFromMarket)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBatchBuyOnMarket(ctx sdk.Context, mpKeeper *Keeper, msg MsgBatchBuyOnMarket) (*sdk.Result, error) {
	mpKeeper.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMsgBatchBuyOnMarket)

	beneficiariesCommission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > mpKeeper.config.MaximumBeneficiaryCommission {
		return nil, fmt.Errorf("failed to BuyNFT: beneficiary commission is too high")
	}

	priceSum := sdk.NewCoins()

	for _, tokenID := range msg.TokenIDs {
		tokenID := tokenID
		token, err := mpKeeper.GetNFT(ctx, tokenID)
		if err != nil {
			return nil, fmt.Errorf("failed to BuyNFT: %v", err)
		}
		if !token.IsOnMarket() {
			return nil, fmt.Errorf("failed to buy: token %v is not on market", token.ID)
		}
		priceSum = priceSum.Add(token.Price...)
	}

	buyerCoins := mpKeeper.coinKeeper.GetAllBalances(ctx, msg.Buyer)
	if !buyerCoins.IsAllGTE(priceSum) {
		return nil, fmt.Errorf("failed to buy batch: not enough funds")
	}

	for _, tokenID := range msg.TokenIDs {
		tokenID := tokenID

		_, err := handleMsgBuyNFT(ctx, mpKeeper, MsgBuyNFT{
			Buyer:       msg.Buyer,
			Beneficiary: msg.Beneficiary,
			TokenID:     tokenID,
		})
		if err != nil {
			ctx.Logger().Info("batch buy error, tokenID:", tokenID, "result:")
			continue
		}
	}

	mpKeeper.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMsgBatchBuyOnMarket)
	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})

	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
