package app

import (
    "context"
    "testing"
    "time"

    "github.com/jaminalder/estimations/internal/domain"
)

// Compile-time conformance checks for ports.
func TestPorts_CompileTypes(t *testing.T) {
    var _ RoomRepo = &fakeRepo{}
    var _ IdGen = fakeIDGen{}
    var _ Clock = fakeClock{}
    var _ Broadcaster = fakeBroadcaster{}

    // Use the vars a bit to avoid unused warnings in future expansions
    _ = context.Background()
}

type fakeRepo struct{}

func (f *fakeRepo) Create(ctx context.Context, room *domain.Room) error { return nil }
func (f *fakeRepo) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
    return nil, false, nil
}
func (f *fakeRepo) Delete(ctx context.Context, id domain.RoomID) error { return nil }

type fakeIDGen struct{}

func (f fakeIDGen) NewRoomID() domain.RoomID          { return domain.RoomID("r") }
func (f fakeIDGen) NewParticipantID() domain.ParticipantID { return domain.ParticipantID("p") }

type fakeClock struct{}

func (f fakeClock) Now() time.Time { return time.Time{} }

type fakeBroadcaster struct{}

func (f fakeBroadcaster) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
    return nil
}
