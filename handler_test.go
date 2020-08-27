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

	// The application is submitted correctly
	msg := types.NewMsgSubmitApplication(validator)
	_, err := handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgSubmitApplication should submit an application, got error %v", err)
	}
	_, found := poaKeeper.GetApplication(ctx, validator.GetOperator())
	if !found {
		t.Errorf("MsgSubmitApplication should submit an application, the application has not been found")
	}
	_, found = poaKeeper.GetApplicationByConsAddr(ctx, validator.GetConsAddr())
	if !found {
		t.Errorf("MsgSubmitApplication should submit an application, the application has not been found by cons addr")
	}

	// A new application with the same validator cannot be created
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrAlreadyApplying.Error() {
		t.Errorf("MsgSubmitApplication with duplicate, error should be %v, got %v", types.ErrAlreadyApplying.Error(), err.Error())
	}

	// Test with quorum=0
	ctx, poaKeeper = poa.MockContext()
	handler = poa.NewHandler(poaKeeper)
	validator, _ = poa.MockValidator()
	poaKeeper.SetParams(ctx, types.NewParams(15, 0))

	// The validator should be directly appended if the quorum is 0
	msg = types.NewMsgSubmitApplication(validator)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgSubmitApplication with quorum 0 should append validator, got error %v", err)
	}
	_, found = poaKeeper.GetValidator(ctx, validator.GetOperator())
	if !found {
		t.Errorf("MsgSubmitApplication with quorum 0 should append validator, the validator has not been found")
	}
	_, found = poaKeeper.GetValidatorByConsAddr(ctx, validator.GetConsAddr())
	if !found {
		t.Errorf("MsgSubmitApplication with quorum 0 should append validator, the validator has not been found by cons addr")
	}
	foundState, found := poaKeeper.GetValidatorState(ctx, validator.GetOperator())
	if !found {
		t.Errorf("MsgSubmitApplication with quorum 0 should append validator, the validator state has not been found")
	}
	if foundState != types.ValidatorStateJoining {
		t.Errorf("MsgSubmitApplication with quorum 0, the validator should have the state joining, if it is appended")
	}

	// A new application cannot be created if the validator already exist
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrAlreadyValidator.Error() {
		t.Errorf("MsgSubmitApplication with duplicate, error should be %v, got %v", types.ErrAlreadyValidator.Error(), err.Error())
	}

	// Test max validators condition
	poaKeeper.SetParams(ctx, types.NewParams(1, 0))
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrMaxValidatorsReached.Error() {
		t.Errorf("MsgSubmitApplication with max validators reached, error should be %v, got %v", types.ErrMaxValidatorsReached.Error(), err.Error())
	}
}

