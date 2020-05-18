package marketplace

import (
	"fmt"
	"strconv"

	"github.com/p2p-org/marketplace/common"
	"github.com/p2p-org/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/tendermint/types/time"
)

func handleMsgPutNFTOnAuction(ctx sdk.Context, k *Keeper, msg types.MsgPutNFTOnAuction) (*sdk.Result, error) {
	k.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgPutNFTOnAuction)

	if !k.IsDenomExist(ctx, msg.OpeningPrice) {
		return nil, fmt.Errorf("failed to PutNFTOnAuction: %v", "denom does not exist")
	}

	if !msg.BuyoutPrice.IsZero() {
		if !k.IsDenomExist(ctx, msg.BuyoutPrice) {
			return nil, fmt.Errorf("failed to PutNFTOnAuction: %v", "denom does not exist")
		}
		if msg.OpeningPrice.IsAnyGTE(msg.BuyoutPrice) {
			return nil, fmt.Errorf("failed to PutNFTOnAuction: %v", "buyout price is too low")
		}
	}

	if err := k.PutNFTOnAuction(ctx, msg.TokenID, msg.Owner, msg.Beneficiary, msg.OpeningPrice,
		msg.BuyoutPrice, msg.TimeToSell); err != nil {
		return nil, fmt.Errorf("failed to PutNFTOnAuction: %v", err)
	}

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.Beneficiary.String()),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
			sdk.NewAttribute(types.AttributeKeyOpeningPrice, msg.OpeningPrice.String()),
			sdk.NewAttribute(types.AttributeKeyBuyoutPrice, msg.BuyoutPrice.String()),
			sdk.NewAttribute(types.AttributeKeyFinishTime, msg.TimeToSell.String()),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	k.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgPutNFTOnAuction)
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgRemoveNFTFromAuction(ctx sdk.Context, k *Keeper, msg MsgRemoveNFTFromAuction) (*sdk.Result, error) {
	k.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgRemoveNFTFromAuction)

	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	nft, err := k.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	if !nft.Owner.Equals(msg.Owner) {
		return nil, fmt.Errorf("auction lot owner: %v and finisher: %v do not match", msg.Owner, nft.Owner)
	}

	// return bid to last bidder if exists
	if lot.LastBid != nil {
		_, err := k.coinKeeper.AddCoins(ctx, lot.LastBid.Bidder, lot.LastBid.Bid)
		if err != nil {
			return nil, err
		}
	}

	// return nft to owner, delete lot
	if err := k.RemoveNFTFromAuction(ctx, msg.TokenID, msg.Owner); err != nil {
		return nil, err
	}

	k.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgRemoveNFTFromAuction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgFinishAuction(ctx sdk.Context, k *Keeper, msg MsgFinishAuction) (*sdk.Result, error) {
	k.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgFinishAuction)
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	nft, err := k.GetNFT(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	// auction time has not expired yet
	if lot.ExpirationTime.After(time.Now().UTC()) {
		if !nft.Owner.Equals(msg.Owner) {
			return nil, fmt.Errorf("auction lot owner: %v and finisher: %v do not match", msg.Owner, nft.Owner)
		}
	}

	owner := nft.Owner.String()

	// no bids on lot
	if lot.LastBid != nil {
		owner = lot.LastBid.Bidder.String()
		if err := k.BuyLotOnAuction(ctx, lot.LastBid.Bidder, lot.LastBid.BuyerBeneficiary,
			lot.LastBid.Bid, lot, lot.LastBid.BeneficiaryCommission); err != nil {
			return nil, err
		}
	} else {
		if err := k.RemoveNFTFromAuction(ctx, msg.TokenID, nft.Owner); err != nil {
			return nil, err
		}
	}

	k.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgFinishAuction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyOwner, owner),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgMakeBidOnAuction(ctx sdk.Context, k *Keeper, msg MsgMakeBidOnAuction) (*sdk.Result, error) {
	k.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgMakeBidOnAuction)

	logger := ctx.Logger()
	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("auction is already finished")
	}

	beneficiariesCommission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(msg.BeneficiaryCommission, 64)
	if err == nil {
		beneficiariesCommission = parsed
	}
	if beneficiariesCommission > k.config.MaximumBeneficiaryCommission {
		return nil, fmt.Errorf("failed to BuyNFT: beneficiary commission is too high")
	}
	beneficiariesCommissionString := fmt.Sprintf("%v", beneficiariesCommission)

	// bid is less than lastBid
	if lot.LastBid != nil {
		if msg.Bid.IsAllLTE(lot.LastBid.Bid) {
			return nil, fmt.Errorf("bid: %+v is lower than last bid: %+v", msg.Bid, lot.LastBid.Bid)
		}
	}

	// bid is less than opening price
	if lot.OpeningPrice.IsAnyGT(msg.Bid) {
		return nil, fmt.Errorf("bid: %+v is lower than opening price: %+v", msg.Bid, lot.OpeningPrice)
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
			return nil, err
		}
	}

	// take coins from new bidder
	_, err = k.coinKeeper.SubtractCoins(ctx, msg.Bidder, msg.Bid)
	if err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return nil, err
	}

	auctionBid := types.NewAuctionBid(msg.Bidder, msg.BuyerBeneficiary, msg.Bid, beneficiariesCommissionString)
	lot.SetLastBid(auctionBid)

	if err := k.UpdateAuctionLot(ctx, lot); err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return nil, err
	}

	attrs := []sdk.Attribute{
		sdk.NewAttribute(types.AttributeKeyBidder, msg.Bidder.String()),
		sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.BuyerBeneficiary.String()),
		sdk.NewAttribute(types.AttributeKeyBid, msg.Bid.String()),
		sdk.NewAttribute(types.AttributeKeyCommission, msg.BeneficiaryCommission),
		sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
	}

	// Bid is more than buyout price. Perform buyout.
	if !lot.BuyoutPrice.IsZero() {
		if msg.Bid.IsAllGTE(lot.BuyoutPrice) {
			err = k.BuyLotOnAuction(ctx, msg.Bidder, msg.BuyerBeneficiary, lot.BuyoutPrice, lot, msg.BeneficiaryCommission)
			if err != nil {
				return nil, err
			}
			attrs = append(attrs, sdk.NewAttribute(types.AttributeKeyIsBuyout, "true"))
		}
	}

	k.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgMakeBidOnAuction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			attrs...,
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Bidder.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}

