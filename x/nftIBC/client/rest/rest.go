package rest

import (
	"github.com/gorilla/mux"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/types/rest"
)

const (
	RestChannelID = "channel-id"
	RestPortID    = "port-id"
)

// RegisterRoutes - Central function to define routes that get registered by the main application
func RegisterRoutes(cliCtx context.CLIContext, r *mux.Router) {
	registerQueryRoutes(cliCtx, r)
	registerTxRoutes(cliCtx, r)
}

type TransferNFTTxReq struct {
	BaseReq    rest.BaseReq `json:"base_req" yaml:"base_req"`
	DestHeight uint64       `json:"dest_height" yaml:"dest_height"`
	ID         string       `json:"id" yaml:"id"`
	Denom      string       `json:"denom" yaml:"denom"`
	Receiver   string       `json:"receiver" yaml:"receiver"`
}
