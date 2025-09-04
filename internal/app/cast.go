package app

import (
	"context"
	"fmt"

	"github.com/jaminalder/estimations/internal/domain"
)

// Cast records a participant's vote in the room and broadcasts VoteCast on success.
func (s *Service) Cast(ctx context.Context, roomID domain.RoomID, participantID domain.ParticipantID, card string) error {
	room, err := s.getRoom(ctx, roomID)
	if err != nil {
		return err
	}
	if err := room.CastVote(participantID, card); err != nil {
		return fmt.Errorf("cast: %w", err)
	}
	if err := s.emit(ctx, roomID, VoteCast{RoomID: roomID, ParticipantID: participantID, Card: card}); err != nil {
		return err
	}
	return nil
}
