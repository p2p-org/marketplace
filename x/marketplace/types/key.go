package types

import (
	"fmt"
)

type NFTStatus int8

func (s NFTStatus) String() string {
	switch s {
	case NFTStatusDefault:
		return "default"
	case NFTStatusOnMarket:
		return "on_market"
	case NFTStatusOnAuction:
		return "on_auction"
	case NFTStatusDeleted:
		return "deleted"
	case NFTStatusUndefined:
		return "undefined"
	}
	return "undefined"
}

func (s NFTStatus) MarshalJSON() ([]byte, error) {
	r := fmt.Sprintf("\"%v\"", s)
	return []byte(r), nil
}

func (s *NFTStatus) UnmarshalJSON(b []byte) error {
	var t string
	t = string(b)
	var e NFTStatus
	switch t {
	case "\"default\"":
		e = NFTStatus(0)
	case "\"on_market\"":
		e = NFTStatus(1)
	case "\"on_auction\"":
		e = NFTStatus(2)
	case "\"deleted\"":
		e = NFTStatus(3)
	case "\"undefined\"":
		e = NFTStatus(4)
	default:
		e = NFTStatus(0)
	}

	*s = e
	return nil
}

const (
	NFTStatusDefault NFTStatus = iota
	NFTStatusOnMarket
	NFTStatusOnAuction
	NFTStatusDeleted
	NFTStatusUndefined
)

const (
	// module name
	ModuleName = "marketplace"

	// StoreKey to be used when creating the KVStore
	StoreKey         = ModuleName
	RegisterCurrency = "register_currency"
	AuctionKey       = "auction"
	DeletedNFTKey    = "deleted_nft"

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

	FlagParamBuyoutPrice      = "buyout"
	FlagParamBuyoutPriceShort = "u"

	DefaultMaximumBeneficiaryCommission = 0.05
	DefaultBeneficiariesCommission      = 0.015
	DefaultValidatorsCommission         = 0.01

	DefaultTokenDenom = "token"

	MaxTokenIDLength     = 36
	MaxNameLength        = 50
	MaxDescriptionLength = 32000
	MaxImageLength       = 32000
	MaxTokenURILength    = 32000
	MaxDenomLength       = 16
	MinDenomLength       = 3
	IBCNFTPort           = "transfernft"

	DefaultFinishAuctionHost    = "localhost"
	DefaultFinishAuctionPort    = 1317
	DefaultChainName            = "mpchain"
	DefaultFinishingAccountName = "dgaming"
	DefaultFinishingAccountPass = "12345678"
	DefaultFinishingAccountAddr = "cosmos1tctr64k4en25uvet2k2tfkwkh0geyrv8fvuvet"
)
