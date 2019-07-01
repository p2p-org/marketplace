package cli

import (
	"github.com/spf13/cobra"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
)

func GetTxCmd(storeKey string, cdc *codec.Codec) *cobra.Command {
	marketplaceTxCmd := &cobra.Command{
		Use:                        types.ModuleName,
		Short:                      "Marketplace transaction subcommands",
		DisableFlagParsing:         true,
		SuggestionsMinimumDistance: 2,
		RunE:                       client.ValidateCmd,
	}

	marketplaceTxCmd.AddCommand(client.PostCommands(
		GetCmdMintNFT(cdc),
	)...)

	return marketplaceTxCmd
}

// GetCmdMintNFT is the CLI command for sending a BuyName transaction
func GetCmdMintNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint [name] [description] [image] [token_uri] [price]",
		Short: "mint a new NFT",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContextWithFrom(args[0]).WithCodec(cdc)

			var (
				owner       = cliCtx.GetFromAddress()
				name        = args[0]
				description = args[1]
				image       = args[2]
				tokenURI    = args[3]
			)
			price, err := sdk.ParseCoin(args[4])
			if err != nil {
				return err
			}

			msg := types.NewMsgMintNFT(owner, name, description, image, tokenURI, price)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
