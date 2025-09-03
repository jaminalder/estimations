package app

import (
    "context"
    "time"

    "github.com/jaminalder/estimations/internal/domain"
)

// RoomRepo is an in-memory repository interface for Room aggregates.
type RoomRepo interface {
    Create(ctx context.Context, room *domain.Room) error
    Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error)
    Delete(ctx context.Context, id domain.RoomID) error
}

// IdGen provides opaque identifiers for rooms and participants.
type IdGen interface {
    NewRoomID() domain.RoomID
    NewParticipantID() domain.ParticipantID
}

// Clock supplies time for TTLs/metadata at the app layer.
type Clock interface {
    Now() time.Time
}

// Broadcaster emits room-scoped events (bridged to SSE by adapter).
type Broadcaster interface {
    Broadcast(ctx context.Context, roomID domain.RoomID, event any) error
}

