package poa_test

import (
	"fmt"
	"reflect"
	"testing"

	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/log"
	tmtypes "github.com/tendermint/tendermint/types"
	dbm "github.com/tendermint/tm-db"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/params"
	"github.com/ltacker/poa"
	"github.com/ltacker/poa/keeper"
	"github.com/ltacker/poa/types"
)

// Context and keeper used for mocking purpose
func mockContext() (sdk.Context, keeper.Keeper) {
	// Store keys
	keys := sdk.NewKVStoreKeys(types.StoreKey, params.StoreKey)
	tKeys := sdk.NewTransientStoreKeys(params.TStoreKey)

	cdc := codec.New()

	// Create the params keeper
	paramsKeeper := params.NewKeeper(cdc, keys[params.StoreKey], tKeys[params.TStoreKey])

	// Create a poa keeper
	poaKeeper := keeper.NewKeeper(cdc, keys[types.StoreKey], paramsKeeper.Subspace(types.ModuleName))

	// Create multiStore in memory
	db := dbm.NewMemDB()
	cms := store.NewCommitMultiStore(db)

	// Mount stores
	cms.MountStoreWithDB(keys[types.StoreKey], sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(keys[params.StoreKey], sdk.StoreTypeIAVL, db)
	cms.MountStoreWithDB(tKeys[params.TStoreKey], sdk.StoreTypeTransient, db)
	cms.LoadLatestVersion()

	// Create context
	ctx := sdk.NewContext(cms, abci.Header{}, false, log.NewNopLogger())

	return ctx, poaKeeper
}

// Create a validator for test
func mockValidator() (types.Validator, string) {
	operatorAddress, err := sdk.ValAddressFromBech32("cosmosvaloper1jdq0qugyrtcp4h0mhv3mx3we7gql3mc68na7f3")
	if err != nil {
		panic(fmt.Sprintf("Cannot get the address: %v", err))
	}
	validatorDescription := types.Description{
		Moniker:         "Moniker",
		Identity:        "Identity",
		Website:         "Website",
		SecurityContact: "SecurityContact",
		Details:         "Details",
	}
	consPubKey := "cosmosvalconspub1zcjduepq78q2f33h5cz32g65jsnflftllg95lrpwhanpr9dnph7mnnlqvh0sxapk0c"
	validator := types.Validator{
		OperatorAddress: operatorAddress,
		ConsensusPubkey: consPubKey,
		Description:     validatorDescription,
	}

	return validator, consPubKey
}

func TestValidateGenesis(t *testing.T) {
	validator, _ := mockValidator()

	// Valid genesis
	validGenesis := types.NewGenesisState(types.DefaultParams(), []types.Validator{validator})
	if types.ValidateGenesis(validGenesis) != nil {
		t.Errorf("The genesis state %v should be valid", validGenesis)
	}

	// A genesis with two validators with the same consensus pukey is invalid
	invalidGenesis := types.NewGenesisState(types.DefaultParams(), []types.Validator{validator, validator})
	if types.ValidateGenesis(invalidGenesis) == nil {
		t.Errorf("The genesis state %v should not be valid", invalidGenesis)
	}

	// Default genesis state
	if types.ValidateGenesis(types.DefaultGenesisState()) != nil {
		t.Errorf("The default genesis state should be valid")
	}
}

func TestInitGenesis(t *testing.T) {
	ctx, poaKeeper := mockContext()
	validator, consPubKey := mockValidator()

	// Test genesis data
	testGenesis := types.NewGenesisState(types.DefaultParams(), []types.Validator{validator})

	// InitGenesis
	validatorUpdates := poa.InitGenesis(ctx, poaKeeper, testGenesis)

	// Only one update
	if len(validatorUpdates) != 1 {
		t.Errorf("Should get exactly one validator update")
	}

	// No weight
	power := validatorUpdates[0].Power
	if power != 1 {
		t.Errorf("power should be 1, got %v", power)
	}

	// Correct public key
	pubKey, err := tmtypes.PB2TM.PubKey(validatorUpdates[0].PubKey)
	if err != nil {
		t.Errorf("Incorrect public key: %v", err)
	}
	pubKeyString := sdk.MustBech32ifyPubKey(sdk.Bech32PubKeyTypeConsPub, pubKey)
	if pubKeyString != consPubKey {
		t.Errorf("validator PubKey should be %v, got %v", consPubKey, pubKeyString)
	}
}

func TestExportGenesis(t *testing.T) {
	ctx, poaKeeper := mockContext()
	validator, _ := mockValidator()

	// Manually set values in keeper
	poaKeeper.SetValidator(ctx, validator)
	poaKeeper.SetParams(ctx, types.DefaultParams())

	exportedGenesis := poa.ExportGenesis(ctx, poaKeeper)

	if !reflect.DeepEqual(exportedGenesis.Params, types.DefaultParams()) {
		t.Errorf("Exported genesis param shoud be: %v, not %v", types.DefaultParams(), exportedGenesis.Params)
	}

	if !reflect.DeepEqual(exportedGenesis.Validators, []types.Validator{validator}) {
		t.Errorf("Exported genesis validators shoud be: %v, not %v", []types.Validator{validator}, exportedGenesis.Validators)
	}
}
