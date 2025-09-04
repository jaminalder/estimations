package domain

import (
	"testing"
)

func TestRoom_Reveal_Rules(t *testing.T) {
	r := NewRoom(RoomID("r1"))

	// Reveal without any votes should error
	if err := r.Reveal(); err == nil {
		t.Fatalf("expected error on reveal with no votes")
	}

	// Add a participant and a vote, then reveal should succeed
	if err := r.Join(ParticipantID("p1"), "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := r.CastVote(ParticipantID("p1"), "8"); err != nil {
		t.Fatalf("cast vote: %v", err)
	}
	if err := r.Reveal(); err != nil {
		t.Fatalf("reveal after vote: %v", err)
	}

	// Idempotent: second reveal is a no-op
	if err := r.Reveal(); err != nil {
		t.Fatalf("second reveal should be idempotent: %v", err)
	}
}

func TestRoom_Reset_ClearsAndReopens(t *testing.T) {
	r := NewRoom(RoomID("r1"))
	if err := r.Join(ParticipantID("p1"), "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := r.Join(ParticipantID("p2"), "Bob"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := r.CastVote(ParticipantID("p1"), "13"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.CastVote(ParticipantID("p2"), "Pass"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.Reveal(); err != nil {
		t.Fatalf("reveal: %v", err)
	}

	before := r.RoundIndex()
	if err := r.Reset(); err != nil {
		t.Fatalf("reset: %v", err)
	}
	after := r.RoundIndex()
	if after != before+1 {
		t.Fatalf("round index not incremented: before=%d after=%d", before, after)
	}

	// After reset, votes cleared: reveal should now fail until someone votes
	if err := r.Reveal(); err == nil {
		t.Fatalf("expected error on reveal with no votes after reset")
	}

	// And voting is open again
	if err := r.CastVote(ParticipantID("p2"), "8"); err != nil {
		t.Fatalf("cast after reset should succeed: %v", err)
	}
}
