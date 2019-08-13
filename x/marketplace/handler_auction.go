package marketplace

import (
	"fmt"
	"strconv"

	"github.com/tendermint/tendermint/types/time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

func handleMsgPutNFTOnAuction(ctx sdk.Context, k Keeper, msg types.MsgPutNFTOnAuction) sdk.Result {
	failMsg := "failed to RemoveNFTFromAuction"
	if !k.IsDenomExist(ctx, msg.OpeningPrice) {
		return wrapError(failMsg, fmt.Errorf("failed to PutNFTOnAuction: %v", "denom does not exist"))
	}

	if !msg.BuyoutPrice.IsZero() && !k.IsDenomExist(ctx, msg.BuyoutPrice) {
		return wrapError(failMsg, fmt.Errorf("failed to PutNFTOnAuction: %v", "denom does not exist"))
	}

	if err := k.PutNFTOnAuction(ctx, msg.TokenID, msg.Owner, msg.Beneficiary, msg.OpeningPrice,
		msg.BuyoutPrice, msg.TimeToSell); err != nil {
		return wrapError(failMsg, fmt.Errorf("failed to PutNFTOnAuction: %v", err))
	}

	return sdk.Result{}
}

func handleMsgRemoveNFTFromAuction(ctx sdk.Context, k Keeper, msg MsgRemoveNFTFromAuction) sdk.Result {
	failMsg := "failed to RemoveNFTFromAuction"
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return wrapError(failMsg, err)
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return wrapError(failMsg, fmt.Errorf("auction is already finished"))
	}

	// return bid to last bidder if exists
	if lot.LastBid != nil {
		_, err := k.coinKeeper.AddCoins(ctx, lot.LastBid.Bidder, lot.LastBid.Bid)
		if err != nil {
			return wrapError(failMsg, err)
		}
	}

	// return nft to owner, delete lot
	if err := k.RemoveNFTFromAuction(ctx, msg.TokenID, msg.Owner); err != nil {
		return wrapError(failMsg, err)
	}

	return sdk.Result{}
}

func handleMsgFinishAuction(ctx sdk.Context, k Keeper, msg MsgFinishAuction) sdk.Result {
	failMsg := "failed to FinishAuction"
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return wrapError(failMsg, err)
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return wrapError(failMsg, fmt.Errorf("auction is already finished"))
	}

	// no bids on lot
	if lot.LastBid != nil {
		if err := k.BuyLotOnAuction(ctx, lot.LastBid.Bidder, lot.LastBid.BuyerBeneficiary, lot.LastBid.Bid, lot); err != nil {
			return wrapError(failMsg, err)
		}
	} else {
		if err := k.RemoveNFTFromAuction(ctx, msg.TokenID, msg.Owner); err != nil {
			return wrapError(failMsg, err)
		}
	}

	return sdk.Result{}
}

func handleMsgMakeBidOnAuction(ctx sdk.Context, k Keeper, msg MsgMakeBidOnAuction) sdk.Result {
	logger := ctx.Logger()
	failMsg := "failed to MakeBidOnAuction"
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return wrapError(failMsg, err)
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return wrapError(failMsg, fmt.Errorf("auction is already finished"))
	}

	beneficiariesCommission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > k.config.MaximumBeneficiaryCommission {
		return wrapError(failMsg, fmt.Errorf("failed to BuyNFT: beneficiary commission is too high"))
	}
	beneficiariesCommissionString := fmt.Sprintf("%v", beneficiariesCommission)

	// bid is less than lastBid
	if lot.LastBid != nil {
		if msg.Bid.IsAllLTE(lot.LastBid.Bid) {
			return wrapError(failMsg, fmt.Errorf("bid: %+v is lower than last bid: %+v", msg.Bid, lot.LastBid.Bid))
		}
	}

	// bid is less than opening price
	if lot.OpeningPrice.IsAnyGT(msg.Bid) {
		return wrapError(failMsg, fmt.Errorf("bid: %+v is lower than opening price: %+v", msg.Bid, lot.OpeningPrice))
	}

	// no buyout, change lastBid
	balances := GetBalances(ctx, k, msg.Bidder, msg.BuyerBeneficiary)
	if lot.LastBid != nil {
		balances = append(balances,
			GetBalances(ctx, k, lot.LastBid.Bidder, lot.LastBid.BuyerBeneficiary)...)
		// return coins to previous bidder
		_, err = k.coinKeeper.AddCoins(ctx, lot.LastBid.Bidder, lot.LastBid.Bid)
		if err != nil {
			RollbackCommissions(ctx, k, logger, balances)
			return wrapError(failMsg, err)
		}
	}

	// take coins from new bidder
	_, err = k.coinKeeper.SubtractCoins(ctx, msg.Bidder, msg.Bid)
	if err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return wrapError(failMsg, err)
	}

	auctionBid := types.NewAuctionBid(msg.Bidder, msg.BuyerBeneficiary, msg.Bid, beneficiariesCommissionString)
	lot.SetLastBid(auctionBid)

	if err := k.UpdateAuctionLot(ctx, lot); err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return wrapError(failMsg, err)
	}

	// bid is more than buyout price. perform buyout
	if !lot.BuyoutPrice.IsZero() {
		if msg.Bid.IsAllGTE(lot.BuyoutPrice) {
			err = k.BuyLotOnAuction(ctx, msg.Bidder, msg.BuyerBeneficiary, msg.Bid, lot)
			if err != nil {
				return wrapError(failMsg, err)
			}
			return sdk.Result{}
		}
	}

	return sdk.Result{}
}

func handleMsgBuyoutOnAuction(ctx sdk.Context, k Keeper, msg MsgBuyoutOnAuction) sdk.Result {
	failMsg := "failed to MakeBidOnAuction"
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return wrapError(failMsg, err)
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return wrapError(failMsg, fmt.Errorf("auction is already finished"))
	}

	if lot.BuyoutPrice.IsZero() {
		return wrapError(failMsg, fmt.Errorf("lot has no buyoutprice"))
	}

	err = k.BuyLotOnAuction(ctx, msg.Buyer, msg.BuyerBeneficiary, lot.BuyoutPrice, lot)
	if err != nil {
		return wrapError(failMsg, err)
	}

	return sdk.Result{}
}

func wrapError(failMsg string, err error) sdk.Result {
	return sdk.Result{
		Code:      sdk.CodeUnknownRequest,
		Codespace: "marketplace",
		Data:      []byte(fmt.Sprintf("%s: %v", failMsg, err)),
	}
}
