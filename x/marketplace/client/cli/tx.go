package cli

import (
	"fmt"
	"strconv"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
	mptypes "github.com/dgamingfoundation/marketplace/x/marketplace/types"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
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
		GetCmdPutNFTOnMarket(cdc),
		GetCmdBuyNFT(cdc),
		GetCmdUpdateNFTParams(cdc),
		GetCmdCreateFungibleToken(cdc),
		GetCmdTransferFungibleTokens(cdc),
		GetCmdUpdateNFTParams(cdc),
	)...)

	return marketplaceTxCmd
}

func GetCmdMintNFT(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "mint [token_id] [name] [description] [image] [token_uri]",
		Short: "mint a new NFT",
		Args:  cobra.ExactArgs(5),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			var (
				owner       = cliCtx.GetFromAddress()
				tokenID     = args[0]
				name        = args[1]
				description = args[2]
				image       = args[3]
				tokenURI    = args[4]
			)
			msg := types.NewMsgMintNFT(tokenID, owner, name, description, image, tokenURI)
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

func GetCmdPutNFTOnMarket(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "put_on_market [token_id] [price] [beneficiary]",
		Short: "put on market an NFT (token can be bought for the specified price)",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			price, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			beneficiary, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}

			msg := types.NewMsgPutOnMarketNFT(cliCtx.GetFromAddress(), beneficiary, args[0], price)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdBuyNFT(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buy [token_id] [beneficiary]",
		Short: "buy an NFT",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			beneficiary, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}
			commission := viper.GetString(mptypes.FlagBeneficiaryCommission)

			msg := types.NewMsgBuyNFT(cliCtx.GetFromAddress(), beneficiary, args[0], commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(mptypes.FlagBeneficiaryCommission, mptypes.FlagBeneficiaryCommissionShort, mptypes.DefaultBeneficiariesCommission,
		"beneficiary fee, if left blank will be set to default")
	return cmd
}

func GetCmdCreateFungibleToken(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "createFT [denom] [amount]",
		Short: "create a fungible token",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			amount, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse amount: %v", err)
			}

			msg := types.NewMsgCreateFungibleToken(cliCtx.GetFromAddress(), args[0], int64(amount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdTransferFungibleTokens(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "transferFT [recipient] [denom] [amount]",
		Short: "transfer fungible tokens to another account",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			recipient, err := sdk.AccAddressFromBech32(args[0])
			if err != nil {
				return fmt.Errorf("failed to parse recipient address: %v", err)
			}

			amount, err := strconv.Atoi(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse amount: %v", err)
			}

			msg := types.NewMsgTransferFungibleTokens(cliCtx.GetFromAddress(), recipient, args[1], int64(amount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdUpdateNFTParams(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update_params [token_id]",
		Short: "update params of an NFT",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)
			params := make([]types.NFTParam, 0)
			if price := viper.GetString(mptypes.FlagParamPrice); price != "" {
				params = append(params, types.NFTParam{Key: mptypes.FlagParamPrice, Value: viper.GetString(mptypes.FlagParamPrice)})
			}
			if name := viper.GetString(mptypes.FlagParamTokenName); name != "" {
				params = append(params, types.NFTParam{Key: mptypes.FlagParamTokenName, Value: name})
			}
			if uri := viper.GetString(mptypes.FlagParamTokenURI); uri != "" {
				params = append(params, types.NFTParam{Key: mptypes.FlagParamTokenURI, Value: uri})
			}
			if img := viper.GetString(mptypes.FlagParamImage); img != "" {
				params = append(params, types.NFTParam{Key: mptypes.FlagParamImage, Value: img})
			}
			if desc := viper.GetString(mptypes.FlagParamDescription); desc != "" {
				params = append(params, types.NFTParam{Key: mptypes.FlagParamDescription, Value: desc})
			}

			msg := types.NewMsgUpdateNFTParams(cliCtx.GetFromAddress(), args[0], params)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(mptypes.FlagParamPrice, mptypes.FlagParamPriceShort, "", "new nft price, if left blank will not be changed")
	cmd.Flags().StringP(mptypes.FlagParamTokenName, mptypes.FlagParamTokenNameShort, "", "new nft name, if left blank will not be changed")
	cmd.Flags().StringP(mptypes.FlagParamImage, mptypes.FlagParamImageShort, "", "new nft image, if left blank will not be changed")
	cmd.Flags().StringP(mptypes.FlagParamTokenURI, mptypes.FlagParamTokenURIShort, "", "new nft uri, if left blank will not be changed")
	cmd.Flags().StringP(mptypes.FlagParamDescription, mptypes.FlagParamDescriptionShort, "", "new nft description, if left blank will not be changed")
	return cmd
}
