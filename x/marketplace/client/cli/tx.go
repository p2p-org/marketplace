package cli

import (
	"fmt"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/auth/client/utils"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"
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
		GetCmdPutNFTOnMarket(cdc),
		GetCmdRemoveNFTFromMarket(cdc),
		GetCmdBuyNFT(cdc),
		GetCmdUpdateNFTParams(cdc),
		GetCmdCreateFungibleToken(cdc),
		GetCmdTransferFungibleTokens(cdc),
		GetCmdUpdateNFTParams(cdc),
		GetCmdPutNFTOnAuction(cdc),
		GetCmdRemoveNFTFromAuction(cdc),
		GetCmdFinishAuction(cdc),
		GetCmdMakeBidOnAuction(cdc),
		GetCmdBuyoutFromAuction(cdc),
		GetCmdBurnFungibleTokens(cdc),
		GetCmdMakeOffer(cdc),
		GetCmdAcceptOffer(cdc),
	)...)

	return marketplaceTxCmd
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

func GetCmdRemoveNFTFromMarket(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove_from_market [token_id]",
		Short: "remove an NFT from market",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgRemoveNFTFromMarket(cliCtx.GetFromAddress(), args[0])
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
			commission := viper.GetString(types.FlagBeneficiaryCommission)

			msg := types.NewMsgBuyNFT(cliCtx.GetFromAddress(), beneficiary, args[0], commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(types.FlagBeneficiaryCommission, types.FlagBeneficiaryCommissionShort, types.DefaultBeneficiariesCommission,
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
			if price := viper.GetString(types.FlagParamPrice); price != "" {
				params = append(params, types.NFTParam{Key: types.FlagParamPrice, Value: viper.GetString(types.FlagParamPrice)})
			}
			if name := viper.GetString(types.FlagParamTokenName); name != "" {
				params = append(params, types.NFTParam{Key: types.FlagParamTokenName, Value: name})
			}
			if uri := viper.GetString(types.FlagParamTokenURI); uri != "" {
				params = append(params, types.NFTParam{Key: types.FlagParamTokenURI, Value: uri})
			}
			if img := viper.GetString(types.FlagParamImage); img != "" {
				params = append(params, types.NFTParam{Key: types.FlagParamImage, Value: img})
			}
			if desc := viper.GetString(types.FlagParamDescription); desc != "" {
				params = append(params, types.NFTParam{Key: types.FlagParamDescription, Value: desc})
			}

			msg := types.NewMsgUpdateNFTParams(cliCtx.GetFromAddress(), args[0], params)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(types.FlagParamPrice, types.FlagParamPriceShort, "", "new nft price, if left blank will not be changed")
	cmd.Flags().StringP(types.FlagParamTokenName, types.FlagParamTokenNameShort, "", "new nft name, if left blank will not be changed")
	cmd.Flags().StringP(types.FlagParamImage, types.FlagParamImageShort, "", "new nft image, if left blank will not be changed")
	cmd.Flags().StringP(types.FlagParamTokenURI, types.FlagParamTokenURIShort, "", "new nft uri, if left blank will not be changed")
	cmd.Flags().StringP(types.FlagParamDescription, types.FlagParamDescriptionShort, "", "new nft description, if left blank will not be changed")
	return cmd
}

func GetCmdPutNFTOnAuction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "put_on_auction [token_id] [opening_price] [beneficiary] [duration]",
		Short: "put on auction an NFT (token will be traded in specified time or returned to owner)",
		Args:  cobra.ExactArgs(4),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			openingPrice, err := sdk.ParseCoins(args[1])
			if err != nil {
				return fmt.Errorf("failed to parce openingPrice: %v", err)
			}

			buyout := viper.GetString(types.FlagParamBuyoutPrice)
			buyoutPrice, err := sdk.ParseCoins(buyout)
			if err != nil {
				return fmt.Errorf("failed to parce buyoutPrice: %v", err)
			}

			beneficiary, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}

			dur, err := time.ParseDuration(args[3])
			if err != nil {
				return fmt.Errorf("failed to parse duration of auction: %v", err)
			}

			msg := types.NewMsgPutNFTOnAuction(cliCtx.GetFromAddress(), beneficiary, args[0], openingPrice, buyoutPrice, time.Now().UTC().Add(dur))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().StringP(types.FlagParamBuyoutPrice, types.FlagParamBuyoutPriceShort, "",
		"buyout price for auction lot, if left blank will have no buyout price")
	return cmd
}

