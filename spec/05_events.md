<!--
order: 5
-->

# Events

The PoA module emits the following events:

## Handlers

### MsgAppendValidator

| Type             | Attribute Key | Attribute Value    |
| ---------------- | ------------- | ------------------ |
| append_validator | validator     | {validatorAddress} |
| message          | module        | poa            |
| message          | action        | append_validator   |
| message          | sender        | {senderAddress}    |

### MsgEditValidatorDescription

| Type           | Attribute Key       | Attribute Value     |
| -------------- | ------------------- | ------------------- |
| edit_validator_description | validator     | {validatorAddress} |
| message        | module              | poa             |
| message        | action              | edit_validator_description      |
| message        | sender              | {senderAddress}     |

### MsgRemoveValidator

| Type     | Attribute Key | Attribute Value    |
| -------- | ------------- | ------------------ |
| remove_validator | validator     | {validatorAddress} |
| message  | module        | poa            |
| message  | action        | remove_validator           |
| message  | sender        | {senderAddress}    |
