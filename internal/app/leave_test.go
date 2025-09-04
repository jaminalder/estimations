package app

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

type leaveRepo struct {
	rooms map[domain.RoomID]*domain.Room
}

func (r *leaveRepo) Create(ctx context.Context, room *domain.Room) error {
	if r.rooms == nil {
		r.rooms = make(map[domain.RoomID]*domain.Room)
	}
	r.rooms[room.ID()] = room
	return nil
}

func (r *leaveRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
	rm, ok := r.rooms[id]
	return rm, ok, nil
}

func (r *leaveRepo) Delete(ctx context.Context, id domain.RoomID) error {
	delete(r.rooms, id)
	return nil
}

type leaveBus struct{ events []any }

func (b *leaveBus) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
	b.events = append(b.events, event)
	return nil
}

func TestLeave_RemovesParticipantAndVote_Broadcasts(t *testing.T) {
	ctx := context.Background()
	repo := &leaveRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	p1 := domain.ParticipantID("p1")
	_ = room.Join(p1, "Alice")
	_ = room.CastVote(p1, "8")
	bus := &leaveBus{}
	svc := &Service{Rooms: repo, Bus: bus}

	if err := svc.Leave(ctx, roomID, p1); err != nil {
		t.Fatalf("leave: %v", err)
	}
	if len(room.Participants()) != 0 {
		t.Fatalf("participant should be removed")
	}
	if len(room.Votes()) != 0 {
		t.Fatalf("vote should be removed on leave")
	}
	if len(bus.events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(bus.events))
	}
	if _, ok := bus.events[0].(ParticipantLeft); !ok {
		t.Fatalf("wrong event type: %T", bus.events[0])
	}
}

func TestLeave_NotParticipant_NoBroadcast(t *testing.T) {
	ctx := context.Background()
	repo := &leaveRepo{}
	roomID := domain.RoomID("r1")
	room := domain.NewRoom(roomID)
	_ = repo.Create(ctx, room)
	bus := &leaveBus{}
	svc := &Service{Rooms: repo, Bus: bus}

	if err := svc.Leave(ctx, roomID, domain.ParticipantID("ghost")); err == nil {
		t.Fatalf("expected error for non-participant")
	}
	if len(bus.events) != 0 {
		t.Fatalf("no broadcast on failure")
	}
}
