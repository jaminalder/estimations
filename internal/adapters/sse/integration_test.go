package sse

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/adapters/memory"
	"github.com/jaminalder/estimations/internal/app"
	"github.com/jaminalder/estimations/internal/domain"
)

type idsFixed struct{ pid domain.ParticipantID }

func (i idsFixed) NewRoomID() domain.RoomID               { return "unused" }
func (i idsFixed) NewParticipantID() domain.ParticipantID { return i.pid }

// Integration: use-case -> hub broadcast -> subscriber receives JSON payload.
func TestIntegration_JoinBroadcastsToSSE(t *testing.T) {
	ctx := context.Background()
	repo := memory.NewRoomRepo()
	hub := NewHub(4)
	svc := &app.Service{Rooms: repo, Ids: idsFixed{pid: "p1"}, Bus: hub}

	roomID := domain.RoomID("room-1")
	if err := repo.Create(ctx, domain.NewRoom(roomID)); err != nil {
		t.Fatalf("create room: %v", err)
	}

	ch, unsubscribe := hub.Subscribe(roomID)
	defer unsubscribe()

	if _, err := svc.Join(ctx, roomID, "Alice"); err != nil {
		t.Fatalf("join: %v", err)
	}
	payload := <-ch
	if len(payload) == 0 {
		t.Fatalf("expected payload bytes from hub")
	}
	// quick check payload mentions the participant id and name
	if string(payload) == "" || !containsAll(string(payload), []string{"\"ParticipantID\":\"p1\"", "\"Name\":\"Alice\""}) {
		t.Fatalf("unexpected payload: %s", string(payload))
	}
}

func containsAll(s string, subs []string) bool {
	for _, sub := range subs {
		if !contains(s, sub) {
			return false
		}
	}
	return true
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (len(sub) == 0 || (func() bool { return indexOf(s, sub) >= 0 })())
}

func indexOf(s, sub string) int {
	// Simple substring search for small payloads to avoid importing strings in this file
	n, m := len(s), len(sub)
	if m == 0 {
		return 0
	}
	if m > n {
		return -1
	}
	for i := 0; i <= n-m; i++ {
		if s[i:i+m] == sub {
			return i
		}
	}
	return -1
}
