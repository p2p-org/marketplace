package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	NFTRecords []*NFT `json:"nft_records"`
}

func NewGenesisState(nftRecords []*NFT) GenesisState {
	return GenesisState{NFTRecords: nftRecords}
}

func ValidateGenesis(data GenesisState) error {
	// TODO: validate genesis.
	return nil
}

func DefaultGenesisState() GenesisState {
	return GenesisState{
		NFTRecords: []*NFT{},
	}
}

func InitGenesis(ctx sdk.Context, keeper Keeper, data GenesisState) []abci.ValidatorUpdate {
	for _, record := range data.NFTRecords {
		if err := keeper.MintNFT(ctx, record); err != nil {
			panic(fmt.Sprintf("failed to InitGenesis: %v", err))
		}
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var records []*NFT
	iterator := k.GetNFTsIterator(ctx)
	for ; iterator.Valid(); iterator.Next() {
		id := string(iterator.Key())
		nft, err := k.GetNFT(ctx, id)
		if err != nil {
			panic(fmt.Sprintf("failed to ExportGenesis: %v", err))
		}
		records = append(records, nft)
	}
	return GenesisState{NFTRecords: records}
}
