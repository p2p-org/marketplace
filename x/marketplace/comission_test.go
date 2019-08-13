package marketplace_test

import (
	"github.com/dgamingfoundation/marketplace/common"
	"io/ioutil"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	distr "github.com/cosmos/cosmos-sdk/x/distribution"
	distrTypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
	"github.com/cosmos/cosmos-sdk/x/mock"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/cosmos/cosmos-sdk/x/slashing"
	"github.com/cosmos/cosmos-sdk/x/staking"
	stakingTypes "github.com/cosmos/cosmos-sdk/x/staking/types"
	app "github.com/dgamingfoundation/marketplace"
	"github.com/dgamingfoundation/marketplace/x/marketplace"
	"github.com/dgamingfoundation/marketplace/x/marketplace/config"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmTypes "github.com/tendermint/tendermint/types"
)

type marketplaceKeeperTest struct {
	ctx sdk.Context

	accountKeeper       auth.AccountKeeper
	bankKeeper          bank.Keeper
	stakingKeeper       staking.Keeper
	feeCollectionKeeper auth.FeeCollectionKeeper
	distrKeeper         distr.Keeper
	slashingKeeper      slashing.Keeper
	ms                  store.CommitMultiStore
	marketKeeper        marketplace.Keeper
	dbDir               string
	addrs               []sdk.AccAddress
}

// clear removes temp dirs
func (mp *marketplaceKeeperTest) clear() error {
	return os.RemoveAll(mp.dbDir)
}

func createMarketplaceKeeperTest() (*marketplaceKeeperTest, error) {
	var (
		err error
	)

	mpKeeperTest := new(marketplaceKeeperTest)

	cdc := app.MakeCodec()

	mpKeeperTest.dbDir, err = ioutil.TempDir("", "goleveldb-app-sim")
	if err != nil {
		return nil, err
	}
	db, err := sdk.NewLevelDB("Simulation", mpKeeperTest.dbDir)
	if err != nil {
		return nil, err
	}

	keyParams := sdk.NewKVStoreKey(params.StoreKey)
	tkeyParams := sdk.NewTransientStoreKey(params.TStoreKey)
	keyAccount := sdk.NewKVStoreKey(auth.StoreKey)
	keyFeeCollection := sdk.NewKVStoreKey(auth.FeeStoreKey)
	keyStaking := sdk.NewKVStoreKey(staking.StoreKey)
	tkeyStaking := sdk.NewTransientStoreKey(staking.TStoreKey)
	keySlashing := sdk.NewKVStoreKey(slashing.StoreKey)
	keyDistr := sdk.NewKVStoreKey(distr.StoreKey)
	keyRegisterCurrency := sdk.NewKVStoreKey("register_currency")

	paramsKeeper := params.NewKeeper(cdc, keyParams, tkeyParams, params.DefaultCodespace)

	authSubspace := paramsKeeper.Subspace(auth.DefaultParamspace)
	bankSupspace := paramsKeeper.Subspace(bank.DefaultParamspace)
	stakingSubspace := paramsKeeper.Subspace(staking.DefaultParamspace)
	distrSubspace := paramsKeeper.Subspace(distr.DefaultParamspace)
	slashingSubspace := paramsKeeper.Subspace(slashing.DefaultParamspace)

	mpKeeperTest.accountKeeper = auth.NewAccountKeeper(
		cdc,
		keyAccount,
		authSubspace,
		auth.ProtoBaseAccount,
	)
	mpKeeperTest.bankKeeper = bank.NewBaseKeeper(
		mpKeeperTest.accountKeeper,
		bankSupspace,
		bank.DefaultCodespace,
	)
	mpKeeperTest.stakingKeeper = staking.NewKeeper(
		cdc,
		keyStaking,
		tkeyStaking,
		mpKeeperTest.bankKeeper,
		stakingSubspace,
		staking.DefaultCodespace,
	)
	mpKeeperTest.feeCollectionKeeper = auth.NewFeeCollectionKeeper(cdc, keyFeeCollection)
	mpKeeperTest.distrKeeper = distr.NewKeeper(
		cdc,
		keyDistr,
		distrSubspace,
		mpKeeperTest.bankKeeper,
		&mpKeeperTest.stakingKeeper,
		mpKeeperTest.feeCollectionKeeper,
		distr.DefaultCodespace,
	)
	mpKeeperTest.slashingKeeper = slashing.NewKeeper(
		cdc,
		keySlashing,
		&mpKeeperTest.stakingKeeper,
		slashingSubspace,
		slashing.DefaultCodespace,
	)
	mpKeeperTest.stakingKeeper = *mpKeeperTest.stakingKeeper.SetHooks(
		staking.NewMultiStakingHooks(
			mpKeeperTest.distrKeeper.Hooks(),
			mpKeeperTest.slashingKeeper.Hooks()),
	)

	mpStore := sdk.NewKVStoreKey(marketplace.StoreKey)
	mpKeeperTest.ms = store.NewCommitMultiStore(db)
	mpKeeperTest.ms.MountStoreWithDB(mpStore, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keyAccount, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keyFeeCollection, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keySlashing, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keyStaking, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keyDistr, sdk.StoreTypeIAVL, db)
	mpKeeperTest.ms.MountStoreWithDB(keyRegisterCurrency, sdk.StoreTypeIAVL, db)
	if err := mpKeeperTest.ms.LoadLatestVersion(); err != nil {
		return nil, err
	}

	mpKeeperTest.marketKeeper = marketplace.NewKeeper(mpKeeperTest.bankKeeper, mpKeeperTest.stakingKeeper,
		mpKeeperTest.distrKeeper, mpStore, keyRegisterCurrency, cdc, config.DefaultMPServerConfig(), common.NewPrometheusMsgMetrics("marketplace"))

	mpKeeperTest.ctx = sdk.NewContext(mpKeeperTest.ms, abci.Header{}, false, log.NewNopLogger())
	mpKeeperTest.marketKeeper.RegisterBasicDenoms(mpKeeperTest.ctx)
	return mpKeeperTest, nil
}

