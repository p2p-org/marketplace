package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

const RouterKey = ModuleName // this was defined in your key.go file

// --------------------------------------------------------------------------
//
// MsgMintNFT
//
// --------------------------------------------------------------------------

type MsgMintNFT struct {
	Owner       sdk.AccAddress `json:"owner"`
	Name        string         `json:"name"`
	Description string         `json:"description"`
	Image       string         `json:"image"`
	TokenURI    string         `json:"token_uri"`
}

func NewMsgMintNFT(
	owner sdk.AccAddress,
	name,
	description,
	image,
	tokenURI string,
) *MsgMintNFT {
	return &MsgMintNFT{
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
	Owner   sdk.AccAddress `json:"owner"`
	TokenID string         `json:"token_id"`
	Price   sdk.Coin       `json:"price"`
}

func NewMsgSellNFT(owner sdk.AccAddress, tokenID string, price sdk.Coin) *MsgSellNFT {
	return &MsgSellNFT{
		Owner:   owner,
		TokenID: tokenID,
		Price:   price,
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
	Buyer   sdk.AccAddress `json:"buyer"`
	TokenID string         `json:"token_id"`
}

func NewMsgBuyNFT(owner sdk.AccAddress, tokenID string) *MsgBuyNFT {
	return &MsgBuyNFT{
		Buyer:   owner,
		TokenID: tokenID,
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
