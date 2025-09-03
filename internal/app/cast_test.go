package app

import (
    "context"
    "testing"

    "github.com/jaminalder/estimations/internal/domain"
)

// Reuse joinRepo and captureBroadcaster from join_test.go style
type castRepo struct{ rooms map[domain.RoomID]*domain.Room }

func (r *castRepo) Create(ctx context.Context, room *domain.Room) error {
    if r.rooms == nil { r.rooms = make(map[domain.RoomID]*domain.Room) }
    r.rooms[room.ID()] = room
    return nil
}
func (r *castRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
    rm, ok := r.rooms[id]
    return rm, ok, nil
}
func (r *castRepo) Delete(ctx context.Context, id domain.RoomID) error { delete(r.rooms, id); return nil }

type eventsSink struct{ events []any }
func (e *eventsSink) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error { e.events = append(e.events, event); return nil }

func TestCast_Success_BroadcastsEvent(t *testing.T) {
    ctx := context.Background()
    repo := &castRepo{}
    roomID := domain.RoomID("r1")
    room := domain.NewRoom(roomID)
    if err := repo.Create(ctx, room); err != nil { t.Fatalf("create: %v", err) }
    // Seed participant directly in domain
    pid := domain.ParticipantID("p1")
    if err := room.Join(pid, "Alice"); err != nil { t.Fatalf("seed join: %v", err) }

    bus := &eventsSink{}
    svc := &Service{Rooms: repo, Bus: bus}

    if err := svc.Cast(ctx, roomID, pid, "8"); err != nil { t.Fatalf("cast: %v", err) }

    if len(bus.events) != 1 { t.Fatalf("expected 1 event, got %d", len(bus.events)) }
    evt, ok := bus.events[0].(VoteCast)
    if !ok { t.Fatalf("wrong event type: %T", bus.events[0]) }
    if evt.RoomID != roomID || evt.ParticipantID != pid || evt.Card != "8" {
        t.Fatalf("event mismatch: %+v", evt)
    }
}

func TestCast_InvalidCard_NoBroadcast(t *testing.T) {
    ctx := context.Background()
    repo := &castRepo{}
    roomID := domain.RoomID("r1")
    room := domain.NewRoom(roomID)
    _ = repo.Create(ctx, room)
    pid := domain.ParticipantID("p1")
    if err := room.Join(pid, "Alice"); err != nil { t.Fatalf("seed join: %v", err) }

    bus := &eventsSink{}
    svc := &Service{Rooms: repo, Bus: bus}

    if err := svc.Cast(ctx, roomID, pid, "42"); err == nil {
        t.Fatalf("expected invalid card error")
    }
    if len(bus.events) != 0 { t.Fatalf("no event should be broadcast on failure") }
}

func TestCast_NotParticipant_NoBroadcast(t *testing.T) {
    ctx := context.Background()
    repo := &castRepo{}
    roomID := domain.RoomID("r1")
    room := domain.NewRoom(roomID)
    _ = repo.Create(ctx, room)

    bus := &eventsSink{}
    svc := &Service{Rooms: repo, Bus: bus}

    if err := svc.Cast(ctx, roomID, domain.ParticipantID("ghost"), "8"); err == nil {
        t.Fatalf("expected non-participant error")
    }
    if len(bus.events) != 0 { t.Fatalf("no event should be broadcast on failure") }
}

