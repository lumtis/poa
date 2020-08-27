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
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator, the validator has not been found")
	}
	_, found = poaKeeper.GetValidatorByConsAddr(ctx, validator.GetConsAddr())
	if !found {
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator, the validator has not been found by cons addr")
	}
	foundState, found := poaKeeper.GetValidatorState(ctx, validator.GetOperator())
	if !found {
		t.Errorf("NewMsgSubmitApplication with quorum 0 should append validator, the validator state has not been found")
	}
	if foundState != types.ValidatorStateJoining {
		t.Errorf("NewMsgSubmitApplication with quorum 0, the validator should have the state joining, if it is appended")
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

func TestHandleMsgVoteApplication(t *testing.T) {
	ctx, poaKeeper := poa.MockContext()
	handler := poa.NewHandler(poaKeeper)
	voter1, _ := poa.MockValidator()
	voter2, _ := poa.MockValidator()
	candidate1, _ := poa.MockValidator()
	candidate2, _ := poa.MockValidator()
	nothing, _ := poa.MockValidator()
	poaKeeper.SetParams(ctx, types.NewParams(15, 100)) // Set quorum to 100%

	// Add voter to validator set
	poaKeeper.AppendValidator(ctx, voter1)
	poaKeeper.AppendValidator(ctx, voter2)

	// Add candidate to application pool
	poaKeeper.AppendApplication(ctx, candidate1)
	poaKeeper.AppendApplication(ctx, candidate2)

	// Cannot vote if candidate is not in application pool
	msg := types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), nothing.GetOperator(), true)
	_, err := handler(ctx, msg)
	if err.Error() != types.ErrNoApplicationFound.Error() {
		t.Errorf("NewMsgVote should fail with %v, got %v", types.ErrNoApplicationFound, err)
	}

	// Cannot vote if the voter is not in validator set
	msg = types.NewMsgVote(types.VoteTypeApplication, nothing.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrVoterNotValidator.Error() {
		t.Errorf("NewMsgVote should fail with %v, got %v", types.ErrVoterNotValidator, err)
	}

	// Can vote an application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgVote should vote on an application, got error %v", err)
	}
	application, found := poaKeeper.GetApplication(ctx, candidate1.GetOperator())
	if !found {
		t.Errorf("NewMsgVote with 1/2 approve should not remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate1.GetOperator())
	if found {
		t.Errorf("NewMsgVote with 1/2 approve should not append the candidate to the validator set")
	}
	if application.GetTotal() != 1 {
		t.Errorf("NewMsgVote with approve should add one vote to the application")
	}
	if application.GetApprovals() != 1 {
		t.Errorf("NewMsgVote with approve should add one approve to the application")
	}

	// Second approve should append the candidate to the validator pool
	msg = types.NewMsgVote(types.VoteTypeApplication, voter2.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgVote 2 should vote on an application, got error %v", err)
	}
	_, found = poaKeeper.GetApplication(ctx, candidate1.GetOperator())
	if found {
		t.Errorf("NewMsgVote with 2/2 approve should remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate1.GetOperator())
	if !found {
		t.Errorf("NewMsgVote with 2/2 approve should append the candidate to the validator set")
	}

	// Quorum 100%: one reject is sufficient to reject the validator application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgVote 3 should vote on an application, got error %v", err)
	}
	_, found = poaKeeper.GetApplication(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("NewMsgVote with 1 reject should reject the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("NewMsgVote application rejected should not append the candidate to the validator set")
	}

	// Reapply and set quorum to 1%
	poaKeeper.AppendApplication(ctx, candidate2)
	poaKeeper.SetParams(ctx, types.NewParams(15, 1))

	// One reject should update the vote but not reject totally the application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("NewMsgVote 4 should vote on an application, got error %v", err)
	}
	application, found = poaKeeper.GetApplication(ctx, candidate2.GetOperator())
	if !found {
		t.Errorf("NewMsgVote with 1/3 reject should not remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("NewMsgVote with 1/3 reject should not append the candidate to the validator set")
	}
	if application.GetTotal() != 1 {
		t.Errorf("NewMsgVote with reject should add one vote to the application")
	}
	if application.GetApprovals() != 0 {
		t.Errorf("NewMsgVote with reject should not add one approve to the application")
	}

	// Cannot vote if validator set is full
	poaKeeper.SetParams(ctx, types.NewParams(3, 1))

	msg = types.NewMsgVote(types.VoteTypeApplication, voter2.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrMaxValidatorsReached.Error() {
		t.Errorf("NewMsgVote should fail with %v, got %v", types.ErrMaxValidatorsReached, err)
	}
}
