package poa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ltacker/poa/keeper"
	"github.com/ltacker/poa/types"
	abci "github.com/tendermint/tendermint/abci/types"
)

// InitGenesis initialize default parameters
// and the keeper's address to pubkey map
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) (res []abci.ValidatorUpdate) {
	keeper.SetParams(ctx, data.Params)

	// Set validators in the storage
	for _, validator := range data.Validators {
		keeper.SetValidator(ctx, validator)
		keeper.SetValidatorByConsAddr(ctx, validator)
		res = append(res, validator.ABCIValidatorUpdateAppend())
	}

	return res
}

// ExportGenesis writes the current store values
// to a genesis file, which can be imported again
// with InitGenesis
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) (data types.GenesisState) {
	return types.GenesisState{
		Params:     keeper.GetParams(ctx),
		Validators: keeper.GetAllValidators(ctx),
	}
}
