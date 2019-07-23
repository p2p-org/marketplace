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

	FlagParamTokenName        = "name"
	FlagParamTokenNameShort   = "n"
	FlagParamDescription      = "description"
	FlagParamDescriptionShort = "d"
	FlagParamImage            = "image"
	FlagParamImageShort       = "i"
	FlagParamTokenURI         = "token_uri"
	FlagParamTokenURIShort    = "u"
	FlagParamPrice            = "price"
	FlagParamPriceShort       = "p"

	DefaultMaximumBeneficiaryCommission = 0.05
	DefaultBeneficiariesCommission      = 0.015
	DefaultValidatorsCommission         = 0.01

	DefaultTokenDenom = "token"
)

