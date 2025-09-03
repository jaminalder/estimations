package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Reset clears all votes, increments round, and reopens voting. Broadcasts RoundReset with new round index.
func (s *Service) Reset(ctx context.Context, roomID domain.RoomID) error {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return fmt.Errorf("room not found: %s", roomID)
    }
    if err := room.Reset(); err != nil {
        return fmt.Errorf("reset: %w", err)
    }
    if s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, RoundReset{RoomID: roomID, Round: room.RoundIndex()}); err != nil {
            return fmt.Errorf("broadcast: %w", err)
        }
    }
    return nil
}

