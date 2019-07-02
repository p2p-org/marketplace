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
