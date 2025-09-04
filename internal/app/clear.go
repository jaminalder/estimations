package app

import (
	"context"
	"fmt"

	"github.com/jaminalder/estimations/internal/domain"
)

// Clear removes a participant's current vote and broadcasts VoteCleared on success.
func (s *Service) Clear(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID) error {
	room, err := s.getRoom(ctx, roomID)
	if err != nil {
		return err
	}
	if err := room.ClearVote(participantID); err != nil {
		return fmt.Errorf("clear: %w", err)
	}
	if err := s.emit(ctx, roomID, VoteCleared{RoomID: roomID, ParticipantID: participantID}); err != nil {
		return err
	}
	return nil
}
