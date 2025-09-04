package app

import (
	"context"
	"fmt"

	"github.com/jaminalder/estimations/internal/domain"
)

// Reveal reveals the votes if at least one vote exists. Emits VotesRevealed once.
func (s *Service) Reveal(ctx context.Context, roomID domain.RoomID) error {
	room, err := s.getRoom(ctx, roomID)
	if err != nil {
		return err
	}
	wasRevealed := room.IsRevealed()
	if err := room.Reveal(); err != nil {
		return fmt.Errorf("reveal: %w", err)
	}
	if !wasRevealed {
		if err := s.emit(ctx, roomID, VotesRevealed{RoomID: roomID}); err != nil {
			return err
		}
	}
	return nil
}