func (mp *marketplaceKeeperTest) updateAccountsWithCoins(coins sdk.Coins) error {
	_, mp.addrs, _, _ = mock.CreateGenAccounts(4, coins)

	for _, addr := range mp.addrs {
		if _, err := mp.bankKeeper.AddCoins(mp.ctx, addr, coins); err != nil {
			return err
		}
	}
	return nil
}

func (mp *marketplaceKeeperTest) updateVoteInfos(validatorsCount int, coins sdk.Coins) ([]abci.VoteInfo, error) {
	voteInfos := make([]abci.VoteInfo, 0, validatorsCount)
	for i := 0; i < validatorsCount; i++ {
		pv := tmTypes.NewMockPV()
		pubKey := pv.GetPubKey()
		voteInfo := abci.NewPopulatedVoteInfo(rand.New(rand.NewSource(time.Now().UnixNano())), true)
		voteInfo.Validator.Address = pubKey.Address()
		voteInfo.SignedLastBlock = true
		if _, err := mp.bankKeeper.AddCoins(mp.ctx, voteInfo.Validator.Address, coins); err != nil {
			return nil, err
		}
		val := stakingTypes.NewValidator(pubKey.Address().Bytes(), pubKey, stakingTypes.Description{})
		mp.stakingKeeper.SetValidator(mp.ctx, val)
		mp.stakingKeeper.SetValidatorByConsAddr(mp.ctx, val)
		mp.distrKeeper.SetValidatorCurrentRewards(mp.ctx, val.GetOperator(),
			distrTypes.NewValidatorCurrentRewards(sdk.NewDecCoins(coins), 10))
		mp.distrKeeper.SetValidatorOutstandingRewards(mp.ctx, val.GetOperator(), sdk.NewDecCoins(coins))
		voteInfos = append(voteInfos, *voteInfo)
	}
	mp.ctx = mp.ctx.WithVoteInfos(voteInfos)
	return voteInfos, nil
}

