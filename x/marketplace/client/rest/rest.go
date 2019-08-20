package rest

import (
	"fmt"
	"net/http"

	"github.com/dgamingfoundation/cosmos-sdk/client/flags"

	"github.com/dgamingfoundation/cosmos-sdk/x/auth"

	"github.com/dgamingfoundation/cosmos-sdk/client/context"
	"github.com/dgamingfoundation/marketplace/x/marketplace/types"

	sdk "github.com/dgamingfoundation/cosmos-sdk/types"
	"github.com/dgamingfoundation/cosmos-sdk/types/rest"
	"github.com/dgamingfoundation/cosmos-sdk/x/auth/client/utils"

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

	r.HandleFunc(fmt.Sprintf("/%s/put_on_market", storeName), putOnMarketHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/buy", storeName), buyHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/update_params", storeName), updateParamsHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/create_ft", storeName), createFTHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/transfer_ft", storeName), transferFTHandler(cliCtx)).Methods("PUT")
	r.HandleFunc(fmt.Sprintf("/%s/burn_ft", storeName), burnFTHandler(cliCtx)).Methods("PUT")
}

// --------------------------------------------------------------------------------------
//
// Tx Handler
//
// --------------------------------------------------------------------------------------

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
		utils.GetTxEncoder(cliCtx.Codec), bq.AccountNumber, bq.Sequence, gas, gasAdj,
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
// Transfer NFT

type TransferReq struct {
	BaseReq rest.BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID   string `json:"token_id"`
	Recipient string `json:"recipient"`
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
	Commission  string `json:"commission"`
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
// Create FT

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
