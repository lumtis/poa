<!--
order: 2
-->

# Messages

In this section we describe the processing of the messages and the corresponding updates to the state.

## MsgAppendValidator

A validator is appended to the validator set using the `MsgAppendValidator` message.

```go
type MsgAppendValidator struct {
    Description    Description
    ValidatorAddr  sdk.ValAddress
    PubKey         crypto.PubKey
}
```

This message is expected to fail if:

- another validator with this operator address is already registered
- another validator with this pubkey is already registered
- the validator set is full
- the description fields are too large

This message creates and stores the `Validator` object at appropriate indexes.

## MsgEditValidatorDescription

The `Description` of a validator can be updated using the
`MsgEditValidator`.  

```go
type MsgEditValidatorDescription struct {
    Description     Description
    ValidatorAddr   sdk.ValAddress
}
```

This message is expected to fail if:

- the description fields are too large

This message stores the description of the updated `Validator` object.

## MsgRemoveValidator

A validator is removed from the validator set using the `MsgRemoveValidator` message.

```go
type MsgAppendValidator struct {
    ValidatorAddr  sdk.ValAddress
}
```

This message is expected to fail if:

- the validator is not present in the validator set