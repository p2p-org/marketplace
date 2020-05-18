package marketplace_test

import (
	"testing"

	"github.com/p2p-org/marketplace/x/marketplace"
	"github.com/p2p-org/marketplace/x/marketplace/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
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

	// TODO: fix expectedCreatorBalanceAfterCreation
	testData := []createAndTransferFT{
		{1000, types.DefaultTokenDenom, "test", 1000, 500,
			990, 500, 500},
		{1000, types.DefaultTokenDenom, "tester", 1000, 350,
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
		_, err = handler(mpKeeperTest.ctx, *createFT)
		require.NoError(t, err)
		require.Equal(t, data.expectedCreatorBalanceAfterCreation,
			mpKeeperTest.bankKeeper.GetAllBalances(mpKeeperTest.ctx, mpKeeperTest.addrs[0]).AmountOf(data.denom).Int64())

		transferFT := types.NewMsgTransferFungibleTokens(mpKeeperTest.addrs[0], mpKeeperTest.addrs[1], data.denomFT, data.transferFTAmount)
		_, err = handler(mpKeeperTest.ctx, *transferFT)
		require.NoError(t, err)
		require.Equal(t, data.expectedRecipientFTBalanceAfterTransfer, mpKeeperTest.bankKeeper.GetAllBalances(mpKeeperTest.ctx,
			mpKeeperTest.addrs[1]).AmountOf(data.denomFT).Int64())
	}
}
