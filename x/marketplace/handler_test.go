package marketplace

import (
	"testing"

	"github.com/dgamingfoundation/marketplace/x/marketplace/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/magiconair/properties/assert"
)

func TestGetCommission(t *testing.T) {
	price := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(150)))

	// Single token case (validators + beneficiaries).
	expectedValsCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1)))
	valsCommission := GetCommission(price, types.DefaultValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)

	expectedBeneficiariesCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(2)))
	beneficiariesCommission := GetCommission(price, types.DefaultBeneficiariesCommission)
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
	valsCommission = GetCommission(price, types.DefaultValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)
}
