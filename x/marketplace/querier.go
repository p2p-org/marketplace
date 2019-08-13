package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// query endpoints supported by the marketplace Querier
const (
	QueryNFT            = "nft"
	QueryNFTs           = "nfts"
	QueryFungibleToken  = "fungible_token"
	QueryFungibleTokens = "fungible_tokens"
	QueryAuctionLot     = "auction_lot"
	QueryAuctionLots    = "auction_lots"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper Keeper) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err sdk.Error) {
		switch path[0] {
		case QueryNFT:
			return queryNFT(ctx, path[1:], req, keeper)
		case QueryNFTs:
			return queryNFTs(ctx, req, keeper)
		case QueryFungibleToken:
			return queryFungibleToken(ctx, path[1:], req, keeper)
		case QueryFungibleTokens:
			return queryFungibleTokens(ctx, req, keeper)
		case QueryAuctionLot:
			return queryAuctionLot(ctx, path[1:], req, keeper)
		case QueryAuctionLots:
			return queryAuctionLots(ctx, req, keeper)
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
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var nft types.NFT
		keeper.cdc.MustUnmarshalJSON(iterator.Value(), &nft)
		nfts.NFTs = append(nfts.NFTs, &nft)
	}

	return keeper.cdc.MustMarshalJSON(nfts), nil
}

func queryFungibleToken(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	name := path[0]
	value, err := keeper.GetFungibleToken(ctx, name)
	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("could not find Fungible Token with name %s: %v", name, err))
	}

	bz := keeper.cdc.MustMarshalJSON(value)
	return bz, nil
}

func queryFungibleTokens(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var fts types.QueryResFungibleTokens
	iterator := keeper.GetFungibleTokensIterator(ctx)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		var ft types.FungibleToken
		keeper.cdc.MustUnmarshalJSON(iterator.Value(), &ft)
		fts.FungibleTokens = append(fts.FungibleTokens, &ft)
	}

	return keeper.cdc.MustMarshalJSON(fts), nil
}

func queryAuctionLot(ctx sdk.Context, path []string, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	id := path[0]
	value, err := keeper.GetAuctionLot(ctx, id)
	if err != nil {
		return []byte{}, sdk.ErrUnknownRequest(fmt.Sprintf("could not find AuctionLot with id %s: %v", id, err))
	}

	bz := keeper.cdc.MustMarshalJSON(value)
	return bz, nil
}

func queryAuctionLots(ctx sdk.Context, req abci.RequestQuery, keeper Keeper) ([]byte, sdk.Error) {
	var (
		lots     types.QueryResAuctionLots
		iterator = keeper.GetAuctionLotsIterator(ctx)
	)
	defer iterator.Close()
	for ; iterator.Valid(); iterator.Next() {
		var lot types.AuctionLot
		keeper.cdc.MustUnmarshalJSON(iterator.Value(), &lot)
		lots.AuctionLots = append(lots.AuctionLots, &lot)
	}

	return keeper.cdc.MustMarshalJSON(lots), nil
}
