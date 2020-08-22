package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ltacker/poa/x/poa/types"
)

// Get a validator
func (k Keeper) GetValidator(ctx sdk.Context, addr sdk.ValAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)

	// Search the value
	value := store.Get(types.GetValidatorKey(addr))
	if value == nil {
		return validator, false
	}

	// Return the value
	validator = types.MustUnmarshalValidator(k.cdc, value)
	return validator, true
}

// Get a validator by consensus address
func (k Keeper) GetValidatorByConsAddr(ctx sdk.Context, consAddr sdk.ConsAddress) (validator types.Validator, found bool) {
	store := ctx.KVStore(k.storeKey)

	opAddr := store.Get(types.GetValidatorByConsAddrKey(consAddr))
	if opAddr == nil {
		return validator, false
	}

	return k.GetValidator(ctx, opAddr)
}

// Set validator details
func (k Keeper) SetValidator(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalValidator(k.cdc, validator)
	store.Set(types.GetValidatorKey(validator.OperatorAddress), bz)
}

// Set validator consensus address
func (k Keeper) SetValidatorByConsAddr(ctx sdk.Context, validator types.Validator) {
	store := ctx.KVStore(k.storeKey)
	store.Set(types.GetValidatorByConsAddrKey(validator.GetConsAddr()), validator.OperatorAddress)
}

// Get the set of all validators
func (k Keeper) GetAllValidators(ctx sdk.Context) (validators []types.Validator) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ValidatorsKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		validator := types.MustUnmarshalValidator(k.cdc, iterator.Value())
		validators = append(validators, validator)
	}

	return validators
}
