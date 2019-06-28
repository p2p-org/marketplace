package types

import (
	"fmt"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
)

type MarketplaceNFT struct {
	NFT    xnft.NFT `json:"nft"`
	Price  sdk.Coin `json:"price"`
	OnSale bool     `json:"on_sale"`
}

func NewMarketplaceNFT(nft xnft.NFT, price sdk.Coin) MarketplaceNFT {
	return MarketplaceNFT{
		NFT:   nft,
		Price: price,
	}
}

func (m MarketplaceNFT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
NFT: %s
Price: %s
On Salse: %t`, m.NFT.GetOwner(), m.NFT, m.Price, m.OnSale))
}

func (m MarketplaceNFT) SetOnSale(status bool) {
	m.OnSale = status
}
