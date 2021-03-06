package marketplace

import (
	"fmt"
	"strconv"
	"time"

	"github.com/corestario/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func (k *Keeper) PutNFTOnAuction(ctx sdk.Context, id string, owner, beneficiary sdk.AccAddress,
	openingPrice, buyoutPrice sdk.Coins, expirationTime time.Time) error {

	token, err := k.GetNFT(ctx, id)

	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if !token.Owner.Equals(owner) {
		return fmt.Errorf("%s is not the owner of NFT #%s", owner.String(), id)
	}

	if token.IsOnSale() {
		return fmt.Errorf("NFT #%s is alredy on sale", id)
	}
	token.SetStatus(types.NFTStatusOnAuction)
	token.SetSellerBeneficiary(beneficiary)
	lot := types.NewAuctionLot(id, openingPrice, buyoutPrice, expirationTime)
	err = k.createAuctionLot(ctx, lot)
	if err != nil {
		return fmt.Errorf("failed to create auction lot: %v", err)
	}
	return k.UpdateNFT(ctx, token)
}

func (k *Keeper) createAuctionLot(ctx sdk.Context, lot *types.AuctionLot) error {
	store := ctx.KVStore(k.auctionStoreKey)
	if store.Has([]byte(lot.NFTID)) {
		return fmt.Errorf("lot already exists")
	}
	bz := k.cdc.MustMarshalJSON(lot)
	store.Set([]byte(lot.NFTID), bz)
	return nil
}

// important! must do all necessary checks before evoking this function
func (k *Keeper) RemoveNFTFromAuction(ctx sdk.Context, id string, owner sdk.AccAddress) error {
	nft, err := k.GetNFT(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if !nft.Owner.Equals(owner) {
		return fmt.Errorf("%s is not the owner of NFT #%s", owner.String(), id)
	}

	return k.removeNFTFromAuction(ctx, nft)
}

func (k *Keeper) removeNFTFromAuction(ctx sdk.Context, nft *NFT) error {
	if nft.Status != types.NFTStatusOnAuction {
		return fmt.Errorf("NFT #%s is not on auction", nft.ID)
	}

	err := k.deleteAuctionLot(ctx, nft.ID)
	if err != nil {
		return err
	}

	nft.SetStatus(types.NFTStatusDefault)
	nft.SetSellerBeneficiary(sdk.AccAddress{})

	return k.UpdateNFT(ctx, nft)
}

func (k *Keeper) deleteAuctionLot(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.auctionStoreKey)
	if !store.Has([]byte(id)) {
		return fmt.Errorf("lot does not exist")
	}
	store.Delete([]byte(id))
	return nil
}

func (k *Keeper) UpdateAuctionLot(ctx sdk.Context, lot *types.AuctionLot) error {
	store := ctx.KVStore(k.auctionStoreKey)
	if !store.Has([]byte(lot.NFTID)) {
		return fmt.Errorf("could not find lot with id %s", lot.NFTID)
	}

	bz := k.cdc.MustMarshalJSON(lot)
	store.Set([]byte(lot.NFTID), bz)
	return nil
}

func (k *Keeper) GetAuctionLot(ctx sdk.Context, id string) (*types.AuctionLot, error) {
	store := ctx.KVStore(k.auctionStoreKey)
	if !store.Has([]byte(id)) {
		return nil, fmt.Errorf("lot does not exist")
	}
	bz := store.Get([]byte(id))
	var lot types.AuctionLot
	k.cdc.MustUnmarshalJSON(bz, &lot)
	return &lot, nil
}

func (k *Keeper) GetAuctionLotsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.auctionStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// buyout the lot
func (k *Keeper) BuyLotOnAuction(ctx sdk.Context, buyer, buyerBeneficiary sdk.AccAddress,
	price sdk.Coins, lot *types.AuctionLot, buyerCommission string) error {
	logger := ctx.Logger()
	nft, err := k.GetNFT(ctx, lot.NFTID)
	if err != nil {
		return err
	}
	if !nft.IsOnAuction() {
		return fmt.Errorf("nft is not on auction")
	}

	commission := types.DefaultBeneficiariesCommission
	parsed, err := strconv.ParseFloat(buyerCommission, 64)
	if err == nil {
		commission = parsed
	}
	balances := GetBalances(ctx, k, buyer, buyerBeneficiary, nft.SellerBeneficiary, nft.Owner)
	if lot.LastBid != nil {
		balances = append(balances,
			GetBalances(ctx, k, lot.LastBid.Bidder, lot.LastBid.BuyerBeneficiary)...)
		_, err = k.coinKeeper.AddCoins(ctx, lot.LastBid.Bidder, lot.LastBid.Bid)
		if err != nil {
			RollbackCommissions(ctx, k, logger, balances)
			return err
		}
	}

	// similar to buyNFTOnMarket
	priceAfterCommission, err := doNFTCommissions(
		ctx,
		k,
		buyer,
		nft.Owner,
		nft.SellerBeneficiary,
		buyerBeneficiary,
		price,
		commission,
	)
	if err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return err
	}

	err = k.coinKeeper.SendCoins(ctx, buyer, nft.Owner, priceAfterCommission)
	if err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return fmt.Errorf("buyer does not have enough coins")
	}

	err = k.deleteAuctionLot(ctx, lot.NFTID)
	if err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return err
	}

	// transfer nfr to new owner
	nft.SetSellerBeneficiary(sdk.AccAddress{})
	nft.Owner = buyer
	nft.SetStatus(types.NFTStatusDefault)

	if err := k.UpdateNFT(ctx, nft); err != nil {
		RollbackCommissions(ctx, k, logger, balances)
		return err
	}
	return nil
}

func (k *Keeper) CheckFinishedAuctions(ctx sdk.Context) {
	// TODO: is error handler necessary here?
	logger := ctx.Logger()
	iterator := k.GetAuctionLotsIterator(ctx)
	timeNow := time.Now().UTC()
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lot types.AuctionLot
		k.cdc.MustUnmarshalJSON(iterator.Value(), &lot)
		if lot.ExpirationTime.Before(timeNow) {
			// there was at least one bid
			addr, err := sdk.AccAddressFromBech32(k.config.FinishingAccountAddr)
			if err != nil {
				logger.Error("failed to get acc address from bench32", "bench", k.config.FinishingAccountAddr, "error", err)
				continue
			}

			acc := k.accKeeper.GetAccount(ctx, addr)
			if err := k.SendFinish(lot.NFTID, acc); err != nil {
				logger.Error("failed to sent finish tx", "lot", lot.NFTID, "error", err)
				continue
			}
		}
	}
}
