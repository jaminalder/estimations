package domain

import (
    "strings"
    "testing"
)

func TestRoom_Join_UniqueNames_Capacity(t *testing.T) {
    room := NewRoom(RoomID("r1"))

    // Join first participant
    if err := room.Join(ParticipantID("p1"), "Alice"); err != nil {
        t.Fatalf("join Alice: unexpected error: %v", err)
    }

    // Duplicate name (case-insensitive) should be rejected
    if err := room.Join(ParticipantID("p2"), "alice"); err == nil {
        t.Fatalf("expected duplicate name rejection, got nil error")
    } else if !strings.Contains(strings.ToLower(err.Error()), "duplicate") {
        t.Fatalf("expected duplicate error, got: %v", err)
    }

    // Fill up to capacity
    for i := 2; i <= 25; i++ { // already have 1 participant
        name := "User" + string(rune('A'+i))
        if err := room.Join(ParticipantID("p"+string(rune('A'+i))), name); err != nil {
            t.Fatalf("join %s: unexpected error: %v", name, err)
        }
    }

    // Exceed capacity
    if err := room.Join(ParticipantID("p_over"), "Zed"); err == nil {
        t.Fatalf("expected capacity error, got nil error")
    }
}

