package app

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

type revealRepo struct {
	rooms map[domain.RoomID]*domain.Room
}

func (r *revealRepo) Create(ctx context.Context, room *domain.Room) error {
	if r.rooms == nil {
		r.rooms = make(map[domain.RoomID]*domain.Room)
	}
	r.rooms[room.ID()] = room
	return nil
}

func (r *revealRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
	rm, ok := r.rooms[id]
	return rm, ok, nil
}

func (r *revealRepo) Delete(ctx context.Context, id domain.RoomID) error {
	delete(r.rooms, id)
	return nil
}

type revealBus struct{ events []any }

func (b *revealBus) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
	b.events = append(b.events, event)
	return nil
}

func TestReveal_Success_EmitsOnce(t *testing.T) {
	ctx := context.Background()
	repo := &revealRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	pid := domain.ParticipantID("p1")
	if err := room.Join(pid, "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := room.CastVote(pid, "8"); err != nil {
		t.Fatalf("cast: %v", err)
	}

	bus := &revealBus{}
	svc := &Service{Rooms: repo, Bus: bus}

	if err := svc.Reveal(ctx, roomID); err != nil {
		t.Fatalf("reveal: %v", err)
	}
	if len(bus.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(bus.events))
	}
	if _, ok := bus.events[0].(VotesRevealed); !ok {
		t.Fatalf("wrong event: %T", bus.events[0])
	}

	// Idempotent: calling again should not emit a second event
	if err := svc.Reveal(ctx, roomID); err != nil {
		t.Fatalf("second reveal should be ok: %v", err)
	}
	if len(bus.events) != 1 {
		t.Fatalf("expected only 1 event after idempotent reveal, got %d", len(bus.events))
	}
}

func TestReveal_NoVotes_NoBroadcast(t *testing.T) {
	ctx := context.Background()
	repo := &revealRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	_ = room.Join(domain.ParticipantID("p1"), "Alice")

	bus := &revealBus{}
	svc := &Service{Rooms: repo, Bus: bus}
	if err := svc.Reveal(ctx, roomID); err == nil {
		t.Fatalf("expected reveal to fail without votes")
	}
	if len(bus.events) != 0 {
		t.Fatalf("no event expected on failure")
	}
}
