package httpadapter

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/jaminalder/estimations/internal/domain"
)

func (h *Handler) Cast(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := chi.URLParam(r, "roomID")
	if err := r.ParseForm(); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	card := strings.TrimSpace(r.FormValue("card"))
	pid := h.readPID(r)
	if pid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.svc.Cast(r.Context(), domain.RoomID(roomID), domain.ParticipantID(pid), card); err != nil {
		http.Error(w, "cast failed", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}

func (h *Handler) Clear(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := chi.URLParam(r, "roomID")
	pid := h.readPID(r)
	if pid == "" {
		http.Error(w, "unauthorized", http.StatusUnauthorized)
		return
	}
	if err := h.svc.Clear(r.Context(), domain.RoomID(roomID), domain.ParticipantID(pid)); err != nil {
		http.Error(w, "clear failed", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}

func (h *Handler) Reveal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := chi.URLParam(r, "roomID")
	if err := h.svc.Reveal(r.Context(), domain.RoomID(roomID)); err != nil {
		http.Error(w, "reveal failed", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}

func (h *Handler) Reset(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	roomID := chi.URLParam(r, "roomID")
	if err := h.svc.Reset(r.Context(), domain.RoomID(roomID)); err != nil {
		http.Error(w, "reset failed", http.StatusBadRequest)
		return
	}
	http.Redirect(w, r, "/rooms/"+roomID, http.StatusSeeOther)
}

func (h *Handler) readPID(r *http.Request) string {
	c, err := r.Cookie("pid")
	if err != nil || c == nil {
		return ""
	}
	return strings.TrimSpace(c.Value)
}
