package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
)

type NFT struct {
	NFT    xnft.BaseNFT `json:"nft"`
	Price  sdk.Coin     `json:"price"`
	OnSale bool         `json:"on_sale"`
}

func NewNFT(nft xnft.BaseNFT, price sdk.Coin) *NFT {
	return &NFT{
		NFT:   nft,
		Price: price,
	}
}

func (m NFT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
NFT: %s
Price: %s
On Salse: %t`, m.NFT.GetOwner(), m.NFT, m.Price, m.OnSale))
}

func (m NFT) SetOnSale(status bool) {
	m.OnSale = status
}
