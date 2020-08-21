package types

import (
	_ sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

var (
	_ = sdkerrors.Register(ModuleName, 1, "custom error message")
)
