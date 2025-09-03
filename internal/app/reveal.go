package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Reveal reveals the votes if at least one vote exists. Emits VotesRevealed once.
func (s *Service) Reveal(ctx context.Context, roomID domain.RoomID) error {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return fmt.Errorf("room not found: %s", roomID)
    }
    wasRevealed := room.IsRevealed()
    if err := room.Reveal(); err != nil {
        return fmt.Errorf("reveal: %w", err)
    }
    if !wasRevealed && s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, VotesRevealed{RoomID: roomID}); err != nil {
            return fmt.Errorf("broadcast: %w", err)
        }
    }
    return nil
}

