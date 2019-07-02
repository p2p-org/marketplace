package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the marketplace Querier
const (
	QueryNFT  = "nft"
	QueryNFTs = "nfts"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryNFT:
			return queryNFT(ctx, path[1:], req, keeper)
		case QueryNFTs:
			return queryNFTs(ctx, req, keeper)
		default:
			return nil, sdk.ErrUnknownRequest("unknown marketplace query endpoint")
		}
	}
}

// nolint: unparam
func queryNFT(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	id := path[0]
	value, err := keeper.GetNFT(ctx, id)
	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("could not find NFT with id %s: %v", id, err))
	}

	bz := keeper.cdc.MustMarshalJSON(value)
	return bz, nil
}

func queryNFTs(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var (
		nfts     types.QueryResNFTs
		iterator = keeper.GetNFTsIterator(ctx)
	)
	for ; iterator.Valid(); iterator.Next() {
		var nft types.NFT
		keeper.cdc.MustUnmarshalJSON(iterator.Value(), &nft)
		nfts.NFTs = append(nfts.NFTs, &nft)
	}

	return keeper.cdc.MustMarshalJSON(nfts), nil
}
