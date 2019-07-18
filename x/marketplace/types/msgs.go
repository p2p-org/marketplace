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
// MsgSellNFT
//
// --------------------------------------------------------------------------

type MsgSellNFT struct {
	Owner sdk.AccAddress `json:"owner"`
	// Beneficiary is the cosmos user who gets the commission for this transaction.
	Beneficiary sdk.AccAddress `json:"beneficiary"`
	TokenID     string         `json:"token_id"`
	Price       sdk.Coins      `json:"price"`
}

func NewMsgSellNFT(owner, beneficiary sdk.AccAddress, tokenID string, price sdk.Coins) *MsgSellNFT {
	return &MsgSellNFT{
		Owner:       owner,
		TokenID:     tokenID,
		Price:       price,
		Beneficiary: beneficiary,
	}
}

// Route should return the name of the module
func (m MsgSellNFT) Route() string { return RouterKey }

// Type should return the action
func (m MsgSellNFT) Type() string { return "sell_nft" }

// ValidateBasic runs stateless checks on the message
func (m MsgSellNFT) ValidateBasic() sdk.Error {
	if m.Owner.Empty() {
		return sdk.ErrInvalidAddress(m.Owner.String())
	}
	if len(m.TokenID) == 0 {
		return sdk.ErrUnknownRequest("TokenID cannot be empty")
	}
	return nil
}

// GetSignBytes encodes the message for signing
func (m MsgSellNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners defines whose signature is required
func (m MsgSellNFT) GetSigners() []sdk.AccAddress {
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
