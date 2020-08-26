<!--
order: 6
-->

# Events

The poa module emits the following events:

## Handlers

### MsgSubmitApplication

**If Quorum > 0%:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| submit_application | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | submit_application               |
| message  | sender        | {senderAddress}    |

**If Quorum = 0%:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| append_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | submit_application               |
| message  | sender        | {senderAddress}    |

### MsgApproveApplication

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| approve_application | voter     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | approve_application               |
| message  | sender        | {senderAddress}    |

**If Quorum reached:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| append_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | approve_application               |
| message  | sender        | {senderAddress}    |

### MsgRejectApplication

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| reject_application | voter     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | reject_application               |
| message  | sender        | {senderAddress}    |

**If Quorum reached for rejection:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| reject_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | reject_application               |
| message  | sender        | {senderAddress}    |

### MsgLeaveValidatorSet

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| leave_validator_set | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | leave_validator_set               |
| message  | sender        | {senderAddress}    |

### MsgProposeKick

**If Quorum > 0%:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| propose_kick | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | propose_kick               |
| message  | sender        | {senderAddress}    |

**If Quorum = 0%:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| kick_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | propose_kick               |
| message  | sender        | {senderAddress}    |

### MsgApproveKickProposal

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| approve_kick_proposal | voter     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | approve_kick_proposal               |
| message  | sender        | {senderAddress}    |

**If Quorum reached:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| kick_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | approve_kick_proposal               |
| message  | sender        | {senderAddress}    |

### MsgRejectKickProposal

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| reject_kick_proposal | voter     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | reject_kick_proposal               |
| message  | sender        | {senderAddress}    |

**If Quorum reached:**

| Type     | Attribute Key | Attribute Value    |
|----------|---------------|--------------------|
| keep_validator | validator     | {validatorAddress} |
| message  | module        | poa               |
| message  | action        | reject_kick_proposal               |
| message  | sender        | {senderAddress}    |
