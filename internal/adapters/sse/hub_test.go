package sse

import (
	"context"
	"testing"

	"github.com/jaminalder/estimations/internal/domain"
)

func TestHub_SubscribeBroadcastUnsubscribe(t *testing.T) {
	h := NewHub(4)
	room := domain.RoomID("r1")

	c1, u1 := h.Subscribe(room)
	c2, u2 := h.Subscribe(room)

	t.Cleanup(func() { u1(); u2() })

	// Broadcast an event
	if err := h.Broadcast(context.Background(), room, struct{ A string }{A: "x"}); err != nil {
		t.Fatalf("broadcast: %v", err)
	}

	// Both should receive something
	got1 := <-c1
	got2 := <-c2
	if len(got1) == 0 || len(got2) == 0 {
		t.Fatalf("expected payloads on both subscribers")
	}

	// Unsubscribe second and broadcast again
	u2()
	if err := h.Broadcast(context.Background(), room, struct{ B int }{B: 2}); err != nil {
		t.Fatalf("broadcast2: %v", err)
	}
	<-c1 // c1 still receives
}
