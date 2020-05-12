package nftIBC

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewHandler returns sdk.Handler for IBC token transfer module messages
func NewHandler(k Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		switch msg := msg.(type) {
		case MsgTransferNFT:
			return handleMsgTransferNFT(ctx, k, msg)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unrecognized ICS-20 transfer message type: %T", msg)
		}
	}
}

func handleMsgTransferNFT(ctx sdk.Context, k Keeper, msg MsgTransferNFT) (*sdk.Result, error) {
	if err := k.SendTransfer(
		ctx, msg.SourcePort, msg.SourceChannel, msg.DestinationHeight, msg.Id, msg.Denom, msg.Sender, msg.Receiver,
	); err != nil {
		return nil, err
	}

	k.Logger(ctx).Info("IBC transfer: %s from %s to %s", msg.Id, msg.Sender, msg.Receiver)

	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Sender.String()),
			sdk.NewAttribute(AttributeKeyReceiver, msg.Receiver),
		),
	)

	return &sdk.Result{
		Events: ctx.EventManager().Events().ToABCIEvents(),
	}, nil
}
