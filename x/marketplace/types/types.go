package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/nft/exported"
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
	ID                string         `json:"id"`
	Denom             string         `json:"denom"`
	Owner             sdk.AccAddress `json:"owner"`
	Price             sdk.Coins      `json:"price"`
	Status            NFTStatus      `json:"status"`
	SellerBeneficiary sdk.AccAddress `json:"seller_beneficiary"`
	TimeCreated       time.Time      `json:"time_created"`
	Offers            []*Offer       `json:"offers"`
}

func NewNFT(id string, denom string, owner sdk.AccAddress, price sdk.Coins) *NFT {
	return &NFT{
		ID:          id,
		Owner:       owner,
		Denom:       denom,
		Price:       price,
		TimeCreated: time.Now().UTC(),
	}
}

func (m NFT) String() string {
	var offers []string
	for _, offer := range m.Offers {
		offerJSON, _ := json.Marshal(offer)
		offers = append(offers, string(offerJSON))
	}
	return strings.TrimSpace(fmt.Sprintf(`ID: %s
Owner: %s
Denom: %s
Price: %s
Status: %v
SellerBeneficiary: %s
TimeCreated: %v
Offers: %v`, m.ID, m.Owner, m.Denom, m.Price, m.Status, m.SellerBeneficiary, m.TimeCreated, offers))
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

func (m *NFT) AddOffer(offer *Offer) {
	m.Offers = append(m.Offers, offer)
}

func (m *NFT) GetOffer(id string) (*Offer, bool) {
	for _, offer := range m.Offers {
		offer := offer
		if offer.ID == id {
			return offer, true
		}
	}

	return nil, false
}

func (m *NFT) RemoveOffer(offerID string, buyer sdk.AccAddress) bool {
	for k, offer := range m.Offers {
		k, offer := k, offer
		if offer.ID == offerID && offer.Buyer.Equals(buyer) {
			// this is common sliceTrick for slice of pointers
			// needed for avoid memory leak
			if k < len(m.Offers)-1 {
				copy(m.Offers[k:], m.Offers[k+1:])
			}
			m.Offers[len(m.Offers)-1] = nil
			m.Offers = m.Offers[:len(m.Offers)-1]
			return true
		}
	}

	return false
}

type Offer struct {
	ID                    string         `json:"id"`
	Buyer                 sdk.AccAddress `json:"buyer"`
	Price                 sdk.Coins      `json:"price"`
	BuyerBeneficiary      sdk.AccAddress `json:"buyer_beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission"`
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

// copy of data got from exported/nft interface
type NFTMetaData struct {
	ID       string         `json:"id"`
	Owner    sdk.AccAddress `json:"owner"`
	TokenURI string         `json:"token_uri"`
}

func (d NFTMetaData) String() string {
	return strings.TrimSpace(fmt.Sprintf(`ID: %s
Owner: %s
TokenURI: %s`, d.ID, d.Owner, d.TokenURI))
}

func NewNFTMetaData(token exported.NFT) *NFTMetaData {
	return &NFTMetaData{
		ID:       token.GetID(),
		Owner:    token.GetOwner(),
		TokenURI: token.GetTokenURI(),
	}
}

type NFTInfo struct {
	MPNFTInfo    *NFT `json:"nft_mp_info"`
	*NFTMetaData `json:"nft_meta_data"`
}

func (i NFTInfo) String() string {
	return strings.TrimSpace(fmt.Sprintf(`MarketPlaceInfo: %s
MetaData: %s`, i.MPNFTInfo, i.NFTMetaData))
}

func NewNFTInfo(nft *NFT, token exported.NFT) *NFTInfo {
	return &NFTInfo{
		MPNFTInfo:   nft,
		NFTMetaData: NewNFTMetaData(token),
	}
}
