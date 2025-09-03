package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Clear removes a participant's current vote and broadcasts VoteCleared on success.
func (s *Service) Clear(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID) error {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return fmt.Errorf("room not found: %s", roomID)
    }
    if err := room.ClearVote(participantID); err != nil {
        return fmt.Errorf("clear: %w", err)
    }
    if s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, VoteCleared{RoomID: roomID, ParticipantID: participantID}); err != nil {
            return fmt.Errorf("broadcast: %w", err)
        }
    }
    return nil
}

