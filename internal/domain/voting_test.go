package domain

import (
	"strings"
	"testing"
)

func TestRoom_CastVote_Rules(t *testing.T) {
	r := NewRoom(RoomID("r1"))
	must := func(err error) {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	}

	must(r.Join(ParticipantID("p1"), "Alice"))
	must(r.Join(ParticipantID("p2"), "Bob"))

	// Reject invalid card
	if err := r.CastVote(ParticipantID("p1"), "42"); err == nil {
		t.Fatalf("expected invalid card error, got nil")
	} else if !strings.Contains(strings.ToLower(err.Error()), "card") {
		t.Fatalf("expected card error, got: %v", err)
	}

	// Accept valid votes (number and special)
	must(r.CastVote(ParticipantID("p1"), "13"))
	must(r.CastVote(ParticipantID("p2"), "Pass"))

	// Non-member cannot vote
	if err := r.CastVote(ParticipantID("ghost"), "8"); err == nil {
		t.Fatalf("expected non-member error, got nil")
	}

	// Reveal locks votes
	must(r.Reveal())
	if err := r.CastVote(ParticipantID("p2"), "8"); err == nil {
		t.Fatalf("expected voting-closed error after reveal, got nil")
	}
}
