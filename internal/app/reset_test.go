package app

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

type resetRepo struct {
	rooms map[domain.RoomID]*domain.Room
}

func (r *resetRepo) Create(ctx context.Context, room *domain.Room) error {
	if r.rooms == nil {
		r.rooms = make(map[domain.RoomID]*domain.Room)
	}
	r.rooms[room.ID()] = room
	return nil
}

func (r *resetRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
	rm, ok := r.rooms[id]
	return rm, ok, nil
}

func (r *resetRepo) Delete(ctx context.Context, id domain.RoomID) error {
	delete(r.rooms, id)
	return nil
}

type resetBus struct{ events []any }

func (b *resetBus) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
	b.events = append(b.events, event)
	return nil
}

func TestReset_IncrementsRound_ClearsVotes_Broadcasts(t *testing.T) {
	ctx := context.Background()
	repo := &resetRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	p1 := domain.ParticipantID("p1")
	p2 := domain.ParticipantID("p2")
	_ = room.Join(p1, "Alice")
	_ = room.Join(p2, "Bob")
	_ = room.CastVote(p1, "5")
	_ = room.CastVote(p2, "8")
	_ = room.Reveal()
	before := room.RoundIndex()

	bus := &resetBus{}
	svc := &Service{Rooms: repo, Bus: bus}
	if err := svc.Reset(ctx, roomID); err != nil {
		t.Fatalf("reset: %v", err)
	}

	if room.RoundIndex() != before+1 {
		t.Fatalf("round index not incremented")
	}
	if len(room.Votes()) != 0 {
		t.Fatalf("votes should be cleared")
	}
	if room.IsRevealed() {
		t.Fatalf("should be back to Voting state")
	}
	if len(bus.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(bus.events))
	}
	evt, ok := bus.events[0].(RoundReset)
	if !ok {
		t.Fatalf("wrong event type: %T", bus.events[0])
	}
	if evt.RoomID != roomID || evt.Round != room.RoundIndex() {
		t.Fatalf("event contents mismatch: %+v", evt)
	}
}
