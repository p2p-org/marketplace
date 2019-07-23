package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
)

type FungibleToken struct {
	Denom          string         `json:"denom"`
	EmissionAmount int64          `json:"emission_amount"`
	Creator        sdk.AccAddress `json:"creator"`
}

func (c FungibleToken) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Creator: %s
EmissionAmount: %d
Denom: %s`, c.Creator.String(), c.EmissionAmount, c.Denom))
}

type NFT struct {
	xnft.BaseNFT      `json:"nft"`
	Price             sdk.Coins      `json:"price"`
	OnSale            bool           `json:"on_sale"`
	SellerBeneficiary sdk.AccAddress `json:"seller_beneficiary"`
	TimeCreated       time.Time      `json:"time_created"`
}

func NewNFT(nft xnft.BaseNFT, price sdk.Coins) *NFT {
	return &NFT{
		BaseNFT:     nft,
		Price:       price,
		TimeCreated: time.Now().UTC(),
	}
}

func (m NFT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
NFT: %s
Price: %s
On Salse: %t`, m.GetOwner(), m.BaseNFT, m.Price, m.OnSale))
}

func (m *NFT) GetPrice() sdk.Coins {
	return m.Price
}

func (m *NFT) SetPrice(price sdk.Coins) {
	m.Price = price
}

func (m *NFT) IsOnSale() bool {
	return m.OnSale
}

func (m *NFT) SetOnSale(status bool) {
	m.OnSale = status
}

func (m *NFT) SetSellerBeneficiary(addr sdk.AccAddress) {
	m.SellerBeneficiary = addr
}

func (m *NFT) GetTimeCreated() time.Time {
	return m.TimeCreated
}
