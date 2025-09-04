package app

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

type repoMem struct {
	rooms map[domain.RoomID]*domain.Room
}

func (r *repoMem) Create(ctx context.Context, room *domain.Room) error {
	if r.rooms == nil {
		r.rooms = make(map[domain.RoomID]*domain.Room)
	}
	r.rooms[room.ID()] = room
	return nil
}

func (r *repoMem) Get(ctx context.Context, id domain.RoomID) (*domain.Room, bool, error) {
	rm, ok := r.rooms[id]
	return rm, ok, nil
}

func (r *repoMem) Delete(ctx context.Context, id domain.RoomID) error {
	delete(r.rooms, id)
	return nil
}

type idFixed struct{ rid domain.RoomID }

func (i idFixed) NewRoomID() domain.RoomID               { return i.rid }
func (i idFixed) NewParticipantID() domain.ParticipantID { return "p-fixed" }

func TestCreateRoom_Basics(t *testing.T) {
	ctx := context.Background()
	repo := &repoMem{}
	ids := idFixed{rid: domain.RoomID("room-123")}

	svc := &Service{Rooms: repo, Ids: ids}
	gotID, err := svc.CreateRoom(ctx)
	if err != nil {
		t.Fatalf("CreateRoom error: %v", err)
	}
	if gotID != domain.RoomID("room-123") {
		t.Fatalf("unexpected id: got %s want %s", gotID, "room-123")
	}

	rm, ok, _ := repo.Get(ctx, gotID)
	if !ok || rm == nil {
		t.Fatalf("room not found in repo after creation")
	}
	if rm.IsRevealed() {
		t.Fatalf("new room should start in Voting (not revealed)")
	}
	if rm.RoundIndex() != 0 {
		t.Fatalf("new room round index should be 0")
	}
	if len(rm.Participants()) != 0 {
		t.Fatalf("new room should have no participants")
	}

	expectedDeck := []string{"0", "1", "2", "3", "5", "8", "13", "21", "34", "?", "∞", "☕", "Pass"}
	deck := rm.Deck()
	if len(deck) != len(expectedDeck) {
		t.Fatalf("deck length mismatch: got %d want %d", len(deck), len(expectedDeck))
	}
	for i := range deck {
		if deck[i] != expectedDeck[i] {
			t.Fatalf("deck mismatch at %d: got %s want %s", i, deck[i], expectedDeck[i])
		}
	}
}
