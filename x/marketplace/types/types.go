package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	xnft "github.com/cosmos/cosmos-sdk/x/nft"
)

type FungibleToken struct {
	Denom          string         `json:"denom"`
	EmissionAmount int64          `json:"emission_amount"`
	Creator        sdk.AccAddress `json:"creator"`
}

func (c FungibleToken) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Creator: %s
EmissionAmount: %d
Denom: %s`, c.Creator.String(), c.EmissionAmount, c.Denom))
}

type NFT struct {
	xnft.BaseNFT      `json:"nft"`
	Price             sdk.Coins      `json:"price"`
	Status            NFTStatus      `json:"status"`
	SellerBeneficiary sdk.AccAddress `json:"seller_beneficiary"`
	TimeCreated       time.Time      `json:"time_created"`
}

func NewNFT(nft xnft.BaseNFT, price sdk.Coins) *NFT {
	return &NFT{
		BaseNFT:     nft,
		Price:       price,
		TimeCreated: time.Now().UTC(),
	}
}

func (m NFT) String() string {
	return strings.TrimSpace(fmt.Sprintf(`Owner: %s
NFT: %s
Price: %s
Status: %v`, m.GetOwner(), m.BaseNFT, m.Price, m.Status))
}

func (m *NFT) GetPrice() sdk.Coins {
	return m.Price
}

func (m *NFT) SetPrice(price sdk.Coins) {
	m.Price = price
}

func (m *NFT) IsOnSale() bool {
	return m.Status > NFTStatusDefault && m.Status < NFTStatusUndefined
}

func (m *NFT) IsOnMarket() bool {
	return m.Status == NFTStatusOnMarket
}

func (m *NFT) IsOnAuction() bool {
	return m.Status == NFTStatusOnAuction
}

func (m *NFT) SetStatus(status NFTStatus) {
	m.Status = status
}

func (m *NFT) SetSellerBeneficiary(addr sdk.AccAddress) {
	m.SellerBeneficiary = addr
}

func (m *NFT) GetTimeCreated() time.Time {
	return m.TimeCreated
}

type AuctionBid struct {
	Bidder                sdk.AccAddress `json:"bidder"`            // account address that made the bid
	BuyerBeneficiary      sdk.AccAddress `json:"buyer_beneficiary"` // account address that will be the beneficiary of the purchase
	BeneficiaryCommission string         `json:"beneficiary_commission"`
	Bid                   sdk.Coins      `json:"bid"`
	TimeCreated           time.Time      `json:"time_created"`
}

type AuctionLot struct {
	NFTID          string      `json:"nft_id"`
	LastBid        *AuctionBid `json:"last_bid"`
	OpeningPrice   sdk.Coins   `json:"opening_price"`
	BuyoutPrice    sdk.Coins   `json:"buyout_price"`
	ExpirationTime time.Time   `json:"expiration_time"`
}

func NewAuctionLot(id string, openingPrice, buyoutPrice sdk.Coins, expTime time.Time) *AuctionLot {
	return &AuctionLot{
		NFTID:          id,
		OpeningPrice:   openingPrice,
		BuyoutPrice:    buyoutPrice,
		ExpirationTime: expTime,
	}
}

func NewAuctionBid(bidder, beneficiary sdk.AccAddress, price sdk.Coins, commission string) *AuctionBid {
	return &AuctionBid{
		Bidder:                bidder,
		BuyerBeneficiary:      beneficiary,
		Bid:                   price,
		TimeCreated:           time.Now().UTC(),
		BeneficiaryCommission: commission,
	}
}

func (lot *AuctionLot) SetLastBid(bid *AuctionBid) {
	lot.LastBid = bid
}

func (lot AuctionLot) String() string {
	base := strings.TrimSpace(fmt.Sprintf(`NFT: %s
OpeningPrice: %v
ExpirationTime: %v
`, lot.NFTID, lot.OpeningPrice, lot.ExpirationTime))

	if lot.BuyoutPrice.IsZero() {
		base += strings.TrimSpace(fmt.Sprintf(`
BuyoutPrice: %v`, nil))
	} else {
		base += strings.TrimSpace(fmt.Sprintf(`
BuyoutPrice: %v`, lot.BuyoutPrice))
	}

	if lot.LastBid != nil {
		base += strings.TrimSpace(fmt.Sprintf(`
LastBid: %v
TimeOfLastBid: %v`, lot.LastBid.Bid, lot.LastBid.TimeCreated))
	}

	return base
}
