package marketplace

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/exported"
)

type AuctionLotRecord struct {
	TokenID        string
	ExpirationTime time.Time
}

type AccountResponse struct {
	Result struct {
		Value struct {
			AccountNumber uint64 `json:"account_number,string"`
			Sequence      uint64 `json:"sequence,string"`
		} `json:"value"`
	} `json:"result"`
}

type BaseReq struct {
	From          string       `json:"from,omitempty"`
	Memo          string       `json:"memo,omitempty"`
	ChainID       string       `json:"chain_id,omitempty"`
	AccountNumber uint64       `json:"account_number,string,omitempty"`
	Sequence      uint64       `json:"sequence,string,omitempty"`
	Fees          sdk.Coins    `json:"fees,omitempty"`
	GasPrices     sdk.DecCoins `json:"gas_prices,omitempty"`
	Gas           string       `json:"gas,omitempty"`
	GasAdjustment string       `json:"gas_adjustment,omitempty"`
	Simulate      bool         `json:"simulate,omitempty"`
}

type FinishAuctionReq struct {
	BaseReq BaseReq `json:"base_req"`

	Name     string `json:"name"`
	Password string `json:"password"`

	TokenID string `json:"token_id"`
}

func (k *Keeper) SendFinish(id string, acc exported.Account) error {
	far := FinishAuctionReq{
		BaseReq: BaseReq{
			Sequence:      acc.GetSequence(),
			ChainID:       k.config.ChainName,
			AccountNumber: acc.GetAccountNumber(),
			From:          k.config.FinishingAccountAddr,
		},
		Name:     k.config.FinishingAccountName,
		Password: k.config.FinishingAccountPass,

		TokenID: id,
	}

	ba, err := json.Marshal(&far)
	if err != nil {
		return err
	}

	finishAuctionAddr := fmt.Sprintf("http://%s:%d/marketplace/finish_auction",
		k.config.FinishAuctionHost, k.config.FinishAuctionPort)

	buf := bytes.NewBuffer(ba)
	req, err := http.NewRequest(http.MethodPut, finishAuctionAddr, buf)
	if err != nil {
		return err
	}

	resp, err := k.httpCli.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		rsp, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return fmt.Errorf("http error, status code: %v, %v", resp.StatusCode, string(rsp))
	}

	return nil
}