type testBuyPutOnMarketNFTData struct {
	numberOfCoins   int64
	priceOfToken    int64
	validatorsCount int

	expectedBeneficiaryCoinsAmount int64
	expectedBuyerCoinsAmount       int64
	expectedSellerAmount           int64
}

func TestBuyPutOnMarketNFT(t *testing.T) {
	var (
		result sdk.Result
	)

	denom := types.DefaultTokenDenom

	testData := []testBuyPutOnMarketNFTData{
		{int64(1000), int64(600), 1,
			int64(1004), int64(400), int64(1586)},
		{int64(1000), int64(650), 1,
			int64(1004), int64(350), int64(1636)},
	}

	mpKeeperTest, err := createMarketplaceKeeperTest()
	defer mpKeeperTest.clear()
	require.Nil(t, err)

	for _, data := range testData {
		coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(data.numberOfCoins)))
		price := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(data.priceOfToken)))

		require.Nil(t, mpKeeperTest.updateAccountsWithCoins(coins))

		_, err := mpKeeperTest.updateVoteInfos(data.validatorsCount, coins)
		require.Nil(t, err)

		nft := createNFT(mpKeeperTest.addrs[0])
		handler := marketplace.NewHandler(mpKeeperTest.marketKeeper)
		require.Nil(t, mpKeeperTest.marketKeeper.MintNFT(mpKeeperTest.ctx, nft))

		putOnMarketNFTMsg := types.NewMsgPutOnMarketNFT(mpKeeperTest.addrs[0], mpKeeperTest.addrs[2], nft.GetID(), price)
		result = handler(mpKeeperTest.ctx, *putOnMarketNFTMsg)
		require.True(t, result.IsOK())

		buyNFTMsg := types.NewMsgBuyNFT(mpKeeperTest.addrs[1], mpKeeperTest.addrs[3], nft.GetID(), "")
		result = handler(mpKeeperTest.ctx, *buyNFTMsg)
		require.True(t, result.IsOK())

		// check seller's balance
		require.Equal(t, data.expectedSellerAmount,
			mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx, mpKeeperTest.addrs[0]).AmountOf(denom).Int64())

		// check buyer's balance
		require.Equal(t, data.expectedBuyerCoinsAmount,
			mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx, mpKeeperTest.addrs[1]).AmountOf(denom).Int64())

		// check of seller's beneficiary balance
		require.Equal(t, data.expectedBeneficiaryCoinsAmount,
			mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx, mpKeeperTest.addrs[2]).AmountOf(denom).Int64())

		// check buyer's beneficiary balance
		require.Equal(t, data.expectedBeneficiaryCoinsAmount,
			mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx, mpKeeperTest.addrs[3]).AmountOf(denom).Int64())
	}
}

type commissionTestData struct {
	amount     int64
	commission int64
	rat        float64
}

func TestCommission(t *testing.T) {
	denom := types.DefaultTokenDenom

	testData := []commissionTestData{
		{100, 1, 0.01},
		{200, 2, 0.01},
		{365, 3, 0.01},
		{1000, 500, 0.5},
		{1488, 982, 0.66},
		{10, 0, 0.001},
		{1000, 0, 0},
	}

	for _, data := range testData {
		price := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(data.amount)))
		require.Equal(t, data.commission, marketplace.GetCommission(price, data.rat).AmountOf(denom).Int64())
	}
}

func createNFT(owner sdk.AccAddress) *types.NFT {
	nft := marketplace.NewNFT(
		xnft.NewBaseNFT(
			uuid.New().String(),
			owner,
			"",
			"",
			"",
			"",
		),
		sdk.NewCoins(sdk.NewCoin(types.DefaultTokenDenom, sdk.NewInt(0))),
	)
	return nft
}
