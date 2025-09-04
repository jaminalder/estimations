package httpadapter

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/jaminalder/estimations/internal/domain"
)

// Lobby renders the lobby page for an existing room.
func (h *Handler) Lobby(w http.ResponseWriter, r *http.Request) {
	roomID := chi.URLParam(r, "roomID")
	if strings.TrimSpace(roomID) == "" {
		http.NotFound(w, r)
		return
	}
	if h.svc != nil {
		room, ok, err := h.svc.Rooms.Get(r.Context(), domain.RoomID(roomID))
		if err != nil {
			http.Error(w, "server error", http.StatusInternalServerError)
			return
		} else if !ok {
			http.NotFound(w, r)
			return
		}
		// If already joined (cookie pid present and matches a participant), redirect to room.
		if pid := h.readPID(r); pid != "" {
			for _, p := range room.Participants() {
				if string(p.ID) == pid {
					http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
					return
				}
			}
		}
	}
	data := struct{ RoomID string }{RoomID: roomID}
	_ = h.r.Render(w, "lobby", data)
}

// Join handles POST join and redirects to the room page.
func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := chi.URLParam(r, "roomID")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	name := strings.TrimSpace(r.FormValue("name"))
	if name == "" {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	if h.svc == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}
	// If already has a participant cookie for this room, just go to the room.
	if pid := h.readPID(r); pid != "" {
		if room, ok, _ := h.svc.Rooms.Get(r.Context(), domain.RoomID(roomID)); ok {
			for _, p := range room.Participants() {
				if string(p.ID) == pid {
					http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
					return
				}
			}
		}
	}
	pid, err := h.svc.Join(r.Context(), domain.RoomID(roomID), name)
	if err != nil {
		http.Error(w, "join failed", http.StatusBadRequest)
		return
	}
	// Scope participant cookie to this room path so multiple rooms don't collide.
	http.SetCookie(w, &http.Cookie{
		Name:     "pid",
		Value:    string(pid),
		Path:     "/rooms/" + roomID,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}
