package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoValidatorFound     = sdkerrors.Register(ModuleName, 1, "no validator found")
	ErrInvalidValidator     = sdkerrors.Register(ModuleName, 2, "invalid validator")
	ErrInvalidQuorumValue   = sdkerrors.Register(ModuleName, 3, "quorum should be a percentage")
	ErrInvalidVoterPoolSize = sdkerrors.Register(ModuleName, 4, "the voter pool size is incorrect")
)
