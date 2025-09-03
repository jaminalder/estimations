package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Leave removes a participant from the room and broadcasts ParticipantLeft.
func (s *Service) Leave(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID) error {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return fmt.Errorf("room not found: %s", roomID)
    }
    if err := room.Leave(participantID); err != nil {
        return fmt.Errorf("leave: %w", err)
    }
    if s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, ParticipantLeft{RoomID: roomID, ParticipantID: participantID}); err != nil {
            return fmt.Errorf("broadcast: %w", err)
        }
    }
    return nil
}

