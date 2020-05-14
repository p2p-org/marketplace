package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	host "github.com/cosmos/cosmos-sdk/x/ibc/24-host"
)

// msg types
const (
	TypeMsgTransfer = "transferNFT"
)

// NewMsgTransfer creates a new MsgTransfer instance
func NewMsgTransferNFT(
	sourcePort, sourceChannel string, destHeight uint64, sender, receiver sdk.AccAddress, id, denom string,
) MsgTransferNFT {
	return MsgTransferNFT{
		SourcePort:        sourcePort,
		SourceChannel:     sourceChannel,
		DestinationHeight: destHeight,
		Sender:            sender,
		Receiver:          receiver,
		Id:                id,
		Denom:             denom,
	}
}

// Route implements sdk.Msg
func (MsgTransferNFT) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (MsgTransferNFT) Type() string {
	return TypeMsgTransfer
}

// ValidateBasic implements sdk.Msg
func (msg MsgTransferNFT) ValidateBasic() error {
	if err := host.DefaultPortIdentifierValidator(msg.SourcePort); err != nil {
		return sdkerrors.Wrap(err, "invalid source port ID")
	}
	if err := host.DefaultChannelIdentifierValidator(msg.SourceChannel); err != nil {
		return sdkerrors.Wrap(err, "invalid source channel ID")
	}
	if msg.Sender == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing sender address")
	}
	if msg.Receiver == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing recipient address")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (msg MsgTransferNFT) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

// GetSigners implements sdk.Msg
func (msg MsgTransferNFT) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}
