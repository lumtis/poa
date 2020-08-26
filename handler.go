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
	// Check max validator is not reached
	allValidators := k.GetAllValidators(ctx)
	maxValidator := k.MaxValidators(ctx)
	if uint16(len(allValidators)) == maxValidator {
		return nil, types.ErrMaxValidatorsReached
	}
	// Candidate should not be a validator
	_, found := k.GetValidator(ctx, msg.Candidate.GetOperator())
	if found {
		return nil, types.ErrAlreadyValidator
	}
	_, found = k.GetValidatorByConsAddr(ctx, msg.Candidate.GetConsAddr())
	if found {
		return nil, types.ErrAlreadyValidator
	}

	// If quorum is 0 the application is immediately approved
	if k.Quorum(ctx) == 0 {
		// The validator is directly appended in the validator set
		k.SetValidator(ctx, msg.Candidate)
		k.SetValidatorByConsAddr(ctx, msg.Candidate)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAppendValidator,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyValidator, msg.Candidate.GetOperator().String()),
			),
		)
	} else {
		// If quorum is more than 0, we create a application vote

		// Candidate should not be already applying
		_, found = k.GetApplication(ctx, msg.Candidate.GetOperator())
		if found {
			return nil, types.ErrAlreadyApplying
		}
		_, found = k.GetApplicationByConsAddr(ctx, msg.Candidate.GetConsAddr())
		if found {
			return nil, types.ErrAlreadyApplying
		}

		// Create the new application
		applicationEmptyVote := types.NewVote(msg.Candidate)
		k.SetApplication(ctx, applicationEmptyVote)
		k.SetApplicationByConsAddr(ctx, applicationEmptyVote)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSubmitApplication,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyValidator, msg.Candidate.GetOperator().String()),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
