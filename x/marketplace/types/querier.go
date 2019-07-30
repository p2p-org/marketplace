package types

import "strings"

type QueryResNFTs struct {
	NFTs []*NFT `json:"nfts"`
}

func (r QueryResNFTs) String() string {
	var out []string
	for _, nft := range r.NFTs {
		out = append(out, nft.String())
	}

	return strings.Join(out, "\n")
}

type QueryResFungibleTokens struct {
	FungibleTokens []*FungibleToken `json:"fungible_tokens"`
}

func (r QueryResFungibleTokens) String() string {
	var out []string
	for _, ft := range r.FungibleTokens {
		out = append(out, ft.String())
	}

	return strings.Join(out, "\n")
}
