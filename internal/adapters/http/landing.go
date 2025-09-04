package httpadapter

import (
    "net/http"
)

// Landing renders the landing page (index).
func (h *Handler) Landing(w http.ResponseWriter, r *http.Request) {
    _ = h.r.Render(w, "index", nil)
}

// CreateRoom handles POST /rooms and redirects to the lobby.
func (h *Handler) CreateRoom(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost { http.Error(w, "method not allowed", http.StatusMethodNotAllowed); return }
    if err := r.ParseForm(); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
    if h.svc == nil { http.Error(w, "service unavailable", http.StatusInternalServerError); return }
    id, err := h.svc.CreateRoom(r.Context())
    if err != nil { http.Error(w, "failed to create room", http.StatusInternalServerError); return }
    http.Redirect(w, r, "/rooms/"+string(id)+"/lobby", http.StatusSeeOther)
}
