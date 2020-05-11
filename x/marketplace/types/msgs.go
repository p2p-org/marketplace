package types

import (
	"fmt"
	"strings"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// --------------------------------------------------------------------------
//
// MsgPutNFTOnMarket
//
// --------------------------------------------------------------------------

type MsgPutNFTOnMarket struct {
	Owner sdk.AccAddress `json:"owner"`
	// Beneficiary is the cosmos user who gets the commission for this transaction.
	Beneficiary sdk.AccAddress `json:"beneficiary"`
	TokenID     string         `json:"token_id"`
	Price       sdk.Coins      `json:"price"`
}

func NewMsgPutOnMarketNFT(owner, beneficiary sdk.AccAddress, tokenID string, price sdk.Coins) *MsgPutNFTOnMarket {
	return &MsgPutNFTOnMarket{
		Owner:       owner,
		TokenID:     tokenID,
		Price:       price,
		Beneficiary: beneficiary,
	}
}

// Route should return the name of the module
func (m MsgPutNFTOnMarket) Route() string { return RouterKey }

// Type should return the action
func (m MsgPutNFTOnMarket) Type() string { return "put_on_market_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgPutNFTOnMarket) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrUnknownAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	if m.Price.IsZero() || m.Price.IsAnyNegative() {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgPutNFTOnMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgPutNFTOnMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgPutNFTOnMarket
//
// --------------------------------------------------------------------------

type MsgRemoveNFTFromMarket struct {
	Owner   sdk.AccAddress `json:"owner"`
	TokenID string         `json:"token_id"`
}

func NewMsgRemoveNFTFromMarket(owner sdk.AccAddress, tokenID string) *MsgRemoveNFTFromMarket {
	return &MsgRemoveNFTFromMarket{
		Owner:   owner,
		TokenID: tokenID,
	}
}

// Route should return the name of the module
func (m MsgRemoveNFTFromMarket) Route() string { return RouterKey }

// Type should return the action
func (m MsgRemoveNFTFromMarket) Type() string { return "remove_from_market_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgRemoveNFTFromMarket) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgRemoveNFTFromMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgRemoveNFTFromMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgBuyNFT
//
// --------------------------------------------------------------------------

type MsgBuyNFT struct {
	Buyer sdk.AccAddress `json:"buyer"`
	// Beneficiary is the cosmos user who gets the commission for this transaction.
	Beneficiary           sdk.AccAddress `json:"beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenID               string         `json:"token_id"`
}

func NewMsgBuyNFT(owner, beneficiary sdk.AccAddress, tokenID string, commission string) *MsgBuyNFT {
	return &MsgBuyNFT{
		Buyer:                 owner,
		Beneficiary:           beneficiary,
		BeneficiaryCommission: commission,
		TokenID:               tokenID,
	}
}

// Route should return the name of the module
func (m MsgBuyNFT) Route() string { return RouterKey }

// Type should return the action
func (m MsgBuyNFT) Type() string { return "buy_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgBuyNFT) ValidateBasic() error {
	if m.Buyer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBuyNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBuyNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Buyer}
}

// --------------------------------------------------------------------------
//
// MsgCreateFungibleToken
//
// --------------------------------------------------------------------------

type MsgCreateFungibleToken struct {
	Creator sdk.AccAddress `json:"creator"`
	Denom   string         `json:"denom"`
	Amount  int64          `json:"amount"`
}

func NewMsgCreateFungibleToken(creator sdk.AccAddress, denom string, amount int64) *MsgCreateFungibleToken {
	return &MsgCreateFungibleToken{
		Creator: creator,
		Denom:   denom,
		Amount:  amount,
	}
}

func (m MsgCreateFungibleToken) Route() string { return RouterKey }

func (m MsgCreateFungibleToken) Type() string { return "create_fungible_token" }

func (m MsgCreateFungibleToken) ValidateBasic() error {
	if m.Creator.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.Denom) < MinDenomLength || len(m.Denom) > MaxDenomLength {
		return sdkerrors.ErrUnknownRequest
	}
	if m.Amount <= 0 {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

func (m MsgCreateFungibleToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgCreateFungibleToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Creator}
}

// --------------------------------------------------------------------------
//
// MsgTransferFungibleTokens
//
// --------------------------------------------------------------------------

type MsgTransferFungibleTokens struct {
	Owner     sdk.AccAddress `json:"owner"`
	Recipient sdk.AccAddress `json:"recipient"`
	Denom     string         `json:"denom"`
	Amount    int64          `json:"amount"`
}

func NewMsgTransferFungibleTokens(owner, recipient sdk.AccAddress, denom string, amount int64) *MsgTransferFungibleTokens {
	return &MsgTransferFungibleTokens{
		Owner:     owner,
		Recipient: recipient,
		Denom:     denom,
		Amount:    amount,
	}
}

func (m MsgTransferFungibleTokens) Route() string { return RouterKey }

func (m MsgTransferFungibleTokens) Type() string { return "transfer_coins" }

func (m MsgTransferFungibleTokens) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if m.Recipient.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.Denom) < MinDenomLength || len(m.Denom) > MaxDenomLength {
		return sdkerrors.ErrUnknownRequest
	}
	if m.Amount <= 0 {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

func (m MsgTransferFungibleTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgTransferFungibleTokens) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgUpdateNFTParams
//
// --------------------------------------------------------------------------

type NFTParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type NFTParams []NFTParam

func (p NFTParams) String() string {
	out := make([]string, 0)
	for _, v := range p {
		v := v
		out = append(out,
			strings.TrimSpace(fmt.Sprintf(`Key: %s; Value: %s`, v.Key, v.Value)))
	}
	return strings.Join(out, "\n")
}

type MsgUpdateNFTParams struct {
	Owner   sdk.AccAddress `json:"owner"`
	Params  NFTParams      `json:"params"`
	TokenID string         `json:"token_id"`
}

func NewMsgUpdateNFTParams(owner sdk.AccAddress, id string, params []NFTParam) *MsgUpdateNFTParams {
	return &MsgUpdateNFTParams{
		Owner:   owner,
		Params:  params,
		TokenID: id,
	}
}

// Route should return the name of the module
func (m MsgUpdateNFTParams) Route() string { return RouterKey }

// Type should return the action
func (m MsgUpdateNFTParams) Type() string { return "update_nft_params" }

// ValidateBasic runs stateless checks on the message
func (m MsgUpdateNFTParams) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if m.TokenID == "" {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgUpdateNFTParams) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgUpdateNFTParams) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgBurnFungibleTokens
//
// --------------------------------------------------------------------------

type MsgBurnFungibleTokens struct {
	Owner  sdk.AccAddress `json:"owner"`
	Denom  string         `json:"denom"`
	Amount int64          `json:"amount"`
}

func NewMsgBurnFungibleTokens(owner sdk.AccAddress, denom string, amount int64) *MsgBurnFungibleTokens {
	return &MsgBurnFungibleTokens{
		Owner:  owner,
		Denom:  denom,
		Amount: amount,
	}
}

func (m MsgBurnFungibleTokens) Route() string { return RouterKey }

func (m MsgBurnFungibleTokens) Type() string { return "burn_coins" }

func (m MsgBurnFungibleTokens) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.Denom) < 3 || len(m.Denom) > 16 {
		return sdkerrors.ErrUnknownRequest
	}
	if m.Amount <= 0 {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

func (m MsgBurnFungibleTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgBurnFungibleTokens) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgPutNFTOnAuction
//
// --------------------------------------------------------------------------

type MsgPutNFTOnAuction struct {
	Owner sdk.AccAddress `json:"owner"`
	// Beneficiary is the cosmos user who gets the commission for this transaction.
	Beneficiary  sdk.AccAddress `json:"beneficiary"`
	TokenID      string         `json:"token_id"`
	OpeningPrice sdk.Coins      `json:"opening_price"`
	BuyoutPrice  sdk.Coins      `json:"buyout_price"`
	TimeToSell   time.Time      `json:"time_to_sell"`
}

func NewMsgPutNFTOnAuction(owner, beneficiary sdk.AccAddress, tokenID string,
	openingPrice, buyoutPrice sdk.Coins, timeToSell time.Time) *MsgPutNFTOnAuction {
	return &MsgPutNFTOnAuction{
		Owner:        owner,
		TokenID:      tokenID,
		OpeningPrice: openingPrice,
		BuyoutPrice:  buyoutPrice,
		Beneficiary:  beneficiary,
		TimeToSell:   timeToSell,
	}
}

// Route should return the name of the module
func (m MsgPutNFTOnAuction) Route() string { return RouterKey }

// Type should return the action
func (m MsgPutNFTOnAuction) Type() string { return "put_on_auction_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgPutNFTOnAuction) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	if m.OpeningPrice.IsZero() || m.OpeningPrice.IsAnyNegative() {
		return sdkerrors.ErrUnknownRequest
	}
	if m.TimeToSell.IsZero() {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgPutNFTOnAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgPutNFTOnAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgRemoveNFTFromAuction
//
// --------------------------------------------------------------------------

type MsgRemoveNFTFromAuction MsgRemoveNFTFromMarket

func NewMsgRemoveNFTFromAuction(owner sdk.AccAddress, tokenID string) *MsgRemoveNFTFromAuction {
	return &MsgRemoveNFTFromAuction{
		Owner:   owner,
		TokenID: tokenID,
	}
}

// Route should return the name of the module
func (m MsgRemoveNFTFromAuction) Route() string { return RouterKey }

// Type should return the action
func (m MsgRemoveNFTFromAuction) Type() string { return "remove_from_auction_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgRemoveNFTFromAuction) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgRemoveNFTFromAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgRemoveNFTFromAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgMakeBidOnAuction
//
// --------------------------------------------------------------------------

type MsgMakeBidOnAuction struct {
	Bidder sdk.AccAddress `json:"bidder"`
	// Beneficiary is the cosmos user who gets the commission for this transaction.
	BuyerBeneficiary      sdk.AccAddress `json:"buyer_beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenID               string         `json:"token_id"`
	Bid                   sdk.Coins      `json:"bid"`
}

func NewMsgMakeBidOnAuction(bidder, buyerBeneficiary sdk.AccAddress, tokenID string, bid sdk.Coins, commission string) *MsgMakeBidOnAuction {
	return &MsgMakeBidOnAuction{
		Bidder:                bidder,
		BuyerBeneficiary:      buyerBeneficiary,
		TokenID:               tokenID,
		Bid:                   bid,
		BeneficiaryCommission: commission,
	}
}

// Route should return the name of the module
func (m MsgMakeBidOnAuction) Route() string { return RouterKey }

// Type should return the action
func (m MsgMakeBidOnAuction) Type() string { return "make_bid_on_auction_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgMakeBidOnAuction) ValidateBasic() error {
	if m.Bidder.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgMakeBidOnAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgMakeBidOnAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Bidder}
}

// --------------------------------------------------------------------------
//
// MsgFinishAuction
//
// --------------------------------------------------------------------------

type MsgFinishAuction MsgRemoveNFTFromMarket

func NewMsgFinishAuction(owner sdk.AccAddress, tokenID string) *MsgFinishAuction {
	return &MsgFinishAuction{
		Owner:   owner,
		TokenID: tokenID,
	}
}

// Route should return the name of the module
func (m MsgFinishAuction) Route() string { return RouterKey }

// Type should return the action
func (m MsgFinishAuction) Type() string { return "finish_auction_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgFinishAuction) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgFinishAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgFinishAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgBuyoutOnAuction
//
// --------------------------------------------------------------------------

type MsgBuyoutOnAuction struct {
	Buyer                 sdk.AccAddress `json:"buyer"`
	BuyerBeneficiary      sdk.AccAddress `json:"buyer_beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenID               string         `json:"token_id"`
}

func NewMsgBuyOutOnAuction(bidder, buyerBeneficiary sdk.AccAddress, tokenID string, commission string) *MsgBuyoutOnAuction {
	return &MsgBuyoutOnAuction{
		Buyer:                 bidder,
		BuyerBeneficiary:      buyerBeneficiary,
		TokenID:               tokenID,
		BeneficiaryCommission: commission,
	}
}

// Route should return the name of the module
func (m MsgBuyoutOnAuction) Route() string { return RouterKey }

// Type should return the action
func (m MsgBuyoutOnAuction) Type() string { return "buyout_on_auction_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgBuyoutOnAuction) ValidateBasic() error {
	if m.Buyer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBuyoutOnAuction) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBuyoutOnAuction) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Buyer}
}

// --------------------------------------------------------------------------
//
// MsgBatchTransfer
//
// --------------------------------------------------------------------------

type MsgBatchTransfer struct {
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	TokenIDs  []string       `json:"token_ids"`
}

func NewMsgBatchTransfer(sender, recipient sdk.AccAddress, tokenIDs []string) *MsgBatchTransfer {
	return &MsgBatchTransfer{
		Sender:    sender,
		Recipient: recipient,
		TokenIDs:  tokenIDs,
	}
}

// Route should return the name of the module
func (m MsgBatchTransfer) Route() string { return RouterKey }

// Type should return the action
func (m MsgBatchTransfer) Type() string { return "batch_transfer" }

// ValidateBasic runs stateless checks on the message
func (m MsgBatchTransfer) ValidateBasic() error {
	if m.Sender.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenIDs) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	for _, tokenID := range m.TokenIDs {
		if len(tokenID) > MaxTokenIDLength {
			return sdkerrors.ErrUnknownRequest
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBatchTransfer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBatchTransfer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

// --------------------------------------------------------------------------
//
// MsgBatchPutOnMarket
//
// --------------------------------------------------------------------------

type MsgBatchPutOnMarket struct {
	Owner       sdk.AccAddress `json:"owner"`
	Beneficiary sdk.AccAddress `json:"beneficiary"`
	TokenIDs    []string       `json:"token_ids"`
	TokenPrices []sdk.Coins    `json:"token_prices"`
}

func NewMsgBatchPutOnMarket(owner, beneficiary sdk.AccAddress, tokenIDs []string, tokenPrices []sdk.Coins) *MsgBatchPutOnMarket {
	return &MsgBatchPutOnMarket{
		Owner:       owner,
		Beneficiary: beneficiary,
		TokenIDs:    tokenIDs,
		TokenPrices: tokenPrices,
	}
}

// Route should return the name of the module
func (m MsgBatchPutOnMarket) Route() string { return RouterKey }

// Type should return the action
func (m MsgBatchPutOnMarket) Type() string { return "batch_put_on_market" }

// ValidateBasic runs stateless checks on the message
func (m MsgBatchPutOnMarket) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}

	if len(m.TokenIDs) == 0 {
		return sdkerrors.ErrUnknownRequest
	}

	if len(m.TokenPrices) == 0 {
		return sdkerrors.ErrUnknownRequest
	}

	if len(m.TokenIDs) != len(m.TokenPrices) {
		return sdkerrors.ErrUnknownRequest
	}
	for _, tokenID := range m.TokenIDs {
		tokenID := tokenID
		if len(tokenID) > MaxTokenIDLength {
			return sdkerrors.ErrUnknownRequest
		}
	}

	for _, price := range m.TokenPrices {
		price := price
		if price.IsZero() || price.IsAnyNegative() {
			return sdkerrors.ErrUnknownRequest
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBatchPutOnMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBatchPutOnMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgBatchRemoveFromMarket
//
// --------------------------------------------------------------------------

type MsgBatchRemoveFromMarket struct {
	Owner    sdk.AccAddress `json:"owner"`
	TokenIDs []string       `json:"token_ids"`
}

func NewMsgBatchRemoveFromMarket(owner sdk.AccAddress, tokenIDs []string) *MsgBatchRemoveFromMarket {
	return &MsgBatchRemoveFromMarket{
		Owner:    owner,
		TokenIDs: tokenIDs,
	}
}

// Route should return the name of the module
func (m MsgBatchRemoveFromMarket) Route() string { return RouterKey }

// Type should return the action
func (m MsgBatchRemoveFromMarket) Type() string { return "batch_remove_from_market" }

// ValidateBasic runs stateless checks on the message
func (m MsgBatchRemoveFromMarket) ValidateBasic() error {
	if m.Owner.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenIDs) == 0 {
		return sdkerrors.ErrUnknownRequest
	}

	for _, tokenID := range m.TokenIDs {
		if len(tokenID) > MaxTokenIDLength {
			return sdkerrors.ErrUnknownRequest
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBatchRemoveFromMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBatchRemoveFromMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgBatchBuyOnMarket
//
// --------------------------------------------------------------------------

type MsgBatchBuyOnMarket struct {
	Buyer                 sdk.AccAddress `json:"owner"`
	Beneficiary           sdk.AccAddress `json:"beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenIDs              []string       `json:"token_ids"`
}

func NewMsgBatchBuyOnMarket(buyer, beneficiary sdk.AccAddress, commission string, tokenIDs []string) *MsgBatchBuyOnMarket {
	return &MsgBatchBuyOnMarket{
		Buyer:                 buyer,
		Beneficiary:           beneficiary,
		BeneficiaryCommission: commission,
		TokenIDs:              tokenIDs,
	}
}

// Route should return the name of the module
func (m MsgBatchBuyOnMarket) Route() string { return RouterKey }

// Type should return the action
func (m MsgBatchBuyOnMarket) Type() string { return "batch_buy_on_market" }

// ValidateBasic runs stateless checks on the message
func (m MsgBatchBuyOnMarket) ValidateBasic() error {
	if m.Buyer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenIDs) == 0 {
		return sdkerrors.ErrUnknownRequest
	}

	for _, tokenID := range m.TokenIDs {
		if len(tokenID) > MaxTokenIDLength {
			return sdkerrors.ErrUnknownRequest
		}
	}

	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgBatchBuyOnMarket) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgBatchBuyOnMarket) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Buyer}
}

