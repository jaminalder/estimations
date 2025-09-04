package app

import (
	"context"
	"fmt"

	"github.com/jaminalder/estimations/internal/domain"
)

func (s *Service) getRoom(ctx context.Context, roomID domain.RoomID) (*domain.Room, error) {
	room, ok, err := s.Rooms.Get(ctx, roomID)
	if err != nil {
		return nil, fmt.Errorf("get room: %w", err)
	}
	if !ok || room == nil {
		return nil, fmt.Errorf("room not found: %s", roomID)
	}
	return room, nil
}

func (s *Service) emit(ctx context.Context, roomID domain.RoomID, event any) error {
	if s.Bus == nil {
		return nil
	}
	if err := s.Bus.Broadcast(ctx, roomID, event); err != nil {
		return fmt.Errorf("broadcast: %w", err)
	}
	return nil
}
