package marketplace_test

import (
	"testing"

	"github.com/corestario/marketplace/x/marketplace"
	"github.com/corestario/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

func TestGetCommission(t *testing.T) {
	price := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(150)))

	// Single token case (validators + beneficiaries).
	expectedValsCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1)))
	valsCommission := marketplace.GetCommission(price, types.DefaultValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)

	expectedBeneficiariesCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(2)))
	beneficiariesCommission := marketplace.GetCommission(price, types.DefaultBeneficiariesCommission)
	assert.Equal(t, beneficiariesCommission, expectedBeneficiariesCommission)

	// Multiple tokens case (validators).
	price = sdk.NewCoins(
		sdk.NewCoin("test1", sdk.NewInt(150)),
		sdk.NewCoin("test2", sdk.NewInt(150)),
	)
	expectedValsCommission = sdk.NewCoins(
		sdk.NewCoin("test1", sdk.NewInt(1)),
		sdk.NewCoin("test2", sdk.NewInt(1)),
	)
	valsCommission = marketplace.GetCommission(price, types.DefaultValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)
}

type sendPair struct {
	Src    sdk.AccAddress
	Dst    sdk.AccAddress
	Amount int64
}

type testRollBackData struct {
	pairs       []sendPair
	startValues int64
}

func createAccAddressSlice(pairs []sendPair) []sdk.AccAddress {
	tmpMap := make(map[string]bool)
	out := make([]sdk.AccAddress, 0)
	for _, v := range pairs {
		v := v
		if _, ok := tmpMap[v.Src.String()]; !ok {
			tmpMap[v.Src.String()] = true
			out = append(out, v.Src)
		}
		if _, ok := tmpMap[v.Dst.String()]; !ok {
			tmpMap[v.Dst.String()] = true
			out = append(out, v.Dst)
		}
	}
	return out
}

func TestRollbackCommission(t *testing.T) {
	logger := log.NewNopLogger()

	mpKeeperTest, err := createMarketplaceKeeperTest()
	defer mpKeeperTest.clear()
	require.Nil(t, err)

	denom := types.DefaultTokenDenom

	coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(1000)))
	require.Nil(t, mpKeeperTest.updateAccountsWithCoins(coins))

	user1 := mpKeeperTest.addrs[0]
	user2 := mpKeeperTest.addrs[1]
	sellerBeneficiary := mpKeeperTest.addrs[2]
	buyerBeneficiary := mpKeeperTest.addrs[3]

	testData := []testRollBackData{
		{[]sendPair{{user1, user2, 100}, {user1, user2, 10}}, 500},
		{[]sendPair{{sellerBeneficiary, buyerBeneficiary, 300},
			{user1, sellerBeneficiary, 100}, {user1, user2, 100},
		}, 400},
		{[]sendPair{{sellerBeneficiary, buyerBeneficiary, 100},
			{sellerBeneficiary, user1, 100}, {sellerBeneficiary, user2, 100},
			{sellerBeneficiary, sellerBeneficiary, 100}, {sellerBeneficiary, buyerBeneficiary, 100},
		}, 1000},
		{[]sendPair{{user2, sellerBeneficiary, 40}, {user1, buyerBeneficiary, 50},
			{user1, user2, 100}, {user1, user2, 700},
			{user1, buyerBeneficiary, 100},
		}, 1500},
	}

	for _, data := range testData {
		data := data
		coins := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(data.startValues)))
		_, err := mpKeeperTest.updateVoteInfos(1, coins)
		require.Nil(t, err)
		require.Nil(t, mpKeeperTest.updateAccountsWithCoins(coins))
		nft := createNFT(user1)
		require.Nil(t, mpKeeperTest.marketKeeper.MintNFT(mpKeeperTest.ctx, nft))

		balanceAddrs := createAccAddressSlice(data.pairs)
		initialBalances := marketplace.GetBalances(mpKeeperTest.ctx, mpKeeperTest.marketKeeper, balanceAddrs...)
		// send coins between pairs
		for _, pair := range data.pairs {
			pair := pair
			coinsToSend := sdk.NewCoins(sdk.NewCoin(denom, sdk.NewInt(pair.Amount)))
			require.Nil(t, mpKeeperTest.bankKeeper.SendCoins(mpKeeperTest.ctx, pair.Src, pair.Dst, coinsToSend))
		}
		// check that balances have changed
		for _, addr := range balanceAddrs {
			addr := addr
			require.NotEqual(t, coins.AmountOf(denom).Int64(),
				mpKeeperTest.bankKeeper.GetAllBalances(mpKeeperTest.ctx, addr).AmountOf(denom).Int64())
		}
		// rollback
		marketplace.RollbackCommissions(mpKeeperTest.ctx, mpKeeperTest.marketKeeper, logger, initialBalances)
		// check that balances are restored
		for _, addr := range mpKeeperTest.addrs {
			addr := addr
			require.Equal(t, coins.AmountOf(denom).Int64(),
				mpKeeperTest.bankKeeper.GetAllBalances(mpKeeperTest.ctx, addr).AmountOf(denom).Int64())
		}
	}
}
