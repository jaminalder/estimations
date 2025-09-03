package app

import (
    "context"
    "fmt"

    "github.com/jaminalder/estimations/internal/domain"
)

// CreateRoom creates a new room with a generated ID and persists it.
func (s *Service) CreateRoom(ctx context.Context) (domain.RoomID, error) {
    id := s.Ids.NewRoomID()
    room := domain.NewRoom(id)
    if err := s.Rooms.Create(ctx, room); err != nil {
        return "", fmt.Errorf("create room: %w", err)
    }
    return id, nil
}

