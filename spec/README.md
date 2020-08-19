<!--
order: 0
title: POA Overview
parent:
  title: "PoA"
-->

# `PoA`

## Abstract

The module enables a Cosmos-SDK based blockchain to use a Proof of Authority system to determine the validator set.

All validators in the system has an equal voting power.

Currently, this system allows a single master to arbitrarily and dynamically determine the effective validator set for the system. I plan implement a set of several masters with the MultSig interface to manage the validator set.

## Contents

1. **[State](01_state.md)**
2. **[Messages](02_messages.md)**
3. **[End-Block ](03_end_block.md)**
4. **[Hooks](04_hooks.md)**
5. **[Events](05_events.md)**
6. **[Parameters](06_params.md)**