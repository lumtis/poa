package poa

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/ltacker/poa/keeper"
	"github.com/ltacker/poa/types"
)

// NewHandler creates an sdk.Handler for all the poa type messages
func NewHandler(k keeper.Keeper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) (*sdk.Result, error) {
		ctx = ctx.WithEventManager(sdk.NewEventManager())
		switch msg := msg.(type) {
		case types.MsgSubmitApplication:
			return handleMsgSubmitApplication(ctx, k, msg)
		default:
			errMsg := fmt.Sprintf("unrecognized %s message type: %T", types.ModuleName, msg)
			return nil, sdkerrors.Wrap(sdkerrors.ErrUnknownRequest, errMsg)
		}
	}
}

// handleMsgSubmitApplication create a new application to become a validator
func handleMsgSubmitApplication(ctx sdk.Context, k keeper.Keeper, msg types.MsgSubmitApplication) (*sdk.Result, error) {
	// validator should not be already applying
	_, found := k.GetApplication(ctx, msg.Candidate.GetOperator())
	if found {
		return nil, types.ErrAlreadyApplying
	}

	// TODO: Verify the validator doesn't already exist
	// TODO: Verify maximum number of validator is not reached

	applicationEmptyVote := types.NewVote(msg.Candidate)
	k.SetApplication(ctx, applicationEmptyVote)

	// TODO: Define your msg events
	// ctx.EventManager().EmitEvent(
	// 	sdk.NewEvent(
	// 		sdk.EventTypeMessage,
	// 		sdk.NewAttribute(sdk.AttributeKeyModule, AttributeValueCategory),
	// 		sdk.NewAttribute(sdk.AttributeKeySender, msg.ValidatorAddr.String()),
	// 	),
	// )

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
