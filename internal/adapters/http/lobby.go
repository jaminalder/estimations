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
    if strings.TrimSpace(roomID) == "" { http.NotFound(w, r); return }
    if h.svc != nil {
        if _, ok, err := h.svc.Rooms.Get(r.Context(), domain.RoomID(roomID)); err != nil {
            http.Error(w, "server error", http.StatusInternalServerError); return
        } else if !ok {
            http.NotFound(w, r); return
        }
    }
    data := struct{ RoomID string }{RoomID: roomID}
    _ = h.r.Render(w, "lobby", data)
}

// Join handles POST join and redirects to the room page.
func (h *Handler) Join(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    roomID := chi.URLParam(r, "roomID")
    if err := r.ParseForm(); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
    name := strings.TrimSpace(r.FormValue("name"))
    if name == "" { http.Error(w, "bad request", http.StatusBadRequest); return }
    if h.svc == nil { http.Error(w, "service unavailable", http.StatusInternalServerError); return }
    if _, err := h.svc.Join(r.Context(), domain.RoomID(roomID), name); err != nil {
        http.Error(w, "join failed", http.StatusBadRequest); return
    }
    http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}
