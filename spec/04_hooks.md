<!--
order: 4
-->

# Hooks

Other modules may register operations to execute when a new validator is appended or removed.  These events can be registered to execute either
right `Before` or `After` the event (as per the hook name). The
following hooks can registered: 

 - `AfterValidatorCreated(Context, ValAddress)`
   - called when a validator is appenned
 - `BeforeValidatorDescriptionModified(Context, ValAddress)`
   - called when a validator's description is changed
 - `AfterValidatorRemoved(Context, ConsAddress, ValAddress)`
   - called when a validator is removed