package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Leave removes a participant from the room and broadcasts ParticipantLeft.
func (s *Service) Leave(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID) error {
    room, err := s.getRoom(ctx, roomID)
    if err != nil { return err }
    if err := room.Leave(participantID); err != nil {
        return fmt.Errorf("leave: %w", err)
    }
    if err := s.emit(ctx, roomID, ParticipantLeft{RoomID: roomID, ParticipantID: participantID}); err != nil { return err }
    return nil
}
