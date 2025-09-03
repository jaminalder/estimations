package app

import (
    "context"
    "errors"
    "testing"

    "github.com/jaminalder/estimations/internal/domain"
)

// fake broadcaster capturing events
type captureBroadcaster struct{ events []any }

func (c *captureBroadcaster) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
    c.events = append(c.events, event)
    return nil
}

// in-memory repo fake (pointer semantics)
type joinRepo struct{ rooms map[domain.RoomID]*domain.Room }

func (r *joinRepo) Create(ctx context.Context, room *domain.Room) error {
    if r.rooms == nil { r.rooms = make(map[domain.RoomID]*domain.Room) }
    if _, exists := r.rooms[room.ID()]; exists {
        return errors.New("room exists")
    }
    r.rooms[room.ID()] = room
    return nil
}
func (r *joinRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
    room, ok := r.rooms[id]
    return room, ok, nil
}
func (r *joinRepo) Delete(ctx context.Context, id domain.RoomID) error {
    delete(r.rooms, id)
    return nil
}

type fixedIDs struct{ nextP domain.ParticipantID }

func (f fixedIDs) NewRoomID() domain.RoomID               { return "unused" }
func (f fixedIDs) NewParticipantID() domain.ParticipantID { return f.nextP }

func TestJoin_Success_BroadcastsEvent(t *testing.T) {
    ctx := context.Background()
    repo := &joinRepo{}
    roomID := domain.RoomID("r1")
    room := domain.NewRoom(roomID)
    if err := repo.Create(ctx, room); err != nil { t.Fatalf("setup create room: %v", err) }

    bus := &captureBroadcaster{}
    ids := fixedIDs{nextP: domain.ParticipantID("p1")}
    svc := &Service{Rooms: repo, Ids: ids, Bus: bus}

    pid, err := svc.Join(ctx, roomID, "Alice")
    if err != nil { t.Fatalf("join error: %v", err) }
    if pid != "p1" { t.Fatalf("participant id mismatch: got %s want p1", pid) }

    // Domain state updated
    if len(room.Participants()) != 1 { t.Fatalf("expected 1 participant") }
    if room.Participants()[0].Name != "Alice" { t.Fatalf("name mismatch: %s", room.Participants()[0].Name) }

    // Event captured
    if len(bus.events) != 1 {
        t.Fatalf("expected 1 event, got %d", len(bus.events))
    }
    evt, ok := bus.events[0].(ParticipantJoined)
    if !ok { t.Fatalf("wrong event type: %T", bus.events[0]) }
    if evt.RoomID != roomID || evt.ParticipantID != pid || evt.Name != "Alice" {
        t.Fatalf("event contents mismatch: %+v", evt)
    }
}

func TestJoin_RoomNotFound(t *testing.T) {
    ctx := context.Background()
    repo := &joinRepo{}
    bus := &captureBroadcaster{}
    ids := fixedIDs{nextP: domain.ParticipantID("p1")}
    svc := &Service{Rooms: repo, Ids: ids, Bus: bus}

    _, err := svc.Join(ctx, domain.RoomID("missing"), "Alice")
    if err == nil { t.Fatalf("expected error for missing room") }
    if len(bus.events) != 0 { t.Fatalf("no event should be broadcast on failure") }
}

func TestJoin_DuplicateName_Fails_NoBroadcast(t *testing.T) {
    ctx := context.Background()
    repo := &joinRepo{}
    roomID := domain.RoomID("r1")
    room := domain.NewRoom(roomID)
    if err := repo.Create(ctx, room); err != nil { t.Fatalf("create: %v", err) }

    bus := &captureBroadcaster{}
    ids := fixedIDs{nextP: domain.ParticipantID("p1")}
    svc := &Service{Rooms: repo, Ids: ids, Bus: bus}

    // First join ok
    if _, err := svc.Join(ctx, roomID, "Alice"); err != nil { t.Fatalf("join1: %v", err) }
    // Second join with same name should fail
    ids.nextP = "p2"
    // Replace ids generator in service for second call
    svc.Ids = ids
    if _, err := svc.Join(ctx, roomID, "alice"); err == nil {
        t.Fatalf("expected duplicate name error")
    }
    if len(bus.events) != 1 { t.Fatalf("only first join should broadcast; got %d", len(bus.events)) }
}
