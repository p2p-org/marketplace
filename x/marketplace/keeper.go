package marketplace

import (
	"fmt"

	"github.com/dgamingfoundation/marketplace/x/marketplace/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/dgamingfoundation/marketplace/x/marketplace/config"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper               bank.Keeper
	stakingKeeper            staking.Keeper
	distrKeeper              distribution.Keeper
	storeKey                 sdk.StoreKey // Unexposed key to access store from sdk.Context
	currencyRegistryStoreKey *sdk.KVStoreKey
	cdc                      *codec.Codec // The wire codec for binary encoding/decoding.
	config                   *config.MPServerConfig
}

// NewKeeper creates new instances of the marketplace Keeper
func NewKeeper(
	coinKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	distrKeeper distribution.Keeper,
	storeKey sdk.StoreKey,
	currencyRegistryStoreKey *sdk.KVStoreKey,
	cdc *codec.Codec,
	cfg *config.MPServerConfig,
) Keeper {
	return Keeper{
		coinKeeper:               coinKeeper,
		stakingKeeper:            stakingKeeper,
		distrKeeper:              distrKeeper,
		storeKey:                 storeKey,
		currencyRegistryStoreKey: currencyRegistryStoreKey,
		cdc:                      cdc,
		config:                   cfg,
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

// Get an iterator over all registered currencies
func (k Keeper) GetRegisteredCurrenciesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
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
	nft.BaseNFT = nft.SetOwner(recipient)

	return k.UpdateNFT(ctx, nft)
}

func (k Keeper) PutNFTOnMarket(ctx sdk.Context, id string, owner, beneficiary sdk.AccAddress, price sdk.Coins) error {
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

func (k Keeper) IsDenomExist(ctx sdk.Context, coins sdk.Coins) bool {
	if coins.Empty() {
		return false
	}
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	for i := 0; i < coins.Len(); i++ {
		if !store.Has([]byte(coins[i].Denom)) {
			return false
		}
	}

	return true
}

// Creates a new fungible token with given supply and denom for FungibleTokenCreationPrice
func (k Keeper) CreateFungibleToken(ctx sdk.Context, creator sdk.AccAddress, denom string, amount int64) error {
	logger := ctx.Logger()

	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if store.Has([]byte(denom)) {
		return fmt.Errorf("currency already exists")
	}

	commissionAddress, err := sdk.AccAddressFromBech32(FungibleCommissionAddress)
	if err != nil {
		return fmt.Errorf("failed to get comissionAddress: %v", err)
	}

	initialBalances := GetBalances(ctx, k, creator, commissionAddress)

	if err := k.coinKeeper.SendCoins(ctx, creator, commissionAddress,
		sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.NewInt(FungibleTokenCreationPrice)))); err != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return fmt.Errorf("failed to send coins to comissionAddress")
	}

	if _, err := k.coinKeeper.AddCoins(ctx, creator, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))); err != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return fmt.Errorf("failed to add coins: %v", err)
	}
	k.registerFungibleTokensCurrency(ctx, FungibleToken{Creator: creator, Denom: denom, EmissionAmount: amount})
	return nil
}

// Should be run just once
func (k *Keeper) RegisterBasicDenoms(ctx sdk.Context) {
	ft := FungibleToken{Creator: []byte{}, Denom: types.DefaultTokenDenom, EmissionAmount: 1}
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	store.Set([]byte(ft.Denom), k.cdc.MustMarshalJSON(ft))
}

// Registers fungible token for prevent double creation
func (k Keeper) registerFungibleTokensCurrency(ctx sdk.Context, ft FungibleToken) {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	store.Set([]byte(ft.Denom), k.cdc.MustMarshalJSON(ft))
}

// Transfers amount of fungible tokens from one account to another
func (k Keeper) TransferFungibleTokens(ctx sdk.Context, currencyOwner, recipient sdk.AccAddress, denom string, amount int64) error {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if !store.Has([]byte(denom)) {
		return fmt.Errorf("unknown currency")
	}

	if err := k.coinKeeper.SendCoins(ctx, currencyOwner, recipient, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))); err != nil {
		return fmt.Errorf("failed to transfer tokens")
	}
	return nil
}

func (k Keeper) GetFungibleToken(ctx sdk.Context, name string) (*FungibleToken, error) {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if !store.Has([]byte(name)) {
		return nil, fmt.Errorf("could not find Fungible Token with name %s", name)
	}

	bz := store.Get([]byte(name))
	var ft FungibleToken
	k.cdc.MustUnmarshalJSON(bz, &ft)

	return &ft, nil
}

// Get an iterator over all Fungible Tokens
func (k Keeper) GetFungibleTokensIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}
