package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// Join adds a participant with the given display name to the room and
// broadcasts a ParticipantJoined event upon success.
func (s *Service) Join(ctx context.Context, roomID domain.RoomID, name string) (domain.ParticipantID, error) {
    room, ok, err := s.Rooms.Get(ctx, roomID)
    if err != nil {
        return "", fmt.Errorf("get room: %w", err)
    }
    if !ok || room == nil {
        return "", fmt.Errorf("room not found: %s", roomID)
    }

    pid := s.Ids.NewParticipantID()
    if err := room.Join(pid, name); err != nil {
        return "", fmt.Errorf("join: %w", err)
    }
    // Broadcast outside the domain; return error if broadcasting fails
    if s.Bus != nil {
        if err := s.Bus.Broadcast(ctx, roomID, ParticipantJoined{RoomID: roomID, ParticipantID: pid, Name: name}); err != nil {
            return "", fmt.Errorf("broadcast: %w", err)
        }
    }
    return pid, nil
}
