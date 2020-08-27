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
		case types.MsgVote:
			return handleMsgVote(ctx, k, msg)
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
		k.AppendValidator(ctx, msg.Candidate)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeAppendValidator,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
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
		k.AppendApplication(ctx, msg.Candidate)

		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeSubmitApplication,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
			),
		)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}

// handleMsgVote handles a vote performed by a validator
func handleMsgVote(ctx sdk.Context, k keeper.Keeper, msg types.MsgVote) (*sdk.Result, error) {
	switch msg.Type {
	case types.VoteTypeApplication:
		return handleMsgVoteApplication(ctx, k, msg)
	default:
		return nil, types.ErrInvalidVoteMsg
	}
}

func handleMsgVoteApplication(ctx sdk.Context, k keeper.Keeper, msg types.MsgVote) (*sdk.Result, error) {
	// Check max validator is not reached. If max validator is reached, not application can be voted
	allValidators := k.GetAllValidators(ctx)
	validatorCount := len(allValidators)
	maxValidator := k.MaxValidators(ctx)
	if uint16(validatorCount) == maxValidator {
		return nil, types.ErrMaxValidatorsReached
	}

	// The voter must be a validator
	_, found := k.GetValidator(ctx, msg.VoterAddr)
	if !found {
		return nil, types.ErrVoterNotValidator
	}

	// Check the application exist
	application, found := k.GetApplication(ctx, msg.CandidateAddr)
	if !found {
		return nil, types.ErrNoApplicationFound
	}

	// Check if already voted and vote
	alreadyVoted := application.AddVote(msg.VoterAddr, msg.Approve)
	if alreadyVoted {
		return nil, types.ErrAlreadyVoted
	}

	// Emit the vote event
	if msg.Approve {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeApproveApplication,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyVoter, msg.VoterAddr),
				sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
			),
		)
	} else {
		ctx.EventManager().EmitEvent(
			sdk.NewEvent(
				types.EventTypeRejectApplication,
				sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
				sdk.NewAttribute(types.AttributeKeyVoter, msg.VoterAddr),
				sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
			),
		)
	}

	// Check if the quorum has been reached
	reached, approved, err := application.CheckQuorum(validatorCount, k.Quorum())
	if err != nil {
		return nil, err
	}

	if reached {
		if approved {
			// Candidate is appended to the validator set
			k.RemoveApplication(ctx, msg.CandidateAddr)
			k.AppendValidator(ctx, application.GetSubject())

			// Emit approved event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeAppendValidator,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
				),
			)
		} else {
			// Candidate is rejected from joining the validator set
			k.RemoveApplication(ctx, msg.CandidateAddr)

			// Emit rejected event
			ctx.EventManager().EmitEvent(
				sdk.NewEvent(
					types.EventTypeRejectValidator,
					sdk.NewAttribute(sdk.AttributeKeyModule, types.AttributeValueCategory),
					sdk.NewAttribute(types.AttributeKeyCandidate, msg.Candidate.GetOperator().String()),
				),
			)
		}
	} else {
		// Quorum has not been reached yet, update the vote
		k.SetApplication(ctx, application)
	}

	return &sdk.Result{Events: ctx.EventManager().Events()}, nil
}
