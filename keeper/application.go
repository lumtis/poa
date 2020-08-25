package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ltacker/poa/types"
)

// Get an application
func (k Keeper) GetApplication(ctx sdk.Context, addr sdk.ValAddress) (application types.Vote, found bool) {
	store := ctx.KVStore(k.storeKey)

	// Search the value
	value := store.Get(types.GetApplicationKey(addr))
	if value == nil {
		return application, false
	}

	// Return the value
	application = types.MustUnmarshalVote(k.cdc, value)
	return application, true
}

// Set application details
func (k Keeper) SetApplication(ctx sdk.Context, application types.Vote) {
	store := ctx.KVStore(k.storeKey)
	bz := types.MustMarshalVote(k.cdc, application)
	store.Set(types.GetApplicationKey(application.GetSubject().GetOperator()), bz)
}

// Get the set of all application
func (k Keeper) GetAllApplications(ctx sdk.Context) (applications []types.Vote) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.ApplicationPoolKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		application := types.MustUnmarshalVote(k.cdc, iterator.Value())
		applications = append(applications, application)
	}

	return applications
}
