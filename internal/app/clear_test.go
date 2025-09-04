package app

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

type clearRepo struct {
	rooms map[domain.RoomID]*domain.Room
}

func (r *clearRepo) Create(ctx context.Context, room *domain.Room) error {
	if r.rooms == nil {
		r.rooms = make(map[domain.RoomID]*domain.Room)
	}
	r.rooms[room.ID()] = room
	return nil
}

func (r *clearRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
	rm, ok := r.rooms[id]
	return rm, ok, nil
}

func (r *clearRepo) Delete(ctx context.Context, id domain.RoomID) error {
	delete(r.rooms, id)
	return nil
}

type clearBus struct{ events []any }

func (b *clearBus) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
	b.events = append(b.events, event)
	return nil
}

func TestClear_Success_BroadcastsEvent(t *testing.T) {
	ctx := context.Background()
	repo := &clearRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	pid := domain.ParticipantID("p1")
	if err := room.Join(pid, "Alice"); err != nil {
		t.Fatalf("seed join: %v", err)
	}
	if err := room.CastVote(pid, "8"); err != nil {
		t.Fatalf("seed cast: %v", err)
	}

	bus := &clearBus{}
	svc := &Service{Rooms: repo, Bus: bus}

	if err := svc.Clear(ctx, roomID, pid); err != nil {
		t.Fatalf("clear: %v", err)
	}
	// Vote is gone
	if _, ok := room.Votes()[pid]; ok {
		t.Fatalf("expected vote to be cleared")
	}
	// Event emitted
	if len(bus.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(bus.events))
	}
	if _, ok := bus.events[0].(VoteCleared); !ok {
		t.Fatalf("wrong event type: %T", bus.events[0])
	}
}

func TestClear_NotParticipant_NoBroadcast(t *testing.T) {
	ctx := context.Background()
	repo := &clearRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)

	bus := &clearBus{}
	svc := &Service{Rooms: repo, Bus: bus}

	if err := svc.Clear(ctx, roomID, domain.ParticipantID("ghost")); err == nil {
		t.Fatalf("expected error for non-participant clear")
	}
	if len(bus.events) != 0 {
		t.Fatalf("no broadcast expected on failure")
	}
}

func TestClear_WhileRevealed_NoBroadcast(t *testing.T) {
	ctx := context.Background()
	repo := &clearRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	pid := domain.ParticipantID("p1")
	if err := room.Join(pid, "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := room.CastVote(pid, "5"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := room.Reveal(); err != nil {
		t.Fatalf("reveal: %v", err)
	}

	bus := &clearBus{}
	svc := &Service{Rooms: repo, Bus: bus}
	if err := svc.Clear(ctx, roomID, pid); err == nil {
		t.Fatalf("expected error while revealed")
	}
	if len(bus.events) != 0 {
		t.Fatalf("no broadcast expected on failure")
	}
}
