package types

const (
	// module name
	ModuleName = "marketplace"

	// StoreKey to be used when creating the KVStore
	StoreKey  = ModuleName
	RouterKey = ModuleName

	FlagMaxCommission              = "max-commission"
	FlagBeneficiaryCommission      = "beneficiary-commission"
	FlagBeneficiaryCommissionShort = "c"
)

const (
	DefaultMaximumBeneficiaryCommission = 0.05
	DefaultBeneficiariesCommission      = 0.015
	DefaultValidatorsCommission         = 0.01
)
