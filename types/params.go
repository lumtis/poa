package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/x/params"
)

// Default parameter namespace
const (
	// Default max number of validators
	DefaultMaxValidators uint32 = 15
)

// Parameter store keys
var (
	KeyMaxValidators = []byte("MaxValidators")
)

// ParamKeyTable for poa module
func ParamKeyTable() params.KeyTable {
	return params.NewKeyTable().RegisterParamSet(&Params{})
}

// Params - used for initializing default parameter for poa at genesis
type Params struct {
	MaxValidators uint32
}

// NewParams creates a new Params object
func NewParams(maxValidators uint32) Params {
	return Params{
		MaxValidators: maxValidators,
	}
}

// String implements the stringer interface for Params
func (p Params) String() string {
	return fmt.Sprintf("%d", p.MaxValidators)
}

// ParamSetPairs - Implements params.ParamSet
func (p *Params) ParamSetPairs() params.ParamSetPairs {
	return params.ParamSetPairs{
		params.NewParamSetPair(KeyMaxValidators, &p.MaxValidators, func(value interface{}) error { return nil }),
	}
}

// DefaultParams defines the parameters for this module
func DefaultParams() Params {
	return NewParams(DefaultMaxValidators)
}

// Validate a set of params
func (p Params) Validate() error {
	if err := validateMaxValidators(p.MaxValidators); err != nil {
		return err
	}
	return nil
}

// Validate maxValidators param
func validateMaxValidators(i interface{}) error {
	v, ok := i.(uint32)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T", i)
	}

	if v == 0 {
		return fmt.Errorf("max validators must be positive: %d", v)
	}

	return nil
}
