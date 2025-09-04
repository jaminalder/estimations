package httpadapter

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/jaminalder/estimations/internal/domain"
)

// Room renders the dynamic room page with participants and deck.
func (h *Handler) Room(w http.ResponseWriter, r *http.Request) {
	roomID := strings.TrimSpace(chi.URLParam(r, "roomID"))
	if roomID == "" {
		http.NotFound(w, r)
		return
	}
	if h.svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}
	room, ok, err := h.svc.Rooms.Get(r.Context(), domain.RoomID(roomID))
	if err != nil {
		http.Error(w, "server error", http.StatusInternalServerError)
		return
	}
	if !ok || room == nil {
		http.NotFound(w, r)
		return
	}

	votes := room.Votes()
	me := h.readPID(r)
	type participantVM struct {
		Name     string
		HasVoted bool
		Card     string
		IsYou    bool
	}
	pvs := make([]participantVM, 0, len(room.Participants()))
	for _, p := range room.Participants() {
		card, has := votes[p.ID]
		pvs = append(pvs, participantVM{Name: p.Name, HasVoted: has, Card: card, IsYou: string(p.ID) == me})
	}
	data := struct {
		RoomID       string
		Participants []participantVM
		Total        int
		Voted        int
		Deck         []string
		Revealed     bool
	}{
		RoomID:       roomID,
		Participants: pvs,
		Total:        len(room.Participants()),
		Voted:        len(votes),
		Deck:         room.Deck(),
		Revealed:     room.IsRevealed(),
	}
	_ = h.r.Render(w, "room", data)
}
