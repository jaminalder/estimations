package httpadapter

import (
	"io/fs"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	static "github.com/jaminalder/estimations/web/static"
)

type (
	serverOpts struct{ logger *log.Logger }
	Option     func(*serverOpts)
)

// WithLogger configures a standard logger for request logging.
func WithLogger(l *log.Logger) Option { return func(o *serverOpts) { o.logger = l } }

// newRouter builds the chi router with routes and middleware.
func newRouter(h *Handler, opts ...Option) http.Handler {
	var cfg serverOpts
	for _, o := range opts {
		o(&cfg)
	}

	r := chi.NewRouter()
	// Basic recoverer; keep logs readable
	r.Use(chimw.Recoverer)
	if cfg.logger != nil {
		r.Use(Logging(cfg.logger))
	}

	// Static assets
	if sub, err := fs.Sub(static.FS, "."); err == nil {
		r.Handle("/static/*", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
	}

	// Landing
	r.Get("/", h.Landing)
	r.Get("/landing", h.Landing)
	r.Post("/rooms", h.CreateRoom)

	// Rooms
	r.Route("/rooms/{roomID}", func(r chi.Router) {
		r.Get("/lobby", h.Lobby)
		r.Post("/join", h.Join)
		r.Get("/", h.Room)
	})

	// Fallbacks for legacy mockup routes
	r.Get("/lobby", h.Lobby) // expects roomID in data; will 404 without one
	r.Get("/room", func(w http.ResponseWriter, r *http.Request) {
		// Render the room template with empty values for mock view
		data := map[string]any{
			"RoomID":       "",
			"Participants": []any{},
			"Total":        0,
			"Voted":        0,
			"Deck":         []string{},
			"Revealed":     false,
		}
		_ = h.r.Render(w, "room", data)
	})

	return r
}
