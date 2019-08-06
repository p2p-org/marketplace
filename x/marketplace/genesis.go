package marketplace

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

type GenesisState struct {
	NFTRecords           []*NFT          `json:"nft_records"`
	RegisteredCurrencies []FungibleToken `json:"registered_tokens"`
}

func NewGenesisState(nftRecords []*NFT) GenesisState {
	return GenesisState{NFTRecords: nftRecords}
}

func ValidateGenesis(data GenesisState) error {
	for _, cur := range data.RegisteredCurrencies {
		if cur.Creator == nil {
			return fmt.Errorf("invalid FungibleToken: Denom: %s. Error: Missing Creator", cur.Denom)
		}
		if cur.Denom == "" {
			return fmt.Errorf("invalid FungibleToken: Creator: %v. Error: Missing Denom", cur.Creator)
		}
	}
	for _, record := range data.NFTRecords {
		if record.Owner == nil {
			return fmt.Errorf("invalid NFTRecord: Name: %v. Error: Missing Owner", record.Name)
		}
		if record.BaseNFT.Owner == nil {
			return fmt.Errorf("invalid NFTRecord: Name: %v. Error: Missing BaseNFT.Owner", record.Name)
		}
		if record.ID == "" {
			return fmt.Errorf("invalid NFTRecord: Name: %v. Error: Missing ID", record.Name)
		}
		if record.Name == "" {
			return fmt.Errorf("invalid NFTRecord: ID: %v. Error: Missing Name", record.ID)
		}
		if record.TimeCreated.IsZero() {
			return fmt.Errorf("invalid NFTRecord: Name: %v. Error: Missing TimeCreated", record.Name)
		}
	}

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

	for _, currency := range data.RegisteredCurrencies {
		keeper.registerFungibleTokensCurrency(ctx, currency)
	}
	return []abci.ValidatorUpdate{}
}

func ExportGenesis(ctx sdk.Context, k Keeper) GenesisState {
	var (
		records    []*NFT
		currencies []FungibleToken
		currency   FungibleToken
	)
	nftIterator := k.GetNFTsIterator(ctx)
	for ; nftIterator.Valid(); nftIterator.Next() {
		id := string(nftIterator.Key())
		nft, err := k.GetNFT(ctx, id)
		if err != nil {
			panic(fmt.Sprintf("failed to ExportGenesis: %v", err))
		}
		records = append(records, nft)
	}

	currIterator := k.GetRegisteredCurrenciesIterator(ctx)
	for ; currIterator.Valid(); currIterator.Next() {
		k.cdc.MustUnmarshalBinaryBare(currIterator.Value(), &currency)
		currencies = append(currencies, currency)
	}
	return GenesisState{NFTRecords: records, RegisteredCurrencies: currencies}
}
