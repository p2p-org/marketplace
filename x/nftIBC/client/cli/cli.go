package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/codec"
)

// GetQueryCmd returns the query commands for IBC fungible token transfer
func GetQueryCmd(cdc *codec.Codec, queryRoute string) *cobra.Command {
	ics20TransferQueryCmd := &cobra.Command{
		Use:   "transferNFT",
		Short: "IBC NFT transfer query subcommands",
	}

	ics20TransferQueryCmd.AddCommand(flags.GetCommands(
		GetCmdQueryNextSequence(cdc, queryRoute),
	)...)

	return ics20TransferQueryCmd
}

// GetTxCmd returns the transaction commands for IBC fungible token transfer
func GetTxCmd(cdc *codec.Codec) *cobra.Command {
	ics20TransferTxCmd := &cobra.Command{
		Use:   "transferNFT",
		Short: "IBC NFT transfer transaction subcommands",
	}

	ics20TransferTxCmd.AddCommand(flags.PostCommands(
		GetTransferNFTTxCmd(cdc),
	)...)

	return ics20TransferTxCmd
}
