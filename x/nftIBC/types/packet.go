package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

func NewNFTPacketData(
	id string, denom string, owner sdk.AccAddress) NFTPacketData {
	return NFTPacketData{
		Id:    id,
		Denom: denom,
		Owner: owner,
	}
}

// ValidateBasic is used for validating the token transfer
func (ftpd *NFTPacketData) ValidateBasic() error {
	if ftpd.Owner == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing owner address")
	}
	if ftpd.Receiver == nil {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidAddress, "missing receiver address")
	}
	return nil
}

// GetBytes is a helper for serialising
func (ftpd NFTPacketData) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ftpd))
}

// GetBytes is a helper for serialising
func (ack *NFTPacketAcknowledgement) GetBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(ack))
}
