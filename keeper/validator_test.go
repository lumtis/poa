package keeper_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/ltacker/poa"
)

func TestGetValidator(t *testing.T) {
	ctx, poaKeeper := poa.MockContext()
	validator1, _ := poa.MockValidator()
	validator2, _ := poa.MockValidator()

	poaKeeper.SetValidator(ctx, validator1)

	// Should find the correct validator
	retrievedValidator, found := poaKeeper.GetValidator(ctx, validator1.GetOperator())
	if !found {
		t.Errorf("GetValidator should find validator if it has been set")
	}

	if !cmp.Equal(validator1, retrievedValidator) {
		t.Errorf("GetValidator should find %v, found %v", validator1, retrievedValidator)
	}

	// Should not find a unset validator
	_, found = poaKeeper.GetValidator(ctx, validator2.GetOperator())
	if found {
		t.Errorf("GetValidator should not find validator if it has not been set")
	}
}

func TestGetValidatorByConsAddr(t *testing.T) {
	ctx, poaKeeper := poa.MockContext()
	validator1, _ := poa.MockValidator()
	validator2, _ := poa.MockValidator()

	poaKeeper.SetValidator(ctx, validator1)
	poaKeeper.SetValidatorByConsAddr(ctx, validator1)

	// Should find the correct validator
	retrievedValidator, found := poaKeeper.GetValidatorByConsAddr(ctx, validator1.GetConsAddr())
	if !found {
		t.Errorf("GetValidatorByConsAddr should find validator if it has been set")
	}

	if !cmp.Equal(validator1, retrievedValidator) {
		t.Errorf("GetValidatorByConsAddr should find %v, found %v", validator1, retrievedValidator)
	}

	// Should not find a unset validator
	_, found = poaKeeper.GetValidator(ctx, validator2.GetOperator())
	if found {
		t.Errorf("GetValidatorByConsAddr should not find validator if it has not been set")
	}

	// Should not find the validator if we call SetValidatorByConsAddr without SetValidator
	poaKeeper.SetValidatorByConsAddr(ctx, validator2)
	_, found = poaKeeper.GetValidator(ctx, validator2.GetOperator())
	if found {
		t.Errorf("GetValidatorByConsAddr should not find validator if it has not been set with SetValidator")
	}
}

func TestGetAllValidators(t *testing.T) {
	ctx, poaKeeper := poa.MockContext()
	validator1, _ := poa.MockValidator()
	validator2, _ := poa.MockValidator()

	poaKeeper.SetValidator(ctx, validator1)
	poaKeeper.SetValidator(ctx, validator2)

	retrievedValidators := poaKeeper.GetAllValidators(ctx)
	if len(retrievedValidators) != 2 {
		t.Errorf("GetAllValidators should find %v validators, found %v", 2, len(retrievedValidators))
	}
}
