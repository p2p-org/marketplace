package marketplace

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/x/bank"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper bank.Keeper
	storeKey   sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc        *codec.Codec // The wire codec for binary encoding/decoding.
}

// NewKeeper creates new instances of the marketplace Keeper
func NewKeeper(coinKeeper bank.Keeper, storeKey sdk.StoreKey, cdc *codec.Codec) Keeper {
	return Keeper{
		coinKeeper: coinKeeper,
		storeKey:   storeKey,
		cdc:        cdc,
	}
}

func (k Keeper) GetNFT(ctx sdk.Context, id string) (*NFT, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return nil, fmt.Errorf("could not find NFT with id %s", id)
	}

	bz := store.Get([]byte(id))
	var nft NFT
	k.cdc.MustUnmarshalJSON(bz, &nft)

	return &nft, nil
}

func (k Keeper) MintNFT(ctx sdk.Context, nft *NFT) error {
	id := nft.NFT.GetID()
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(id)) {
		return fmt.Errorf("nft with ID %s already exists", id)
	}

	bz := k.cdc.MustMarshalJSON(nft)
	store.Set([]byte(id), bz)
	return nil
}

// Get an iterator over all NFTs.
func (k Keeper) GetNFTsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
