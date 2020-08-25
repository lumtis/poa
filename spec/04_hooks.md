<!--
order: 4
-->

# Hooks

Other modules may register operations to execute when an event has occurred within the module (new application, new validator,...). These events can be registered to execute either right `After` the event. The following hooks can be registered: 

- `AfterApplicationSubmitted(Context, ValAddress)`
  - called when a new application to become validator has been submitted
- `AfterApplicationApproved(Context, ValAddress)`
  - called when an application has been approved and a new validator is appended into the validator set
- `AfterApplicationRejected(Context, ValAddress)`
  - called when an application has been rejected
- `AfterValidatorLeft(Context, ValAddress)`
  - called when a validator left the validator set
- `AfterKickProposed(Context, ValAddress)`
  - called when a new kick proposal has been created
- `AfterKickProposalApproved(Context, ValAddress)`
  - called when a kick proposal has been approved and the validator is removed from the validator set
- `AfterKickProposalRejected(Context, ValAddress)`
  - called when a kick proposal has been rejected