package app

import "github.com/jaminalder/estimations/internal/domain"

// ParticipantJoined is emitted after a participant successfully joins a room.
type ParticipantJoined struct {
    RoomID        domain.RoomID
    ParticipantID domain.ParticipantID
    Name          string
}

// ParticipantLeft is emitted after a participant leaves.
type ParticipantLeft struct {
    RoomID        domain.RoomID
    ParticipantID domain.ParticipantID
}

// VoteCast is emitted when a participant casts a vote.
type VoteCast struct {
    RoomID        domain.RoomID
    ParticipantID domain.ParticipantID
    Card          string
}

// VoteCleared is emitted when a participant clears their vote.
type VoteCleared struct {
    RoomID        domain.RoomID
    ParticipantID domain.ParticipantID
}

// VotesRevealed is emitted when votes are revealed for the current round.
type VotesRevealed struct {
    RoomID domain.RoomID
}

// RoundReset is emitted when a new round starts (after reset).
type RoundReset struct {
    RoomID domain.RoomID
    Round  int
}