func TestHandleMsgProposeKick(t *testing.T) {
	// Test with maxValidator=15, quorum=66
	ctx, poaKeeper := poa.MockContext()
	handler := poa.NewHandler(poaKeeper)
	validator1, _ := poa.MockValidator()
	validator2, _ := poa.MockValidator()
	nothing, _ := poa.MockValidator()
	poaKeeper.SetParams(ctx, types.DefaultParams())

	// Add validators to validator set
	poaKeeper.AppendValidator(ctx, validator1)
	poaKeeper.AppendValidator(ctx, validator2)

	// The kick proposal is created correctly
	msg := types.NewMsgProposeKick(validator1.GetOperator(), validator2.GetOperator())
	_, err := handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgProposeKick should create a kick proposal, got error %v", err)
	}
	_, found := poaKeeper.GetKickProposal(ctx, validator1.GetOperator())
	if !found {
		t.Errorf("MsgProposeKick should create a kick proposal, the kick proposal has not been found")
	}

	// A new application with the same validator cannot be created
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrAlreadyInKickProposal.Error() {
		t.Errorf("MsgProposeKick with duplicate, error should be %v, got %v", types.ErrAlreadyInKickProposal.Error(), err.Error())
	}

	// A non validator cannot create a kick proposal
	msg = types.NewMsgProposeKick(validator2.GetOperator(), nothing.GetOperator())
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrProposerNotValidator.Error() {
		t.Errorf("MsgProposeKick sent by a non validator, error should be %v, got %v", types.ErrProposerNotValidator.Error(), err.Error())
	}

	// A non validator cannot be proposed to be kicked
	msg = types.NewMsgProposeKick(nothing.GetOperator(), validator2.GetOperator())
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrNotValidator.Error() {
		t.Errorf("MsgProposeKick propose a non validator, error should be %v, got %v", types.ErrNotValidator.Error(), err.Error())
	}

	// Test with quorum=0
	ctx, poaKeeper = poa.MockContext()
	handler = poa.NewHandler(poaKeeper)
	validator1, _ = poa.MockValidator()
	validator2, _ = poa.MockValidator()
	poaKeeper.SetParams(ctx, types.NewParams(15, 0))

	// Add validators to validator set
	poaKeeper.AppendValidator(ctx, validator1)
	poaKeeper.AppendValidator(ctx, validator2)

	// The validator should be directly appended if the quorum is 0
	msg = types.NewMsgProposeKick(validator1.GetOperator(), validator2.GetOperator())
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgProposeKick with quorum 0 should kick validator, got error %v", err)
	}
	// Check state is leaving
	foundState, found := poaKeeper.GetValidatorState(ctx, validator1.GetOperator())
	if !found {
		t.Errorf("MsgProposeKick with quorum 0 should not directly remove the validator from the validator set")
	}
	if foundState != types.ValidatorStateLeaving {
		t.Errorf("MsgProposeKick with quorum 0, the validator state should be leaving")
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
		t.Errorf("MsgVoteApplication should fail with %v, got %v", types.ErrNoApplicationFound, err)
	}

	// Cannot vote if the voter is not in validator set
	msg = types.NewMsgVote(types.VoteTypeApplication, nothing.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrVoterNotValidator.Error() {
		t.Errorf("MsgVoteApplication should fail with %v, got %v", types.ErrVoterNotValidator, err)
	}

	// Can vote an application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgVoteApplication should vote on an application, got error %v", err)
	}
	application, found := poaKeeper.GetApplication(ctx, candidate1.GetOperator())
	if !found {
		t.Errorf("MsgVoteApplication with 1/2 approve should not remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate1.GetOperator())
	if found {
		t.Errorf("MsgVoteApplication with 1/2 approve should not append the candidate to the validator set")
	}
	if application.GetTotal() != 1 {
		t.Errorf("MsgVoteApplication with approve should add one vote to the application")
	}
	if application.GetApprovals() != 1 {
		t.Errorf("MsgVoteApplication with approve should add one approve to the application")
	}

	// Second approve should append the candidate to the validator pool
	msg = types.NewMsgVote(types.VoteTypeApplication, voter2.GetOperator(), candidate1.GetOperator(), true)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgVoteApplication 2 should vote on an application, got error %v", err)
	}
	_, found = poaKeeper.GetApplication(ctx, candidate1.GetOperator())
	if found {
		t.Errorf("MsgVoteApplication with 2/2 approve should remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate1.GetOperator())
	if !found {
		t.Errorf("MsgVoteApplication with 2/2 approve should append the candidate to the validator set")
	}

	// Quorum 100%: one reject is sufficient to reject the validator application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgVoteApplication 3 should vote on an application, got error %v", err)
	}
	_, found = poaKeeper.GetApplication(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("MsgVoteApplication with 1 reject should reject the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("MsgVoteApplication application rejected should not append the candidate to the validator set")
	}

	// Reapply and set quorum to 1%
	poaKeeper.AppendApplication(ctx, candidate2)
	poaKeeper.SetParams(ctx, types.NewParams(15, 1))

	// One reject should update the vote but not reject totally the application
	msg = types.NewMsgVote(types.VoteTypeApplication, voter1.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err != nil {
		t.Errorf("MsgVoteApplication 4 should vote on an application, got error %v", err)
	}
	application, found = poaKeeper.GetApplication(ctx, candidate2.GetOperator())
	if !found {
		t.Errorf("MsgVoteApplication with 1/3 reject should not remove the application")
	}
	_, found = poaKeeper.GetValidator(ctx, candidate2.GetOperator())
	if found {
		t.Errorf("MsgVoteApplication with 1/3 reject should not append the candidate to the validator set")
	}
	if application.GetTotal() != 1 {
		t.Errorf("MsgVoteApplication with reject should add one vote to the application")
	}
	if application.GetApprovals() != 0 {
		t.Errorf("MsgVoteApplication with reject should not add one approve to the application")
	}

	// Cannot vote if validator set is full
	poaKeeper.SetParams(ctx, types.NewParams(3, 1))

	msg = types.NewMsgVote(types.VoteTypeApplication, voter2.GetOperator(), candidate2.GetOperator(), false)
	_, err = handler(ctx, msg)
	if err.Error() != types.ErrMaxValidatorsReached.Error() {
		t.Errorf("MsgVoteApplication should fail with %v, got %v", types.ErrMaxValidatorsReached, err)
	}
}
