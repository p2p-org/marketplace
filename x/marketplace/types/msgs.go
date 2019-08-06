package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// --------------------------------------------------------------------------
//
// MsgMintNFT
//
// --------------------------------------------------------------------------

type MsgMintNFT struct {
	TokenID     string         `json:"token_id"`
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	TokenURI    string         `json:"token_uri"`
}

func NewMsgMintNFT(
	TokenID string,
	owner sdk.AccAddress,
	name,
	description,
	image,
	tokenURI string,
) *MsgMintNFT {
	return &MsgMintNFT{
		TokenID:     TokenID,
		Owner:       owner,
		Name:        name,
		Description: description,
		Image:       image,
		TokenURI:    tokenURI,
	}
}

// Route should return the name of the module
func (m MsgMintNFT) Route() string { return RouterKey }

// Type should return the action
func (m MsgMintNFT) Type() string { return "mint_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgMintNFT) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if len(m.Name) == 0 || len(m.Description) == 0 {
		return sdk.ErrUnknownRequest("Name and/or Description cannot be empty")
	}
	if len(m.Name) > MaxNameLength {
		return sdk.ErrUnknownRequest("Name has invalid format")
	}
	if len(m.Description) > MaxDescriptionLength {
		return sdk.ErrUnknownRequest("Description has invalid format")
	}
	if len(m.TokenURI) > MaxTokenURILength {
		return sdk.ErrUnknownRequest("TokenURI has invalid format")
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdk.ErrUnknownRequest("TokenID has invalid format")
	}
	if len(m.Image) > MaxImageLength {
		return sdk.ErrUnknownRequest("Image has invalid format")
	}
	if !isTokenURIValid(m.TokenURI) {
		return sdk.ErrUnknownRequest("TokenURI has invalid format")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgMintNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgMintNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

// --------------------------------------------------------------------------
//
// MsgTransferNFT
//
// --------------------------------------------------------------------------

type MsgTransferNFT struct {
	TokenID   string         `json:"token_id"`
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
}

func NewMsgTransferNFT(tokenID string, sender, recipient sdk.AccAddress) *MsgTransferNFT {
	return &MsgTransferNFT{
		TokenID:   tokenID,
		Sender:    sender,
		Recipient: recipient,
	}
}

// Route should return the name of the module
func (m MsgTransferNFT) Route() string { return RouterKey }

// Type should return the action
func (m MsgTransferNFT) Type() string { return "transfer_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgTransferNFT) ValidateBasic() sdk.Error {
	if m.Sender.Empty() {
		return sdk.ErrInvalidAddress(m.Sender.String())
	}
	if m.Recipient.Empty() {
		return sdk.ErrInvalidAddress(m.Recipient.String())
	}
	if len(m.TokenID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgTransferNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Sender}
}

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
func (m MsgPutNFTOnMarket) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if len(m.TokenID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
	}
	if len(m.TokenID) > MaxTokenIDLength {
		return sdk.ErrUnknownRequest("TokenID has invalid format")
	}
	if m.Price.IsZero() || m.Price.IsAnyNegative() {
		return sdk.ErrUnknownRequest("Price cannot be zero or negative")
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
func (m MsgBuyNFT) ValidateBasic() sdk.Error {
	if m.Buyer.Empty() {
		return sdk.ErrInvalidAddress(m.Buyer.String())
	}
	if len(m.TokenID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
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

func (m MsgCreateFungibleToken) ValidateBasic() sdk.Error {
	if m.Creator.Empty() {
		return sdk.ErrInvalidAddress(m.Creator.String())
	}
	if len(m.Denom) < MinDenomLength || len(m.Denom) > MaxDenomLength {
		return sdk.ErrUnknownRequest("denom is not valid")
	}
	if m.Amount <= 0 {
		return sdk.ErrUnknownRequest("amount is invalid")
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

func (m MsgTransferFungibleTokens) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if m.Recipient.Empty() {
		return sdk.ErrInvalidAddress(m.Recipient.String())
	}
	if len(m.Denom) < MinDenomLength || len(m.Denom) > MaxDenomLength {
		return sdk.ErrUnknownRequest("denom is not valid")
	}
	if m.Amount <= 0 {
		return sdk.ErrUnknownRequest("amount is invalid")
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

type MsgUpdateNFTParams struct {
	Owner   sdk.AccAddress `json:"owner"`
	Params  []NFTParam     `json:"params"`
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
func (m MsgUpdateNFTParams) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if m.TokenID == "" {
		return sdk.ErrUnknownRequest(m.TokenID)
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

func (m MsgBurnFungibleTokens) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if len(m.Denom) < 3 || len(m.Denom) > 16 {
		return sdk.ErrUnknownRequest("denom is not valid")
	}
	if m.Amount <= 0 {
		return sdk.ErrUnknownRequest("amount is invalid")
	}
	return nil
}

func (m MsgBurnFungibleTokens) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

func (m MsgBurnFungibleTokens) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.Owner}
}

func isTokenURIValid(tokenURI string) bool {
	return true
}
