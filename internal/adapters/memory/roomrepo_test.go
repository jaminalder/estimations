package memory

import (
    "context"
    "testing"

    "github.com/jaminalder/estimations/internal/domain"
)

func TestRoomRepo_CreateGetDelete_Basics(t *testing.T) {
    ctx := context.Background()
    repo := NewRoomRepo()
    id := domain.RoomID("r1")
    room := domain.NewRoom(id)

    // Create
    if err := repo.Create(ctx, room); err != nil {
        t.Fatalf("create: %v", err)
    }

    // Get
    got, ok, err := repo.Get(ctx, id)
    if err != nil { t.Fatalf("get: %v", err) }
    if !ok || got == nil { t.Fatalf("expected room to exist") }
    if got.ID() != id { t.Fatalf("id mismatch: %s", got.ID()) }

    // Pointer semantics: mutate via returned pointer
    if err := got.Join(domain.ParticipantID("p1"), "Alice"); err != nil {
        t.Fatalf("mutate via pointer: %v", err)
    }
    again, ok, _ := repo.Get(ctx, id)
    if !ok || len(again.Participants()) != 1 { t.Fatalf("expected mutation to persist in repo") }

    // Delete
    if err := repo.Delete(ctx, id); err != nil { t.Fatalf("delete: %v", err) }
    if _, ok, _ := repo.Get(ctx, id); ok { t.Fatalf("expected room to be gone after delete") }
}

func TestRoomRepo_Create_Duplicate(t *testing.T) {
    ctx := context.Background()
    repo := NewRoomRepo()
    id := domain.RoomID("dup")
    if err := repo.Create(ctx, domain.NewRoom(id)); err != nil {
        t.Fatalf("first create: %v", err)
    }
    if err := repo.Create(ctx, domain.NewRoom(id)); err == nil {
        t.Fatalf("expected duplicate create to error")
    }
}

