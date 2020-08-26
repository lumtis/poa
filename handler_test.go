package poa_test

import (
	"testing"

	"github.com/ltacker/poa"
	"github.com/ltacker/poa/types"
)

func TestHandleMsgSubmitApplication(t *testing.T) {
	// Test with maxValidator=15, quorum=66
	ctx, poaKeeper := poa.MockContext()
	handler := poa.NewHandler(poaKeeper)
	validator, _ := poa.MockValidator()
	poaKeeper.SetParams(ctx, types.DefaultParams())

	msg := types.NewMsgSubmitApplication(validator)

	// The application is submitted correctly
	_, err := handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgSubmitApplication should submit an application, got error %v", err)
	}
	_, found := poaKeeper.GetApplication(ctx, validator.GetOperator())
	if !found {
		t.Errorf("NewMsgSubmitApplication should submit an application, the application has not been found")
	}
	_, found = poaKeeper.GetApplicationByConsAddr(ctx, validator.GetConsAddr())
	if !found {
		t.Errorf("NewMsgSubmitApplication should submit an application, the application has not been found by cons addr")
	}

	// A new application with the same validator cannot be created
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrAlreadyApplying.Error() {
		t.Errorf("NewMsgSubmitApplication with duplicate, error should be %v, got %v", types.ErrAlreadyApplying.Error(), err.Error())
	}

	// Test with quorum=0
	ctx, poaKeeper = poa.MockContext()
	handler = poa.NewHandler(poaKeeper)
	validator, _ = poa.MockValidator()
	poaKeeper.SetParams(ctx, types.NewParams(15, 0))

	msg = types.NewMsgSubmitApplication(validator)

	// The validator should be directly appended if the quorum is 0
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator, got error %v", err)
	}
	_, found = poaKeeper.GetValidator(ctx, validator.GetOperator())
	if !found {
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator,, the validator has not been found")
	}
	_, found = poaKeeper.GetValidatorByConsAddr(ctx, validator.GetConsAddr())
	if !found {
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator,, the validator has not been found by cons addr")
	}

	// A new applicationcannot be created if the validator already exist
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrAlreadyValidator.Error() {
		t.Errorf("NewMsgSubmitApplication with duplicate, error should be %v, got %v", types.ErrAlreadyValidator.Error(), err.Error())
	}

	// Test max validators condition
	poaKeeper.SetParams(ctx, types.NewParams(1, 0))
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrMaxValidatorsReached.Error() {
		t.Errorf("NewMsgSubmitApplication with max validators reached, error should be %v, got %v", types.ErrMaxValidatorsReached.Error(), err.Error())
	}
}
