package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Cast records a participant's vote in the room and broadcasts VoteCast on success.
func (s *Service) Cast(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID, card string) error {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return fmt.Errorf("room not found: %s", roomID)
    }
    if err := room.CastVote(participantID, card); err != nil {
        return fmt.Errorf("cast: %w", err)
    }
    if s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, VoteCast{RoomID: roomID, ParticipantID: participantID, Card: card}); err != nil {
            return fmt.Errorf("broadcast: %w", err)
        }
    }
    return nil
}

