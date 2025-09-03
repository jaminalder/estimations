package idgen

import (
    "crypto/rand"
    "encoding/base64"

    "github.com/jaminalder/estimations/internal/app"
    "github.com/jaminalder/estimations/internal/domain"
)

// Random generates URL-safe opaque identifiers using crypto/rand
// and base64 URL encoding without padding.
type Random struct {
    roomBytes int
    partBytes int
}

// NewRandom creates a Random id generator with byte lengths for room and participant ids.
// Recommended: roomBytes=10 (≈16 chars), partBytes=8 (≈11 chars).
func NewRandom(roomBytes, participantBytes int) *Random {
    return &Random{roomBytes: roomBytes, partBytes: participantBytes}
}

var _ app.IdGen = (*Random)(nil)

func (r *Random) NewRoomID() domain.RoomID {
    return domain.RoomID(randString(r.roomBytes))
}

func (r *Random) NewParticipantID() domain.ParticipantID {
    return domain.ParticipantID(randString(r.partBytes))
}

func randString(n int) string {
    if n <= 0 {
        n = 8
    }
    b := make([]byte, n)
    _, _ = rand.Read(b)
    // URL-safe, no padding; characters are [A-Za-z0-9-_]
    return base64.RawURLEncoding.EncodeToString(b)
}