func handleMsgBuyoutOnAuction(ctx sdk.Context, k *Keeper, msg MsgBuyoutOnAuction) (*sdk.Result, error) {
	k.increaseCounter(common.PrometheusValueReceived, common.PrometheusValueMsgBuyoutFromAuction)

	lot, err := k.GetAuctionLot(ctx, msg.TokenID)
	if err != nil {
		return nil, err
	}

	if lot.ExpirationTime.Before(time.Now().UTC()) {
		return nil, fmt.Errorf("auction is already finished")
	}

	if lot.BuyoutPrice.IsZero() {
		return nil, fmt.Errorf("lot has no buyoutprice")
	}

	err = k.BuyLotOnAuction(ctx, msg.Buyer, msg.BuyerBeneficiary, lot.BuyoutPrice, lot, msg.BeneficiaryCommission)
	if err != nil {
		return nil, err
	}

	k.increaseCounter(common.PrometheusValueAccepted, common.PrometheusValueMsgBuyoutFromAuction)

	ctx.EventManager().EmitEvents(sdk.Events{
		sdk.NewEvent(
			msg.Type(),
			sdk.NewAttribute(types.AttributeKeyBidder, msg.Buyer.String()),
			sdk.NewAttribute(types.AttributeKeyBeneficiary, msg.BuyerBeneficiary.String()),
			sdk.NewAttribute(types.AttributeKeyCommission, msg.BeneficiaryCommission),
			sdk.NewAttribute(types.AttributeKeyNFTID, msg.TokenID),
		),
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Buyer.String()),
		),
	})
	return &sdk.Result{Events: ctx.EventManager().ABCIEvents()}, nil
}
