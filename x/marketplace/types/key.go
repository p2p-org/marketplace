package types

const (
	// module name
	ModuleName = "marketplace"

	// StoreKey to be used when creating the KVStore
	StoreKey = ModuleName

	FungibleTokenCreationPrice = 10 // TODO: price or commission
	FungibleCommissionAddress  = "" // TODO: create account for commissions

	RouterKey = ModuleName

	FlagMaxCommission              = "max-commission"
	FlagBeneficiaryCommission      = "beneficiary-commission"
	FlagBeneficiaryCommissionShort = "c"

	DefaultMaximumBeneficiaryCommission = 0.05
	DefaultBeneficiariesCommission      = 0.015
	DefaultValidatorsCommission         = 0.01
)

