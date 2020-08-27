package poa_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	"github.com/ltacker/poa"
	"github.com/ltacker/poa/types"
)

func TestEndBlocker(t *testing.T) {
	ctx, poaKeeper := poa.MockContext()
	validator1, _ := poa.MockValidator()
	validator2, _ := poa.MockValidator()
	validator3, _ := poa.MockValidator()
	validator4, _ := poa.MockValidator()
	validator5, _ := poa.MockValidator()

	// Set the validator in the keepr
	poaKeeper.AppendValidator(ctx, validator1)
	poaKeeper.AppendValidator(ctx, validator2)
	poaKeeper.AppendValidator(ctx, validator3)
	poaKeeper.AppendValidator(ctx, validator4)
	poaKeeper.AppendValidator(ctx, validator5)

	// Simulate validator 2 as if it is already in the validator set
	poaKeeper.SetValidatorState(ctx, validator2, types.ValidatorStateJoined)

	// Simulate validator 4 and 5 as if those are leaving the validator set
	poaKeeper.SetValidatorState(ctx, validator4, types.ValidatorStateLeaving)
	poaKeeper.SetValidatorState(ctx, validator5, types.ValidatorStateLeaving)

	updates := poa.EndBlocker(ctx, poaKeeper)

	// There should be 4 updates
	if len(updates) != 4 {
		t.Errorf("EndBlocker should perform 4 updates, found %v updates", len(updates))
	}

	// Check the updates
	val1Update := validator1.ABCIValidatorUpdateAppend()
	val3Update := validator3.ABCIValidatorUpdateAppend()
	val4Update := validator4.ABCIValidatorUpdateRemove()
	val5Update := validator5.ABCIValidatorUpdateRemove()
	for _, update := range updates {

		// Check if the update has the correct power
		switch {
		case cmp.Equal(update.GetPubKey(), val1Update.GetPubKey()):
			if update.GetPower() != 1 {
				t.Errorf("Validator 1 should join")
			}
		case cmp.Equal(update.GetPubKey(), val3Update.GetPubKey()):
			if update.GetPower() != 1 {
				t.Errorf("Validator 3 should join")
			}
		case cmp.Equal(update.GetPubKey(), val4Update.GetPubKey()):
			if update.GetPower() != 0 {
				t.Errorf("Validator 4 should join")
			}
		case cmp.Equal(update.GetPubKey(), val5Update.GetPubKey()):
			if update.GetPower() != 0 {
				t.Errorf("Validator 5 should join")
			}
		default:
			t.Errorf("EndBlocker returns a unknown update: %v", update)
		}
	}
}
