package idgen

import (
    "regexp"
    "testing"
)

func TestRandomIdGen_FormatAndUniqueness(t *testing.T) {
    g := NewRandom(10, 8)
    re := regexp.MustCompile(`^[A-Za-z0-9_-]+$`)

    // Room IDs
    seen := make(map[string]struct{})
    for i := 0; i < 200; i++ {
        id := string(g.NewRoomID())
        if len(id) < 12 { // expect around 16 chars for 10 bytes
            t.Fatalf("room id too short: %q (len=%d)", id, len(id))
        }
        if !re.MatchString(id) {
            t.Fatalf("room id not URL-safe: %q", id)
        }
        if _, dup := seen[id]; dup {
            t.Fatalf("duplicate room id generated: %q", id)
        }
        seen[id] = struct{}{}
    }

    // Participant IDs
    seen = make(map[string]struct{})
    for i := 0; i < 200; i++ {
        id := string(g.NewParticipantID())
        if len(id) < 10 { // expect around 11 chars for 8 bytes
            t.Fatalf("participant id too short: %q (len=%d)", id, len(id))
        }
        if !re.MatchString(id) {
            t.Fatalf("participant id not URL-safe: %q", id)
        }
        if _, dup := seen[id]; dup {
            t.Fatalf("duplicate participant id generated: %q", id)
        }
        seen[id] = struct{}{}
    }
}

