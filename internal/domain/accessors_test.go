package domain

import (
	"reflect"
	"testing"
)

func TestRoom_Accessors_SnapshotsAndDeck(t *testing.T) {
	r := NewRoom(RoomID("r1"))

	// Initially
	if r.IsRevealed() {
		t.Fatalf("expected not revealed initially")
	}
	if r.RoundIndex() != 0 {
		t.Fatalf("expected round index 0 initially, got %d", r.RoundIndex())
	}
	if r.ID() != RoomID("r1") {
		t.Fatalf("expected room id r1, got %s", r.ID())
	}

	// Join two participants
	if err := r.Join(ParticipantID("p1"), "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := r.Join(ParticipantID("p2"), "Bob"); err != nil {
		t.Fatalf("join: %v", err)
	}

	ps := r.Participants()
	if len(ps) != 2 {
		t.Fatalf("expected 2 participants, got %d", len(ps))
	}

	// Cast a vote and verify votes snapshot is a copy
	if err := r.CastVote(ParticipantID("p1"), "8"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	votesCopy := r.Votes()
	// Mutate the returned map; should not affect internal state
	delete(votesCopy, ParticipantID("p1"))
	if err := r.Reveal(); err != nil {
		t.Fatalf("reveal should still succeed (internal vote intact): %v", err)
	}

	// Deck exposure
	expectedDeck := []string{"0", "1", "2", "3", "5", "8", "13", "21", "34", "?", "∞", "☕", "Pass"}
	if !reflect.DeepEqual(r.Deck(), expectedDeck) {
		t.Fatalf("deck mismatch.\nwant: %v\n got: %v", expectedDeck, r.Deck())
	}
}
