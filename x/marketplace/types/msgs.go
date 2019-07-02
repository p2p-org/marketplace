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
	Price       sdk.Coin       `json:"price"`
}

func NewMsgMintNFT(
	owner sdk.AccAddress,
	name,
	description,
	image,
	tokenURI string,
	price sdk.Coin,
) *MsgMintNFT {
	return &MsgMintNFT{
		Owner:       owner,
		Name:        name,
		Description: description,
		Image:       image,
		TokenURI:    tokenURI,
		Price:       price,
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
	Sender    sdk.AccAddress `json:"sender"`
	Recipient sdk.AccAddress `json:"recipient"`
	TokenID   string         `json:"token_id"`
}

func NewMsgTransferNFT(sender, recipient sdk.AccAddress, tokenID string) *MsgTransferNFT {
	return &MsgTransferNFT{
		Sender:    sender,
		Recipient: recipient,
		TokenID:   tokenID,
	}
}

// Route should return the name of the module
func (m MsgTransferNFT) Route() string { return RouterKey }

// Type should return the action
func (m MsgTransferNFT) Type() string { return "transfer_nfr" }

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
