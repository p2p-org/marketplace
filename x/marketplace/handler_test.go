package marketplace

import (
	"testing"

	"github.com/magiconair/properties/assert"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

func TestGetCommission(t *testing.T) {
	price := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(150)))

	// Single token case (validators + beneficiaries).
	expectedValsCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(1)))
	valsCommission := getCommission(price, ValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)

	expectedBeneficiariesCommission := sdk.NewCoins(sdk.NewCoin("test", sdk.NewInt(2)))
	beneficiariesCommission := getCommission(price, BeneficiariesCommission)
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
	valsCommission = getCommission(price, ValidatorsCommission)
	assert.Equal(t, valsCommission, expectedValsCommission)
}
