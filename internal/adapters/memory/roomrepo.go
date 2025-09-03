package memory

import (
    "context"
    "fmt"
    "sync"

    "github.com/jaminalder/estimations/internal/app"
    "github.com/jaminalder/estimations/internal/domain"
)

// RoomRepo is an in-memory implementation of app.RoomRepo.
type RoomRepo struct {
    mu    sync.RWMutex
    rooms map[domain.RoomID]*domain.Room
}

func NewRoomRepo() *RoomRepo {
    return &RoomRepo{rooms: make(map[domain.RoomID]*domain.Room)}
}

var _ app.RoomRepo = (*RoomRepo)(nil)

func (r *RoomRepo) Create(ctx context.Context, room *domain.Room) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    id := room.ID()
    if _, exists := r.rooms[id]; exists {
        return fmt.Errorf("room exists: %s", id)
    }
    r.rooms[id] = room
    return nil
}

func (r *RoomRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    rm, ok := r.rooms[id]
    return rm, ok, nil
}

func (r *RoomRepo) Delete(ctx context.Context, id domain.RoomID) error {
    r.mu.Lock()
    defer r.mu.Unlock()
    delete(r.rooms, id)
    return nil
}

