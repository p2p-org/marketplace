package cli_test

import (
	"encoding/json"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
)

func TestGetCmdBatchPutOnMarket(t *testing.T) {
	str := `{"token1" : "123token",
	"token2" : "340stake"}`
	var pricesUnparsed map[string]string
	rm := json.RawMessage(str)
	err := json.Unmarshal(rm, &pricesUnparsed)
	assert.Nil(t, err)

	prices := make(map[string]sdk.Coins)
	for k, v := range pricesUnparsed {
		k, v := k, v
		price, err := sdk.ParseCoins(v)
		assert.Nil(t, err)
		prices[k] = price
	}

	t.Log(prices)
}
