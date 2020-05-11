package rest

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/corestario/marketplace/x/marketplace/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/flags"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	"github.com/cosmos/cosmos-sdk/x/auth"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
	authTypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/cosmos/modules/incubator/nft"

	"github.com/gorilla/mux"
)

const (
	restName = "name"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router, storeName string) {
	r.HandleFunc(fmt.Sprintf("/%s/nfts", storeName), nftsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/nfts/{%s}", storeName, restName), nftHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/fungible_tokens", storeName), fungibleTokensHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/fungible_tokens/{%s}", storeName, restName), fungibleTokenHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/auction_lots", storeName), auctionLotsHandler(cliCtx, storeName)).Methods("GET")
	r.HandleFunc(fmt.Sprintf("/%s/auction_lots/{%s}", storeName, restName), auctionLotHandler(cliCtx, storeName)).Methods("GET")

	r.HandleFunc(fmt.Sprintf("/%s/mint", storeName), mintHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/transfer", storeName), transferHandler(cliCtx)).Methods("PUT")

	r.HandleFunc(fmt.Sprintf("/%s/put_on_market", storeName), putOnMarketHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/buy", storeName), buyHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/update_params", storeName), updateParamsHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/remove_from_market", storeName), removeNFTFromMarketHandler(cliCtx)).Methods("PUT")

	r.HandleFunc(fmt.Sprintf("/%s/create_ft", storeName), createFTHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/transfer_ft", storeName), transferFTHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/burn_ft", storeName), burnFTHandler(cliCtx)).Methods("PUT")

	r.HandleFunc(fmt.Sprintf("/%s/put_on_auction", storeName), putOnAuctionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/remove_from_auction", storeName), removeNFTFromAuctionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/finish_auction", storeName), finishAuctionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/bid_on_auction", storeName), bidOnAuctionHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/buyout_auction", storeName), buyoutAuctionHandler(cliCtx)).Methods("PUT")

	r.HandleFunc(fmt.Sprintf("/%s/txs", storeName), unifiedHandler(cliCtx)).Methods("POST")
}

// --------------------------------------------------------------------------------------
//
// Tx Handler
//
// --------------------------------------------------------------------------------------

func unifiedHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		msgBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		var ttx authTypes.StdTx

		err = cliCtx.Codec.UnmarshalJSON(msgBytes, &ttx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		txBytes, err := cliCtx.Codec.MarshalBinaryLengthPrefixed(ttx)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusInternalServerError, err.Error())
			return
		}
		resp, err := cliCtx.BroadcastTxSync(txBytes)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, resp)
	}
}

func broadcastTransaction(
	cliCtx context.CLIContext,
	w http.ResponseWriter,
	msg sdk.Msg,
	bq rest.BaseReq,
	name,
	password string) {

	gasAdj, ok := rest.ParseFloat64OrReturnBadRequest(w, bq.GasAdjustment, flags.DefaultGasAdjustment)
	if !ok {
		return
	}

	_, gas, err := flags.ParseGas(bq.Gas)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	txBldr := auth.NewTxBuilder(
		authclient.GetTxEncoder(cliCtx.Codec), bq.AccountNumber, bq.Sequence, gas, gasAdj,
		bq.Simulate, bq.ChainID, bq.Memo, bq.Fees, bq.GasPrices,
	)

	msgBytes, err := txBldr.BuildAndSign(name, password, []sdk.Msg{msg})
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	resp, err := cliCtx.BroadcastTxCommit(msgBytes)
	if err != nil {
		rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	rest.PostProcessResponse(w, cliCtx, resp)
}

// --------------------------------------------------------------------------------------
// Mint NFT

type MintReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID    string `json:"token_id"`
	TokenDenom string `json:"token_denom"`
	TokenURI   string `json:"token_uri"`
}

func mintHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req MintReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := nft.NewMsgMintNFT(owner, owner, req.TokenID, req.TokenDenom, req.TokenURI)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Transfer NFT

type TransferReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID    string `json:"token_id"`
	TokenDenom string `json:"token_denom"`
	Recipient  string `json:"recipient"`
}

func transferHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := nft.NewMsgTransferNFT(owner, recipient, req.TokenDenom, req.TokenID)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Put NFT on Market

type PutOnMarketReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID     string `json:"token_id"`
	Beneficiary string `json:"beneficiary"`
	Price       string `json:"price"`
}

func putOnMarketHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutOnMarketReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		beneficiary, err := sdk.AccAddressFromBech32(req.Beneficiary)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		coins, err := sdk.ParseCoins(req.Price)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the message
		msg := types.NewMsgPutOnMarketNFT(owner, beneficiary, req.TokenID, coins)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Buy NFT

type BuyReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID     string `json:"token_id"`
	Beneficiary string `json:"beneficiary"`
	Commission  string `json:"commission,omitempty"`
}

func buyHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BuyReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		beneficiary, err := sdk.AccAddressFromBech32(req.Beneficiary)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgBuyNFT(owner, beneficiary, req.TokenID, req.Commission)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Buy NFT

type UpdateParamsReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID     string `json:"token_id"`
	Price       string `json:"price,omitempty"`
	TokenName   string `json:"token_name,omitempty"`
	Image       string `json:"image,omitempty"`
	TokenUri    string `json:"token_uri,omitempty"`
	Description string `json:"description,omitempty"`
}

func updateParamsHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req UpdateParamsReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		params := make([]types.NFTParam, 0)
		if req.Price != "" {
			params = append(params, types.NFTParam{Key: types.FlagParamPrice, Value: req.Price})
		}
		if req.Image != "" {
			params = append(params, types.NFTParam{Key: types.FlagParamImage, Value: req.Image})
		}
		if req.Description != "" {
			params = append(params, types.NFTParam{Key: types.FlagParamDescription, Value: req.Description})
		}
		if req.TokenName != "" {
			params = append(params, types.NFTParam{Key: types.FlagParamTokenName, Value: req.TokenName})
		}
		if req.TokenUri != "" {
			params = append(params, types.NFTParam{Key: types.FlagParamTokenURI, Value: req.TokenUri})
		}

		// create the message
		msg := types.NewMsgUpdateNFTParams(owner, req.TokenID, params)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Remove NFT from market

type RemoveNFTFromMarketReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID string `json:"token_id"`
}

func removeNFTFromMarketHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveNFTFromMarketReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := types.NewMsgRemoveNFTFromMarket(owner, req.TokenID)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Create FT

// --------------------------------------------------------------------------------------
//
// FT Handlers
//
// --------------------------------------------------------------------------------------

type CreateFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

func createFTHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req CreateFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := types.NewMsgCreateFungibleToken(owner, req.Denom, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Transfer FT

type TransferFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	Denom     string `json:"denom"`
	Amount    int64  `json:"amount"`
	Recipient string `json:"recipient"`
}

func transferFTHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req TransferFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		recipient, err := sdk.AccAddressFromBech32(req.Recipient)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgTransferFungibleTokens(owner, recipient, req.Denom, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Burn FT

type BurnFTReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	Denom  string `json:"denom"`
	Amount int64  `json:"amount"`
}

func burnFTHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BurnFTReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := types.NewMsgBurnFungibleTokens(owner, req.Denom, req.Amount)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
//
// Auction Handlers
//
// --------------------------------------------------------------------------------------

// --------------------------------------------------------------------------------------
// Put NFT on Auction

type PutOnAuctionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID      string `json:"token_id"`
	Beneficiary  string `json:"beneficiary"`
	OpeningPrice string `json:"opening_price"`
	BuyoutPrice  string `json:"buyout_price,omitempty"`
	Duration     string `json:"duration"`
}

func putOnAuctionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req PutOnAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		beneficiary, err := sdk.AccAddressFromBech32(req.Beneficiary)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		coins, err := sdk.ParseCoins(req.OpeningPrice)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		buyout, err := sdk.ParseCoins(req.BuyoutPrice)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		dur, err := time.ParseDuration(req.Duration)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		// create the message
		msg := types.NewMsgPutNFTOnAuction(owner, beneficiary, req.TokenID, coins, buyout, time.Now().UTC().Add(dur))
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Remove NFT from auction

type RemoveNFTFromAuctionReq RemoveNFTFromMarketReq

func removeNFTFromAuctionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req RemoveNFTFromAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := types.NewMsgRemoveNFTFromAuction(owner, req.TokenID)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Finish NFT auction

type FinishAuctionReq RemoveNFTFromMarketReq

func finishAuctionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req FinishAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		// create the message
		msg := types.NewMsgFinishAuction(owner, req.TokenID)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Bid on auction NFT

type BidOnAuctionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID     string `json:"token_id"`
	Beneficiary string `json:"beneficiary"`
	Bid         string `json:"bid"`
	Commission  string `json:"commission,omitempty"`
}

func bidOnAuctionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BidOnAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		beneficiary, err := sdk.AccAddressFromBech32(req.Beneficiary)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		bid, err := sdk.ParseCoins(req.Bid)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgMakeBidOnAuction(owner, beneficiary, req.TokenID, bid, req.Commission)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
// Buyout auction NFT

type BuyoutAuctionReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID     string `json:"token_id"`
	Beneficiary string `json:"beneficiary"`
	Commission  string `json:"commission,omitempty"`
}

func buyoutAuctionHandler(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req BuyoutAuctionReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			rest.WriteErrorResponse(w, http.StatusBadRequest, "failed to parse request")
			return
		}

		baseReq := req.BaseReq.Sanitize()
		if !baseReq.ValidateBasic(w) {
			return
		}

		owner, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}
		cliCtx.FromName = req.Name
		cliCtx.FromAddress = owner

		beneficiary, err := sdk.AccAddressFromBech32(req.Beneficiary)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgBuyOutOnAuction(owner, beneficiary, req.TokenID, req.Commission)
		err = msg.ValidateBasic()
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		broadcastTransaction(cliCtx, w, msg, req.BaseReq, req.Name, req.Password)
	}
}

// --------------------------------------------------------------------------------------
//
// Query Handlers
//
// --------------------------------------------------------------------------------------

func nftsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/nfts", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func nftHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nftID := vars[restName]
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/nft/%s", storeName, nftID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func fungibleTokensHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/fungible_tokens", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func fungibleTokenHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		ftName := vars[restName]
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/fungible_token/%s", storeName, ftName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func auctionLotsHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/auction_lots", storeName), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}
		rest.PostProcessResponse(w, cliCtx, res)
	}
}

func auctionLotHandler(cliCtx context.CLIContext, storeName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		nftID := vars[restName]
		res, _, err := cliCtx.QueryWithData(fmt.Sprintf("custom/%s/auction_lot/%s", storeName, nftID), nil)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusNotFound, err.Error())
			return
		}

		rest.PostProcessResponse(w, cliCtx, res)
	}
}
