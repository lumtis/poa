package poa

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/ltacker/poa/x/poa/keeper"
	abci "github.com/tendermint/tendermint/abci/types"
)

// BeginBlocker check for infraction evidence or downtime of validators
// on every begin block
func BeginBlocker(ctx sdk.Context, req abci.RequestBeginBlock, k keeper.Keeper) {
	return
}

// EndBlocker called every block, process inflation, update validator set.
func EndBlocker(ctx sdk.Context, k keeper.Keeper) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}