// --------------------------------------------------------------------------
//
// MsgMakeOffer
//
// --------------------------------------------------------------------------

type MsgMakeOffer struct {
	Buyer                 sdk.AccAddress `json:"buyer"`
	Price                 sdk.Coins      `json:"price"`
	BuyerBeneficiary      sdk.AccAddress `json:"buyer_beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenID               string         `json:"token_id"`
}

func NewMsgMakeOffer(bidder, buyerBeneficiary sdk.AccAddress, price sdk.Coins, tokenID string, commission string) *MsgMakeOffer {
	return &MsgMakeOffer{
		Buyer:                 bidder,
		Price:                 price,
		BuyerBeneficiary:      buyerBeneficiary,
		TokenID:               tokenID,
		BeneficiaryCommission: commission,
	}
}

// Route should return the name of the module
func (m MsgMakeOffer) Route() string { return RouterKey }

// Type should return the action
func (m MsgMakeOffer) Type() string { return "make_offer" }

// ValidateBasic runs stateless checks on the message
func (m MsgMakeOffer) ValidateBasic() error {
	if m.Buyer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgMakeOffer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgMakeOffer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Buyer}
}

// --------------------------------------------------------------------------
//
// MsgAcceptOffer
//
// --------------------------------------------------------------------------

type MsgAcceptOffer struct {
	Seller                sdk.AccAddress `json:"seller"`
	SellerBeneficiary     sdk.AccAddress `json:"seller_beneficiary"`
	BeneficiaryCommission string         `json:"beneficiary_commission,omitempty"`
	TokenID               string         `json:"token_id"`
	OfferID               string         `json:"offer_id"`
}

func NewMsgAcceptOffer(seller, sellerBeneficiary sdk.AccAddress, tokenID, offerID string, commission string) *MsgAcceptOffer {
	return &MsgAcceptOffer{
		Seller:                seller,
		SellerBeneficiary:     sellerBeneficiary,
		TokenID:               tokenID,
		OfferID:               offerID,
		BeneficiaryCommission: commission,
	}
}

// Route should return the name of the module
func (m MsgAcceptOffer) Route() string { return RouterKey }

// Type should return the action
func (m MsgAcceptOffer) Type() string { return "accept_offer" }

// ValidateBasic runs stateless checks on the message
func (m MsgAcceptOffer) ValidateBasic() error {
	if m.Seller.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgAcceptOffer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgAcceptOffer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Seller}
}

// --------------------------------------------------------------------------
//
// MsgRemoveOffer
//
// --------------------------------------------------------------------------

type MsgRemoveOffer struct {
	Buyer   sdk.AccAddress `json:"buyer"`
	TokenID string         `json:"token_id"`
	OfferID string         `json:"offer_id"`
}

func NewMsgRemoveOffer(buyer sdk.AccAddress, tokenID, offerID string) *MsgRemoveOffer {
	return &MsgRemoveOffer{
		Buyer:   buyer,
		TokenID: tokenID,
		OfferID: offerID,
	}
}

// Route should return the name of the module
func (m MsgRemoveOffer) Route() string { return RouterKey }

// Type should return the action
func (m MsgRemoveOffer) Type() string { return "remove_offer" }

// ValidateBasic runs stateless checks on the message
func (m MsgRemoveOffer) ValidateBasic() error {
	if m.Buyer.Empty() {
		return sdkerrors.ErrInvalidAddress
	}
	if len(m.OfferID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) == 0 {
		return sdkerrors.ErrUnknownRequest
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdkerrors.ErrUnknownRequest
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgRemoveOffer) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgRemoveOffer) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Buyer}
}

type MsgTransferNFTByIBC struct {
	SourcePort    string         `json:"source_port" yaml:"source_port"`       // the port on which the packet will be sent
	SourceChannel string         `json:"source_channel" yaml:"source_channel"` // the channel by which the packet will be sent
	TokenID       string         `json:"token_id" yaml:"token_id"`             // the id of token to be transferred
	Sender        sdk.AccAddress `json:"sender" yaml:"sender"`                 // the sender address
	Receiver      sdk.AccAddress `json:"receiver" yaml:"receiver"`             // the recipient address on the destination chain
	Source        bool           `json:"source" yaml:"source"`                 // indicates if the sending chain is the source chain of the tokens to be transferred
}

// NewMsgTransfer creates a new MsgTransfer instance
func NewMsgTransferNFTByIBC(sourcePort, sourceChannel string, tokenID string, sender, receiver sdk.AccAddress, source bool,
) MsgTransferNFTByIBC {
	return MsgTransferNFTByIBC{
		SourcePort:    sourcePort,
		SourceChannel: sourceChannel,
		TokenID:       tokenID,
		Sender:        sender,
		Receiver:      receiver,
		Source:        source,
	}
}

// Route implements sdk.Msg
func (MsgTransferNFTByIBC) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgTransferNFTByIBC) Type() string {
	return "transferNFT"
}

// ValidateBasic implements sdk.Msg
func (msg MsgTransferNFTByIBC) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgTransferNFTByIBC) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgTransferNFTByIBC) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
