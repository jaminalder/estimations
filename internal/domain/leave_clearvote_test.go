package domain

import "testing"

func TestRoom_Leave_RemovesVoteAndMember(t *testing.T) {
	r := NewRoom(RoomID("r1"))
	if err := r.Join(ParticipantID("p1"), "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	if err := r.Join(ParticipantID("p2"), "Bob"); err != nil {
		t.Fatalf("join: %v", err)
	}

	// p1 votes, then leaves → vote removed and p1 cannot vote anymore
	if err := r.CastVote(ParticipantID("p1"), "8"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.Leave(ParticipantID("p1")); err != nil {
		t.Fatalf("leave: %v", err)
	}

	if err := r.CastVote(ParticipantID("p1"), "13"); err == nil {
		t.Fatalf("expected error: p1 should no longer be a participant")
	}

	// With p1 gone and p2 not voted, reveal should fail (no votes present)
	if err := r.Reveal(); err == nil {
		t.Fatalf("expected reveal error due to no votes after leave")
	}

	// Room still functions: p2 can vote and reveal
	if err := r.CastVote(ParticipantID("p2"), "5"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.Reveal(); err != nil {
		t.Fatalf("reveal: %v", err)
	}
}

func TestRoom_ClearVote_Rules(t *testing.T) {
	r := NewRoom(RoomID("r1"))
	if err := r.Join(ParticipantID("p1"), "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}

	// Cast, then clear while Voting → reveal fails
	if err := r.CastVote(ParticipantID("p1"), "8"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.ClearVote(ParticipantID("p1")); err != nil {
		t.Fatalf("clear: %v", err)
	}
	if err := r.Reveal(); err == nil {
		t.Fatalf("expected reveal error after clearing all votes")
	}

	// Cast again and reveal; then clear should fail while Revealed
	// Re-open voting by casting then revealing (we need ≥1 vote)
	if err := r.CastVote(ParticipantID("p1"), "13"); err != nil {
		t.Fatalf("cast: %v", err)
	}
	if err := r.Reveal(); err != nil {
		t.Fatalf("reveal: %v", err)
	}
	if err := r.ClearVote(ParticipantID("p1")); err == nil {
		t.Fatalf("expected clear to fail while revealed")
	}
}
