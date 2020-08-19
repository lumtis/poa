<!--
order: 1
-->

# State

## LastTotalPower

LastTotalPower tracks the total number of validator (since 1 validator = 1 weight) recorded during the previous end block.

- LastTotalPower: `0x12 -> amino(sdk.Int)`

## Params

Params is a module-wide configuration structure that stores system parameters
and defines overall functioning of the staking module.

- Params: `Paramsspace("poa") -> amino(params)`

```go
type Params struct {
    MaxValidators   uint16          // maximum number of validators
    Master          crypto.PubKey   // public key of the master that can manage the validator set
}
```

## Validator

Validators objects should be primarily stored and accessed by the
`OperatorAddr`, an SDK validator address for the operator of the validator. `ValidatorByConsAddr` is maintained per validator object in order to fulfill
required lookups for slashing and validator-set updates.

- Validators: `0x21 | OperatorAddr -> amino(validator)`
- ValidatorsByConsAddr: `0x22 | ConsAddr -> OperatorAddr`

`Validators` is the primary index - it ensures that each operator can have only one
associated validator, where the public key of that validator can change in the
future. Delegators can refer to the immutable operator of the validator, without
concern for the changing public key.

`ValidatorByConsAddr` is an additional index that enables lookups for slashing.
When Tendermint reports evidence, it provides the validator address, so this
map is needed to find the operator. Note that the `ConsAddr` corresponds to the
address which can be derived from the validator's `ConsPubKey`.

Each validator's state is stored in a `Validator` struct:

```go
type Validator struct {
    OperatorAddress         sdk.ValAddress  // address of the validator's operator; bech encoded in JSON
    ConsPubKey              crypto.PubKey   // the consensus public key of the validator; bech encoded in JSON
    Description             Description     // description terms for the validator
}

type Description struct {
    Moniker          string // name
    Identity         string // optional identity signature (ex. UPort or Keybase)
    Website          string // optional website link
    SecurityContact  string // optional email for security contact
    Details          string // optional details
}
```
