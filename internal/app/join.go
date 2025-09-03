package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Join adds a participant with the given display name to the room and
// broadcasts a ParticipantJoined event upon success.
func (s *Service) Join(ctx context.Context, roomID domain.RoomID, name string) (domain.ParticipantID, error) {
    room, err := s.getRoom(ctx, roomID)
    if err != nil { return "", err }

    pid := s.Ids.NewParticipantID()
    if err := room.Join(pid, name); err != nil {
        return "", fmt.Errorf("join: %w", err)
    }
    if err := s.emit(ctx, roomID, ParticipantJoined{RoomID: roomID, ParticipantID: pid, Name: name}); err != nil { return "", err }
    return pid, nil
}
