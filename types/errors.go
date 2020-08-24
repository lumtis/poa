package types

import (
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoValidatorFound = sdkerrors.Register(ModuleName, 1, "no validator found")
	ErrInvalidValidator = sdkerrors.Register(ModuleName, 2, "invalid validator")
)
