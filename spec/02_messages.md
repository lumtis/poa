<!--
order: 2
-->

# Messages

In this section, we describe the processing of the messages and the corresponding updates to the state.

## MsgSubmitApplication

A new application to become a validator is made using the `MsgSubmitApplication` message.

```go
type MsgSubmitApplication struct {
    Candidate Validator
}
```

This message is expected to fail if:

- another application with the operator address is already registered
- another application with the pubkey is already registered
- another validator with the operator address is already registered
- another validator with the pubkey is already registered
- the validator set is full
- the description fields are too large

This message creates and stores a new `Vote` object in the application pool.

## MsgApproveApplication

An application of a new validator can be approved using the
`MsgApproveApplication` message.  

```go
type MsgApproveApplication struct {
    VoterAddr       sdk.ValAddress
    CandidateAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the voter has already voted
- the candidate address is not in the application pool

This message updates the vote status of the application. If the approval quorum is reached, the candidate is appended into the validator set.

## MsgRejectApplication

An application of a new validator can be rejected using the
`MsgRejectApplication` message.  

```go
type MsgRejectApplication struct {
    VoterAddr       sdk.ValAddress
    CandidateAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the voter has already voted
- the candidate address is not in the application pool

This message updates the vote status of the application. If the rejection quorum is reached, the candidate is rejected from joining the validator set.

## MsgLeaveValidatorSet

A current validator arbitrarily leaves the validator set using the MsgLeaveValidatorSet message.

```go
type MsgLeaveValidatorSet struct {
    ValidatorAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the validator address is not in the validator set

The message removes the validator from the validator set.

## MsgProposeKick

A new kick proposal to kick a validator is made using the `MsgProposeKick` message.

```go
type MsgProposeKick struct {
    ValidatorAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the validator address is not in the validator set
- the validator address is already in the kick proposal pool

This message creates and stores a new `Vote` object in the kick proposal pool.

## MsgApproveKickProposal

A kick proposal can be approved using the
`MsgApproveKickProposal` message.  

```go
type MsgApproveKickProposal struct {
    VoterAddr       sdk.ValAddress
    ValidatorAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the voter has already voted
- the validator address is not in the kick proposal pool

This message updates the vote status of the kick proposal. If the approval quorum is reached, the validator is removed from the validator set.

## MsgRejectKickProposal

A kick proposal can be rejected using the
`MsgRejectKickProposal` message.  

```go
type MsgRejectApplication struct {
    VoterAddr       sdk.ValAddress
    ValidatorAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the voter has already voted
- the candidate address is not in the kick proposal pool

This message updates the vote status of the application. If the rejection quorum is reached, the validator is not removed from the validator set.