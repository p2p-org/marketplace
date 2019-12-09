package marketplace

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	"github.com/cosmos/cosmos-sdk/x/distribution"
	"github.com/cosmos/cosmos-sdk/x/staking"
	"github.com/cosmos/cosmos-sdk/x/supply"
	"github.com/cosmos/modules/incubator/nft"
	"github.com/corestario/marketplace/common"
	"github.com/corestario/marketplace/x/marketplace/config"
	"github.com/corestario/marketplace/x/marketplace/types"
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
	supplyKeeper             *supply.Keeper
	accKeeper                *auth.AccountKeeper
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
	supplyKeeper *supply.Keeper,
	accKeeper *auth.AccountKeeper,
) *Keeper {
	return &Keeper{
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
		supplyKeeper:             supplyKeeper,
		accKeeper:                accKeeper,
	}
}

func (k *Keeper) increaseCounter(labels ...string) {
	counter, err := k.msgMetr.NumMsgs.GetMetricWithLabelValues(labels...)
	if err != nil {
		pl.Errorf("get metrics with label values error: %v", err)
		return
	}
	counter.Inc()
}

func (k *Keeper) GetNFT(ctx sdk.Context, id string) (*NFT, error) {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return nil, fmt.Errorf("could not find NFT with id %s", id)
	}

	bz := store.Get([]byte(id))
	var token NFT
	k.cdc.MustUnmarshalJSON(bz, &token)

	return &token, nil
}

func (k *Keeper) MintNFT(ctx sdk.Context, nft *NFT) error {
	id := nft.ID
	store := ctx.KVStore(k.storeKey)
	if store.Has([]byte(id)) {
		return fmt.Errorf("nft with ID %s already exists", id)
	}

	bz := k.cdc.MustMarshalJSON(nft)
	store.Set([]byte(id), bz)
	return nil
}

func (k *Keeper) BurnNFT(ctx sdk.Context, id string) error {
	store := ctx.KVStore(k.storeKey)
	if !store.Has([]byte(id)) {
		return fmt.Errorf("could not find NFT with id %s", id)
	}

	store.Delete([]byte(id))
	return nil
}

// Get an iterator over all NFTs.
func (k *Keeper) GetNFTsIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.storeKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

// Get an iterator over all registered currencies
func (k *Keeper) GetRegisteredCurrenciesIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k *Keeper) PutNFTOnMarket(ctx sdk.Context, id string, owner, beneficiary sdk.AccAddress, price sdk.Coins) error {
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

func (k *Keeper) RemoveNFTFromMarket(ctx sdk.Context, id string, owner sdk.AccAddress) error {
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

func (k *Keeper) UpdateNFT(ctx sdk.Context, newToken *NFT) error {
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

func (k *Keeper) IsDenomExist(ctx sdk.Context, coins sdk.Coins) bool {
	if coins.Empty() {
		return false
	}
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	for _, v := range coins {
		v := v
		if !store.Has([]byte(v.Denom)) {
			return false
		}
	}

	return true
}

// Creates a new fungible token with given supply and denom for FungibleTokenCreationPrice
func (k *Keeper) CreateFungibleToken(ctx sdk.Context, creator sdk.AccAddress, denom string, amount int64) error {
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
	mintedCoins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))
	sdkErr := k.supplyKeeper.MintCoins(ctx, bank.ModuleName, mintedCoins)
	if sdkErr != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return fmt.Errorf("failed to mint fungible tokens: %v", sdkErr.Error())
	}

	sdkErr = k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, bank.ModuleName, creator, mintedCoins)
	if sdkErr != nil {
		RollbackCommissions(ctx, k, logger, initialBalances)
		return fmt.Errorf("failed to add coins: %v", sdkErr.Error())
	}

	k.registerFungibleTokensCurrency(ctx, FungibleToken{Creator: creator, Denom: denom, EmissionAmount: amount})

	return nil
}

// Should be run just once
func (k *Keeper) RegisterBasicDenoms(ctx sdk.Context) {
	ft := FungibleToken{Creator: []byte{}, Denom: types.DefaultTokenDenom, EmissionAmount: 1}
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	k.supplyKeeper.SetSupply(ctx, supply.NewSupply(sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.OneInt()))))

	store.Set([]byte(ft.Denom), k.cdc.MustMarshalJSON(ft))
}

// Registers fungible token for prevent double creation
func (k *Keeper) registerFungibleTokensCurrency(ctx sdk.Context, ft FungibleToken) {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	store.Set([]byte(ft.Denom), k.cdc.MustMarshalJSON(ft))
}

// Transfers amount of fungible tokens from one account to another
func (k *Keeper) TransferFungibleTokens(ctx sdk.Context, currencyOwner, recipient sdk.AccAddress, denom string, amount int64) error {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if !store.Has([]byte(denom)) {
		return fmt.Errorf("unknown currency")
	}
	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, currencyOwner, bank.ModuleName, coins); err != nil {
		return fmt.Errorf("failed to send tokens to module")
	}
	if err := k.supplyKeeper.SendCoinsFromModuleToAccount(ctx, bank.ModuleName, recipient, coins); err != nil {
		return fmt.Errorf("failed to send tokens to account")
	}
	return nil
}

func (k *Keeper) GetFungibleToken(ctx sdk.Context, name string) (*FungibleToken, error) {
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
func (k *Keeper) GetFungibleTokensIterator(ctx sdk.Context) sdk.Iterator {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	return sdk.KVStorePrefixIterator(store, nil)
}

func (k *Keeper) BurnFungibleTokens(ctx sdk.Context, currencyOwner sdk.AccAddress, denom string, amount int64) error {
	store := ctx.KVStore(k.currencyRegistryStoreKey)
	if !store.Has([]byte(denom)) {
		return fmt.Errorf("unknown currency")
	}

	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount)))
	if err := k.supplyKeeper.SendCoinsFromAccountToModule(ctx, currencyOwner, bank.ModuleName, coins); err != nil {
		return fmt.Errorf("failed to send tokens to module")
	}

	err := k.supplyKeeper.BurnCoins(ctx, bank.ModuleName, sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(amount))))
	if err != nil {
		return fmt.Errorf("failed to burn fungible tokens")
	}
	return nil
}

func (k *Keeper) TransferNFT(ctx sdk.Context, id string, sender, recipient sdk.AccAddress) error {
	token, err := k.GetNFT(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to GetNFT: %v", err)
	}

	if token.IsOnSale() {
		return fmt.Errorf("failed to transferNFT: NFT is on sale")
	}

	if !token.Owner.Equals(sender) {
		return fmt.Errorf("%s is not the owner of NFT #%s", sender.String(), id)
	}
	token.Owner = recipient

	return k.UpdateNFT(ctx, token)
}