func GetCmdRemoveNFTFromAuction(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "remove_from_auction [token_id]",
		Short: "remove an NFT from action",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgRemoveNFTFromAuction(cliCtx.GetFromAddress(), args[0])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdFinishAuction(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "finish_auction [token_id]",
		Short: "finish an NFT action",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			msg := types.NewMsgFinishAuction(cliCtx.GetFromAddress(), args[0])
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdMakeBidOnAuction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "bid [token_id] [beneficiary] [price]",
		Short: "make a bid for an NFT on auction",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			beneficiary, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}
			commission := viper.GetString(types.FlagBeneficiaryCommission)
			price, err := sdk.ParseCoins(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse price: %v", err)
			}

			msg := types.NewMsgMakeBidOnAuction(cliCtx.GetFromAddress(), beneficiary, args[0], price, commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(types.FlagBeneficiaryCommission, types.FlagBeneficiaryCommissionShort, types.DefaultBeneficiariesCommission,
		"beneficiary fee, if left blank will be set to default")
	return cmd
}

func GetCmdBuyoutFromAuction(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "buyout [token_id] [beneficiary]",
		Short: "buyout an NFT from auction",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			beneficiary, err := sdk.AccAddressFromBech32(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}

			commission := viper.GetString(types.FlagBeneficiaryCommission)
			msg := types.NewMsgBuyOutOnAuction(cliCtx.GetFromAddress(), beneficiary, args[0], commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(types.FlagBeneficiaryCommission, types.FlagBeneficiaryCommissionShort, types.DefaultBeneficiariesCommission,
		"beneficiary fee, if left blank will be set to default")
	return cmd
}

func GetCmdBurnFungibleTokens(cdc *codec.Codec) *cobra.Command {
	return &cobra.Command{
		Use:   "burnFT [denom] [amount]",
		Short: "burn some amount of owned fungible tokens",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			amount, err := strconv.Atoi(args[1])
			if err != nil {
				return fmt.Errorf("failed to parse amount: %v", err)
			}

			msg := types.NewMsgBurnFungibleTokens(cliCtx.GetFromAddress(), args[0], int64(amount))
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
}

func GetCmdMakeOffer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "offer [token_id] [price] [beneficiary]",
		Short: "offer a price for an NFT that is not currently on sale",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			beneficiary, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}

			price, err := sdk.ParseCoins(args[1])
			if err != nil {
				return err
			}

			commission := viper.GetString(types.FlagBeneficiaryCommission)
			msg := types.NewMsgMakeOffer(cliCtx.GetFromAddress(), beneficiary, price, args[0], commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(types.FlagBeneficiaryCommission, types.FlagBeneficiaryCommissionShort, types.DefaultBeneficiariesCommission,
		"beneficiary fee, if left blank will be set to default")
	return cmd
}

func GetCmdAcceptOffer(cdc *codec.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "accept_offer [token_id] [offer_id] [beneficiary]",
		Short: "accept an offer",
		Args:  cobra.ExactArgs(3),
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := auth.NewTxBuilderFromCLI().WithTxEncoder(utils.GetTxEncoder(cdc))
			cliCtx := context.NewCLIContext().WithCodec(cdc)

			beneficiary, err := sdk.AccAddressFromBech32(args[2])
			if err != nil {
				return fmt.Errorf("failed to parse beneficiary address: %v", err)
			}

			tokenID, offerID := args[0], args[1]

			commission := viper.GetString(types.FlagBeneficiaryCommission)
			msg := types.NewMsgAcceptOffer(cliCtx.GetFromAddress(), beneficiary, tokenID, offerID, commission)
			if err := msg.ValidateBasic(); err != nil {
				return err
			}

			return utils.GenerateOrBroadcastMsgs(cliCtx, txBldr, []sdk.Msg{msg})
		},
	}
	cmd.Flags().Float64P(types.FlagBeneficiaryCommission, types.FlagBeneficiaryCommissionShort, types.DefaultBeneficiariesCommission,
		"beneficiary fee, if left blank will be set to default")
	return cmd
}
