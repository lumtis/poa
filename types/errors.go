package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoValidatorFound     = sdkerrors.Register(ModuleName, 1, "no validator found")
	ErrInvalidValidator     = sdkerrors.Register(ModuleName, 2, "invalid validator")
	ErrInvalidQuorumValue   = sdkerrors.Register(ModuleName, 3, "quorum should be a percentage")
	ErrInvalidVoterPoolSize = sdkerrors.Register(ModuleName, 4, "the voter pool size is incorrect")
	ErrAlreadyApplying      = sdkerrors.Register(ModuleName, 5, "the candidate is already applying to become a validator")
	ErrAlreadyValidator     = sdkerrors.Register(ModuleName, 6, "the candidate is already a validator")
	ErrMaxValidatorsReached = sdkerrors.Register(ModuleName, 7, "the maximum number of validators has been reached")
	ErrInvalidVoteMsg       = sdkerrors.Register(ModuleName, 8, "the vote message is invalid")
	ErrVoterNotValidator    = sdkerrors.Register(ModuleName, 9, "the voter is not a validator")
	ErrNoApplicationFound   = sdkerrors.Register(ModuleName, 10, "no application found")
	ErrAlreadyVoted         = sdkerrors.Register(ModuleName, 11, "the validator already voted")
)
