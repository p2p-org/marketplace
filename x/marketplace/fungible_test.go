package marketplace_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgamingfoundation/marketplace/x/marketplace"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	"github.com/stretchr/testify/require"
	"testing"
)

type createAndTransferFT struct {
	startBalance int64
	denom        string

	denomFT          string
	emissionFTAmount int64
	transferFTAmount int64

	expectedCreatorBalanceAfterCreation     int64
	expectedCreatorFTBalanceAfterTransfer   int64
	expectedRecipientFTBalanceAfterTransfer int64
}

func TestCreateAndTransferFT(t *testing.T) {
	var (
		result sdk.Result
	)

	// TODO: fix expectedCreatorBalanceAfterCreation
	testData := []createAndTransferFT{
		{1000, "token", "test", 1000, 500,
			990, 500, 500},
		{1000, "token", "tester", 1000, 350,
			990, 650, 350},
	}

	mpKeeperTest, err := createMarketplaceKeeperTest()
	defer mpKeeperTest.clear()
	require.Nil(t, err)

	for _, data := range testData {
		coins := sdk.NewCoins(sdk.NewCoin(data.denom, sdk.NewInt(data.startBalance)))

		require.Nil(t, mpKeeperTest.updateAccountsWithCoins(coins))

		_, err := mpKeeperTest.updateVoteInfos(1, coins)
		require.Nil(t, err)

		handler := marketplace.NewHandler(mpKeeperTest.marketKeeper)

		createFT := types.NewMsgCreateFungibleToken(mpKeeperTest.addrs[0], data.denomFT, data.emissionFTAmount)
		result = handler(mpKeeperTest.ctx, *createFT)
		require.True(t, result.IsOK())
		require.Equal(t, data.expectedCreatorBalanceAfterCreation,
			mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx, mpKeeperTest.addrs[0]).AmountOf(data.denom).Int64())

		transferFT := types.NewMsgTransferFungibleTokens(mpKeeperTest.addrs[0], mpKeeperTest.addrs[1], data.denomFT, data.transferFTAmount)
		result = handler(mpKeeperTest.ctx, *transferFT)
		require.True(t, result.IsOK())
		require.Equal(t, data.expectedRecipientFTBalanceAfterTransfer, mpKeeperTest.bankKeeper.GetCoins(mpKeeperTest.ctx,
			mpKeeperTest.addrs[1]).AmountOf(data.denomFT).Int64())
	}
}
