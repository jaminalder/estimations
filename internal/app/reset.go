package app

import (
	"context"
	"fmt"

	"github.com/jaminalder/estimations/internal/domain"
)

// Reset clears all votes, increments round, and reopens voting. Broadcasts RoundReset with new round index.
func (s *Service) Reset(ctx context.Context, roomID domain.RoomID) error {
	room, err := s.getRoom(ctx, roomID)
	if err != nil {
		return err
	}
	if err := room.Reset(); err != nil {
		return fmt.Errorf("reset: %w", err)
	}
	if err := s.emit(ctx, roomID, RoundReset{RoomID: roomID, Round: room.RoundIndex()}); err != nil {
		return err
	}
	return nil
}
