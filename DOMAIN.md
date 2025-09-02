# DOMAIN.md — Estimation Poker Domain

Last updated: 2025-09-02

## Purpose
Authoritative domain model for the in-memory, SSR estimation poker app. Drives use-cases, tests, and adapters.

## Entities
- Room: aggregate root; holds participants, fixed deck, current round, and state.
- Participant: display name + ParticipantID; belongs to exactly one room (session-scoped).
- Round: current-only; tracks votes and state; increments on reset.

## Value Objects
- RoomID, ParticipantID: opaque identifiers.
- Deck: fixed set for v1 — Fibonacci cards `[0,1,2,3,5,8,13,21,34]` plus specials `["?", "∞", "☕", "Pass"]`.
- Card/Vote: a chosen card from the deck; vote can be unset.

## Relationships
- Room → Participants: 1..25. Names unique per room (case-insensitive).
- Room → Deck: exactly 1, immutable in v1.
- Room → current Round: exactly 1; Round maps `ParticipantID → Vote`.

## States & Lifecycle
- Round state machine: `Voting → Revealed → (Reset) → Voting (new round index)`.
- Flow: create room → join participants → cast/clear votes while Voting → reveal (requires ≥1 vote) → reset (clears votes, new round).
- Leave: removes participant immediately and deletes their vote.

## Behaviors (commands)
- Room.Join(name) → ParticipantID
- Room.Leave(participantID)
- Room.CastVote(participantID, card)
- Room.ClearVote(participantID)
- Room.Reveal()
- Room.Reset() // starts a new round with no votes

## Invariants & Rules
- Names: trimmed, non-empty, unique per room (case-insensitive). Duplicate join is rejected.
- Capacity: maximum 25 participants per room.
- Voting: only joined participants can vote; exactly one current vote per participant; votes are mutable only while state=Voting.
- Card validity: vote card must exist in the current deck (including specials like "Pass", "?", "∞", "☕").
- Reveal: allowed only if at least one vote exists (specials count toward the threshold). After reveal, votes are locked (no cast/clear).
- Reset: clears all votes, increments round index, sets state=Voting; deck remains unchanged.
- Deck: immutable in v1 (single built-in deck defined above).

## Domain Events (for SSE bridge)
- ParticipantJoined, ParticipantLeft
- VoteCast, VoteCleared
- VotesRevealed
- RoundReset

## Defaults & Omissions (v1)
- No ownership/admin/permissions; any participant may reveal/reset.
- No multi-round history persistence (in-memory only; app may log externally if needed).
- No round label/metadata (deferred).
- One browser session = one participant; no multi-tab/session consolidation.

## Open Integration Concerns (outside domain)
- Name length/character policy: suggest 1..32 chars; allow unicode; enforce in adapter.
- Room GC/TTL when empty: handled by app layer; domain agnostic.
- Participant ordering for UI: stable by join time; adapter concern.

