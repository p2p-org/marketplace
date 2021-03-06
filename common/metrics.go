package common

import (
	"github.com/prometheus/client_golang/prometheus"
)

type MsgMetrics struct {
	NumMsgs *prometheus.CounterVec
}

const (
	PrometheusLabelStatus                      = "status"
	PrometheusLabelMsgType                     = "msg_type"
	PrometheusValueReceived                    = "Received"
	PrometheusValueAccepted                    = "Accepted"
	PrometheusValueCommon                      = "Common"
	PrometheusValueMsgMintNFT                  = "MsgMintNFT"
	PrometheusValueMsgBurnNFT                  = "MsgBurnNFT"
	PrometheusValueMsgPutNFTOnMarket           = "MsgPutNFTOnMarket"
	PrometheusValueMsgBuyNFT                   = "MsgBuyNFT"
	PrometheusValueMsgTransferNFT              = "MsgTransferNFT"
	PrometheusValueMsgCreateFungibleToken      = "MsgCreateFungibleToken"
	PrometheusValueMsgTransferFungibleTokens   = "MsgTransferFungibleTokens"
	PrometheusValueMsgUpdateNFTParams          = "MsgUpdateNFTParams"
	PrometheusValueMsgBurnFT                   = "MsgBurnFT"
	PrometheusValueMsgPutNFTOnAuction          = "MsgPutNFTOnAuction"
	PrometheusValueMsgRemoveNFTFromAuction     = "MsgRemoveNFTFromAuction"
	PrometheusValueMsgRemoveNFTFromMarket      = "MsgRemoveNFTFromMarket"
	PrometheusValueMsgFinishAuction            = "MsgFinishAuction"
	PrometheusValueMsgMakeBidOnAuction         = "MsgMakeBidOnAuction"
	PrometheusValueMsgBuyoutFromAuction        = "MsgBuyoutFromAuction"
	PrometheusValueMsgBatchTransfer            = "MsgMsgBatchTransfer"
	PrometheusValueMsgMsgBatchPutOnMarket      = "MsgMsgBatchPutOnMarket"
	PrometheusValueMsgMsgBatchRemoveFromMarket = "MsgMsgBatchRemoveFromMarket"
	PrometheusValueMsgMsgBatchBuyOnMarket      = "MsgMsgBatchBuyOnMarket"
	PrometheusValueMsgMakeOffer                = "MsgMakeOffer"
	PrometheusValueMsgAcceptOffer              = "MsgAcceptOffer"
	PrometheusValueMsgRemoveOffer              = "MsgRemoveOffer"
)

func NewPrometheusMsgMetrics(module string) *MsgMetrics {
	numMsgs := prometheus.NewCounterVec(prometheus.CounterOpts{
		Namespace: "Marketplace",
		Subsystem: module + "_MetricsSubsystem",
		Name:      "NumMsgs",
		Help:      "number of messages since start",
	},
		[]string{PrometheusLabelStatus, PrometheusLabelMsgType},
	)
	prometheus.MustRegister(numMsgs)
	return &MsgMetrics{
		NumMsgs: numMsgs,
	}
}
