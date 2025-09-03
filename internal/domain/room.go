package domain

import (
    "errors"
    "fmt"
    "strings"
)

type RoomID string
type ParticipantID string

const MaxParticipants = 25

type Participant struct {
    ID   ParticipantID
    Name string
}

type Room struct {
    id           RoomID
    participants map[ParticipantID]Participant
    names        map[string]ParticipantID // lowercase name → ID
    votes        map[ParticipantID]string // current round votes
    state        roundState
    round        int // increments on each Reset
}

func NewRoom(id RoomID) *Room {
    return &Room{
        id:           id,
        participants: make(map[ParticipantID]Participant),
        names:        make(map[string]ParticipantID),
        votes:        make(map[ParticipantID]string),
        state:        stateVoting,
        round:        0,
    }
}

// Join adds a participant with the given ID and display name.
// Name must be unique within the room (case-insensitive) and non-empty after trim.
func (r *Room) Join(id ParticipantID, name string) error {
    trimmed := strings.TrimSpace(name)
    if trimmed == "" {
        return errors.New("invalid name: empty")
    }
    if len(r.participants) >= MaxParticipants {
        return fmt.Errorf("capacity reached: max %d participants", MaxParticipants)
    }
    key := strings.ToLower(trimmed)
    if _, exists := r.names[key]; exists {
        return fmt.Errorf("duplicate name: %q", trimmed)
    }
    r.participants[id] = Participant{ID: id, Name: trimmed}
    r.names[key] = id
    return nil
}

type roundState int

const (
    stateVoting roundState = iota
    stateRevealed
)

var allowedCards = map[string]struct{}{
    "0": {}, "1": {}, "2": {}, "3": {}, "5": {}, "8": {}, "13": {}, "21": {}, "34": {},
    "?": {}, "∞": {}, "☕": {}, "Pass": {},
}

// CastVote records a participant's vote while in Voting state.
func (r *Room) CastVote(id ParticipantID, card string) error {
    if r.state != stateVoting {
        return errors.New("voting is closed")
    }
    if _, ok := r.participants[id]; !ok {
        return errors.New("not a participant")
    }
    if _, ok := allowedCards[card]; !ok {
        return fmt.Errorf("invalid card: %s", card)
    }
    r.votes[id] = card
    return nil
}

// Reveal reveals votes; requires at least one vote to exist; transitions to Revealed.
func (r *Room) Reveal() error {
    if r.state != stateVoting {
        return nil // idempotent
    }
    if len(r.votes) == 0 {
        return errors.New("cannot reveal: no votes")
    }
    r.state = stateRevealed
    return nil
}

// Reset starts a new round: clears all votes and re-opens voting.
func (r *Room) Reset() error {
    r.votes = make(map[ParticipantID]string)
    r.state = stateVoting
    r.round++
    return nil
}

// RoundIndex returns the current round index, starting at 0.
func (r *Room) RoundIndex() int { return r.round }

// ClearVote clears a participant's current vote in Voting state.
func (r *Room) ClearVote(id ParticipantID) error {
    if r.state != stateVoting {
        return errors.New("voting is closed")
    }
    if _, ok := r.participants[id]; !ok {
        return errors.New("not a participant")
    }
    delete(r.votes, id)
    return nil
}

// Leave removes a participant from the room and clears any vote.
func (r *Room) Leave(id ParticipantID) error {
    p, ok := r.participants[id]
    if !ok {
        return errors.New("not a participant")
    }
    // Remove vote if present
    delete(r.votes, id)
    // Remove name index and participant record
    key := strings.ToLower(strings.TrimSpace(p.Name))
    delete(r.names, key)
    delete(r.participants, id)
    return nil
}
