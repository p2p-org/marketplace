package rest

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/corestario/marketplace/x/nftIBC/types"
	"github.com/cosmos/cosmos-sdk/client/context"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/rest"
	authclient "github.com/cosmos/cosmos-sdk/x/auth/client"
)

func registerTxRoutes(cliCtx context.CLIContext, r *mux.Router) {
	r.HandleFunc(fmt.Sprintf("/ibc/ports/{%s}/channels/{%s}/transferNFT", RestPortID, RestChannelID), transferNFTHandlerFn(cliCtx)).Methods("POST")
}

// transferHandlerFn implements a transfer handler
//
// @Summary Transfer token
// @Tags IBC
// @Accept  json
// @Produce  json
// @Param port-id path string true "Port ID"
// @Param channel-id path string true "Channel ID"
// @Param body body rest.TransferTxReq true "Transfer token request body"
// @Success 200 {object} PostTransfer "OK"
// @Failure 400 {object} rest.ErrorResponse "Invalid port id or channel id"
// @Failure 500 {object} rest.ErrorResponse "Internal Server Error"
// @Router /ibc/ports/{port-id}/channels/{channel-id}/transfer [post]
func transferNFTHandlerFn(cliCtx context.CLIContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		portID := vars[RestPortID]
		channelID := vars[RestChannelID]

		var req TransferNFTTxReq
		if !rest.ReadRESTReq(w, r, cliCtx.Codec, &req) {
			return
		}

		req.BaseReq = req.BaseReq.Sanitize()
		if !req.BaseReq.ValidateBasic(w) {
			return
		}

		fromAddr, err := sdk.AccAddressFromBech32(req.BaseReq.From)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		receiver, err := sdk.AccAddressFromBech32(req.Receiver)
		if err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		// create the message
		msg := types.NewMsgTransferNFT(
			portID,
			channelID,
			req.DestHeight,
			fromAddr,
			receiver,
			req.ID,
			req.Denom,
		)

		if err := msg.ValidateBasic(); err != nil {
			rest.WriteErrorResponse(w, http.StatusBadRequest, err.Error())
			return
		}

		authclient.WriteGenerateStdTxResponse(w, cliCtx, req.BaseReq, []sdk.Msg{msg})
	}
}
