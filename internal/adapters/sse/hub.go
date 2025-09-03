package sse

import (
    "context"
    "encoding/json"
    "fmt"
    "sync"

    "github.com/jaminalder/estimations/internal/app"
    "github.com/jaminalder/estimations/internal/domain"
)

// Hub is a simple in-memory SSE broadcaster implementing app.Broadcaster.
// It manages room-scoped subscribers that receive marshaled event payloads.
type Hub struct {
    mu       sync.RWMutex
    rooms    map[domain.RoomID]map[chan []byte]struct{}
    bufSize  int
    marshal  func(v any) ([]byte, error)
}

// NewHub creates a hub with the given per-subscriber channel buffer size.
func NewHub(bufSize int) *Hub {
    return &Hub{
        rooms:   make(map[domain.RoomID]map[chan []byte]struct{}),
        bufSize: bufSize,
        marshal: json.Marshal,
    }
}

var _ app.Broadcaster = (*Hub)(nil)

// Subscribe registers a new subscriber for a room and returns a receive-only
// channel and an unsubscribe function. The channel will receive JSON-encoded
// event payloads.
func (h *Hub) Subscribe(roomID domain.RoomID) (<-chan []byte, func()) {
    ch := make(chan []byte, h.bufSize)
    h.mu.Lock()
    subs, ok := h.rooms[roomID]
    if !ok {
        subs = make(map[chan []byte]struct{})
        h.rooms[roomID] = subs
    }
    subs[ch] = struct{}{}
    h.mu.Unlock()

    // Unsubscribe closes the channel and removes it from the set.
    unsubscribe := func() {
        h.mu.Lock()
        if subs, ok := h.rooms[roomID]; ok {
            if _, present := subs[ch]; present {
                delete(subs, ch)
                close(ch)
                if len(subs) == 0 {
                    delete(h.rooms, roomID)
                }
            }
        }
        h.mu.Unlock()
    }
    return ch, unsubscribe
}

// Broadcast marshals the event to JSON and fan-outs to all subscribers
// registered to the given roomID. It is best-effort: if a subscriber's
// channel buffer is full, the message is dropped for that subscriber to avoid
// blocking other recipients.
func (h *Hub) Broadcast(ctx context.Context, roomID domain.RoomID, event any) error {
    // Marshal outside the lock
    payload, err := h.marshal(event)
    if err != nil {
        return fmt.Errorf("marshal event: %w", err)
    }

    h.mu.RLock()
    subs := h.rooms[roomID]
    // Copy keys to avoid holding lock while sending
    var chans []chan []byte
    for ch := range subs {
        chans = append(chans, ch)
    }
    h.mu.RUnlock()

    for _, ch := range chans {
        select {
        case <-ctx.Done():
            return ctx.Err()
        case ch <- payload:
            // sent
        default:
            // drop if subscriber is slow
        }
    }
    return nil
}

