package marketplace

import (
	"fmt"

	"github.com/dgamingfoundation/marketplace/x/marketplace/config"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper    bank.Keeper
	stakingKeeper staking.Keeper
	distrKeeper   distribution.Keeper
	storeKey      sdk.StoreKey // Unexposed key to access store from sdk.Context
	cdc           *codec.Codec // The wire codec for binary encoding/decoding.
	config        *config.MPServerConfig
}

// NewKeeper creates new instances of the marketplace Keeper
func NewKeeper(
	coinKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	distrKeeper distribution.Keeper,
	storeKey sdk.StoreKey,
	cdc *codec.Codec,
	cfg *config.MPServerConfig,
) Keeper {
	return Keeper{
		coinKeeper:    coinKeeper,
		stakingKeeper: stakingKeeper,
		distrKeeper:   distrKeeper,
		storeKey:      storeKey,
		cdc:           cdc,
		config:        cfg,
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
	id := nft.GetID()
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

func (k Keeper) TransferNFT(ctx sdk.Context, id string, sender, recipient sdk.AccAddress) error {
	nft, err := k.GetNFT(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if !nft.GetOwner().Equals(sender) {
		return fmt.Errorf("%s is not the owner of NFT #%s", sender.String(), id)
	}
	nft.SetOwner(recipient)

	return k.UpdateNFT(ctx, nft)
}

func (k Keeper) SellNFT(ctx sdk.Context, id string, owner, beneficiary sdk.AccAddress, price sdk.Coins) error {
	nft, err := k.GetNFT(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if !nft.GetOwner().Equals(owner) {
		return fmt.Errorf("%s is not the owner of NFT #%s", owner.String(), id)
	}
	nft.SetPrice(price)
	nft.SetOnSale(true)
	nft.SetSellerBeneficiary(beneficiary)

	return k.UpdateNFT(ctx, nft)
}

func (k Keeper) UpdateNFT(ctx sdk.Context, newToken *NFT) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(newToken.GetID())) {
		return fmt.Errorf("could not find NFT with id %s", newToken.GetID())
	}

	bz := k.cdc.MustMarshalJSON(newToken)
	store.Set([]byte(newToken.GetID()), bz)
	return nil
}
