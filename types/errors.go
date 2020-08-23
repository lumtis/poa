package types

import (
	_ sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	ErrNoValidatorFound = sdkerrors.Register(ModuleName, 1, "no validator found")
)
