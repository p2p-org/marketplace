package marketplace

import (
	"fmt"

	"github.com/dgamingfoundation/cosmos-sdk/x/nft"

	"github.com/dgamingfoundation/marketplace/common"

	"github.com/dgamingfoundation/marketplace/x/marketplace/types"

	"github.com/dgamingfoundation/cosmos-sdk/codec"
	sdk "github.com/dgamingfoundation/cosmos-sdk/types"
	"github.com/dgamingfoundation/cosmos-sdk/x/bank"
	"github.com/dgamingfoundation/cosmos-sdk/x/distribution"
	"github.com/dgamingfoundation/cosmos-sdk/x/staking"
	"github.com/dgamingfoundation/marketplace/x/marketplace/config"
	pl "github.com/prometheus/common/log"
)

// Keeper maintains the link to data storage and exposes getter/setter methods for the various parts of the state machine
type Keeper struct {
	coinKeeper               bank.Keeper
	stakingKeeper            staking.Keeper
	distrKeeper              distribution.Keeper
	storeKey                 sdk.StoreKey // Unexposed key to access store from sdk.Context
	currencyRegistryStoreKey *sdk.KVStoreKey
	auctionStoreKey          *sdk.KVStoreKey
	cdc                      *codec.Codec // The wire codec for binary encoding/decoding.
	config                   *config.MPServerConfig
	msgMetr                  *common.MsgMetrics
	nftKeeper                *nft.Keeper
}

// NewKeeper creates new instances of the marketplace Keeper
func NewKeeper(
	coinKeeper bank.Keeper,
	stakingKeeper staking.Keeper,
	distrKeeper distribution.Keeper,
	storeKey sdk.StoreKey,
	currencyRegistryStoreKey *sdk.KVStoreKey,
	auctionStoreKey *sdk.KVStoreKey,
	cdc *codec.Codec,
	cfg *config.MPServerConfig,
	msgMetr *common.MsgMetrics,
	nftKeeper *nft.Keeper,
) Keeper {
	return Keeper{
		coinKeeper:               coinKeeper,
		stakingKeeper:            stakingKeeper,
		distrKeeper:              distrKeeper,
		storeKey:                 storeKey,
		currencyRegistryStoreKey: currencyRegistryStoreKey,
		auctionStoreKey:          auctionStoreKey,
		cdc:                      cdc,
		config:                   cfg,
		msgMetr:                  msgMetr,
		nftKeeper:                nftKeeper,
	}
}

func (k Keeper) increaseCounter(labels ...string) {
	counter, err := k.msgMetr.NumMsgs.GetMetricWithLabelValues(labels...)
	if err != nil {
		pl.Errorf("get metrics with label values error: %v", err)
		return
	}
	counter.Inc()
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
	id := nft.ID
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

func (k Keeper) PutNFTOnMarket(ctx sdk.Context, id string, owner, beneficiary sdk.AccAddress, price sdk.Coins) error {
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
	token.SetPrice(price)
	token.SetStatus(types.NFTStatusOnMarket)
	token.SetSellerBeneficiary(beneficiary)

	return k.UpdateNFT(ctx, token)
}

func (k Keeper) RemoveNFTFromMarket(ctx sdk.Context, id string, owner sdk.AccAddress) error {
	token, err := k.GetNFT(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if !token.Owner.Equals(owner) {
		return fmt.Errorf("%s is not the owner of NFT #%s", owner.String(), id)
	}

	if !token.IsOnMarket() {
		return fmt.Errorf("NFT #%s is not on market", id)
	}
	token.SetPrice(sdk.Coins{})
	token.SetStatus(types.NFTStatusDefault)
	token.SetSellerBeneficiary(sdk.AccAddress{})

	return k.UpdateNFT(ctx, token)
}

func (k Keeper) UpdateNFT(ctx sdk.Context, newToken *NFT) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(newToken.ID)) {
		return fmt.Errorf("could not find NFT with id %s", newToken.ID)
	}

	bz := k.cdc.MustMarshalJSON(newToken)
	store.Set([]byte(newToken.ID), bz)

	newBaseToken, err := k.nftKeeper.GetNFT(ctx, newToken.Denom, newToken.ID)
	if err != nil {
		return fmt.Errorf("failed to get base token: %v", err)
	}
	newBaseToken.SetOwner(newToken.Owner)
	if err := k.nftKeeper.UpdateNFT(ctx, newToken.Denom, newBaseToken); err != nil {
		return fmt.Errorf("failed to update base token: %v", err)
	}

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

func (k Keeper) BurnFungibleTokens(ctx sdk.Context, currencyOwner sdk.AccAddress, denom string, amount int64) error {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if !store.Has([]byte(denom)) {
		return fmt.Errorf("unknown currency")
	}

	_, err := k.coinKeeper.SubtractCoins(ctx, currencyOwner, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount))))
	if err != nil {
		return fmt.Errorf("failed to burn fungible tokens")
	}
	return nil
}
