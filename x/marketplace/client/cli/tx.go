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
		GetCmdTransferNFT(cdc),
		GetCmdSellNFT(cdc),
		GetCmdBuyNFT(cdc),
	)...)

	return marketplaceTxCmd
}

func GetCmdMintNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint [name] [description] [image] [token_uri] [price]",
		Short: "mint a new NFT",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				owner       = cliCtx.GetFromAddress()
				name        = args[0]
				description = args[1]
				image       = args[2]
				tokenURI    = args[3]
			)
			//price, err := sdk.ParseCoin(args[4])
			//if err != nil {
			//	return err
			//}

			msg := types.NewMsgMintNFT(owner, name, description, image, tokenURI)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdTransferNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transfer [token_id] [recipient]",
		Short: "transfer an NFT from one account to another",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			recipient, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgTransferNFT(args[0], cliCtx.GetFromAddress(), recipient)
			if err = msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdSellNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "sell [token_id] [price]",
		Short: "sell an NFT (token can be bought for the specified price)",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			price, err := sdk.ParseCoin(args[1])
			if err != nil {
				return err
			}

			msg := types.NewMsgSellNFT(cliCtx.GetFromAddress(), args[0], price)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdBuyNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "buy [token_id]",
		Short: "buy an NFT",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgBuyNFT(cliCtx.GetFromAddress(), args[0])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}
