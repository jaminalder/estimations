package httpadapter

import (
    "html/template"
    "io/fs"
    "log"
    "net/http"
    "path/filepath"
    "strings"
    "time"

    templates "github.com/jaminalder/estimations/web/templates"
    static "github.com/jaminalder/estimations/web/static"
    "github.com/jaminalder/estimations/internal/app"
    "github.com/jaminalder/estimations/internal/domain"
)

// Renderer holds parsed templates and renders pages.
type Renderer struct {
    pages map[string]*template.Template
}

// NewRenderer parses templates from the given glob patterns.
func NewRenderer() (*Renderer, error) {
    // Parse base layouts once
    base, err := template.ParseFS(templates.FS, "layouts/*.tmpl.html")
    if err != nil {
        return nil, err
    }
    // Build a per-page template set to avoid block name collisions
    matches, err := fs.Glob(templates.FS, "pages/*.tmpl.html")
    if err != nil {
        return nil, err
    }
    pages := make(map[string]*template.Template, len(matches))
    for _, m := range matches {
        // derive page name from filename (e.g., pages/index.tmpl.html -> index)
        _, file := filepath.Split(m)
        name := strings.SplitN(file, ".", 2)[0]
        clone, err := base.Clone()
        if err != nil {
            return nil, err
        }
        // Parse the single page into the clone
        if _, err := clone.ParseFS(templates.FS, m); err != nil {
            return nil, err
        }
        pages[name] = clone
    }
    return &Renderer{pages: pages}, nil
}

func (r *Renderer) Render(w http.ResponseWriter, page string, data any) error {
    w.Header().Set("Content-Type", "text/html; charset=utf-8")
    tpl, ok := r.pages[page]
    if !ok {
        http.Error(w, "template not found", http.StatusNotFound)
        return nil
    }
    return tpl.ExecuteTemplate(w, page, data)
}

type serverOpts struct{ logger *log.Logger }
type Option func(*serverOpts)

// WithLogger configures a standard logger for request logging.
func WithLogger(l *log.Logger) Option { return func(o *serverOpts) { o.logger = l } }

// NewServer wires HTTP routes. svc provides app use-cases.
func NewServer(svc *app.Service, r *Renderer, opts ...Option) http.Handler {
    mux := http.NewServeMux()

    // Static assets (embedded)
    if sub, err := fs.Sub(static.FS, "."); err == nil {
        mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(sub))))
    }

    // Landing page (static mockup)
    mux.HandleFunc("/landing", func(w http.ResponseWriter, req *http.Request) {
        if err := r.Render(w, "index", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    // Lobby page (default at /)
    mux.HandleFunc("/lobby", func(w http.ResponseWriter, req *http.Request) {
        if err := r.Render(w, "lobby", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    // Room mockup page
    mux.HandleFunc("/room", func(w http.ResponseWriter, req *http.Request) {
        if err := r.Render(w, "room", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    // Root loads landing
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        if err := r.Render(w, "index", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })

    // Create room (POST) -> redirect to lobby
    mux.HandleFunc("/rooms", func(w http.ResponseWriter, req *http.Request) {
        if req.Method != http.MethodPost {
            http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
            return
        }
        if err := req.ParseForm(); err != nil {
            http.Error(w, "bad request", http.StatusBadRequest)
            return
        }
        _ = strings.TrimSpace(req.FormValue("title")) // reserved for future use
        if svc == nil {
            http.Error(w, "service unavailable", http.StatusInternalServerError)
            return
        }
        id, err := svc.CreateRoom(req.Context())
        if err != nil {
            http.Error(w, "failed to create room", http.StatusInternalServerError)
            return
        }
        http.Redirect(w, req, "/rooms/"+string(id)+"/lobby", http.StatusSeeOther)
    })

    // /rooms/* routes: lobby, room, and join
    mux.HandleFunc("/rooms/", func(w http.ResponseWriter, req *http.Request) {
        path := strings.TrimPrefix(req.URL.Path, "/rooms/")
        // Handle join first for POST
        if req.Method == http.MethodPost && strings.HasSuffix(path, "/join") {
            id := strings.TrimSuffix(path, "/join")
            id = strings.TrimSuffix(id, "/")
            if id == "" { http.NotFound(w, req); return }
            if err := req.ParseForm(); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
            name := strings.TrimSpace(req.FormValue("name"))
            if name == "" { http.Error(w, "bad request", http.StatusBadRequest); return }
            if svc == nil { http.Error(w, "service unavailable", http.StatusInternalServerError); return }
            if _, err := svc.Join(req.Context(), domain.RoomID(id), name); err != nil {
                http.Error(w, "join failed", http.StatusBadRequest)
                return
            }
            http.Redirect(w, req, "/rooms/"+id, http.StatusSeeOther)
            return
        }
        // Only handle /rooms/{id}/lobby and /rooms/{id}
        switch {
        case strings.HasSuffix(path, "/lobby"):
            id := strings.TrimSuffix(path, "/lobby")
            id = strings.TrimSuffix(id, "/")
            if id == "" { http.NotFound(w, req); return }
            if svc != nil {
                _, ok, err := svc.Rooms.Get(req.Context(), domain.RoomID(id))
                if err != nil { http.Error(w, "server error", http.StatusInternalServerError); return }
                if !ok { http.NotFound(w, req); return }
            }
            data := struct{ RoomID string }{RoomID: id}
            if err := r.Render(w, "lobby", data); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            return
        case strings.Contains(path, "/") && !strings.Contains(path, "/"):
            // unreachable placeholder
        }
        // If path is exactly {id}
        if !strings.Contains(path, "/") {
            id := strings.TrimSuffix(path, "/")
            if id == "" { http.NotFound(w, req); return }
            if svc == nil { http.Error(w, "service unavailable", http.StatusInternalServerError); return }
            room, ok, err := svc.Rooms.Get(req.Context(), domain.RoomID(id))
            if err != nil { http.Error(w, "server error", http.StatusInternalServerError); return }
            if !ok || room == nil { http.NotFound(w, req); return }
            // Build minimal view model
            votes := room.Votes()
            type participantVM struct{
                Name string
                HasVoted bool
            }
            pvs := make([]participantVM, 0, len(room.Participants()))
            for _, p := range room.Participants() {
                _, has := votes[p.ID]
                pvs = append(pvs, participantVM{Name: p.Name, HasVoted: has})
            }
            data := struct{
                RoomID string
                Participants []participantVM
                Total int
                Voted int
                Deck []string
                Revealed bool
            }{
                RoomID: id,
                Participants: pvs,
                Total: len(room.Participants()),
                Voted: len(votes),
                Deck: room.Deck(),
                Revealed: room.IsRevealed(),
            }
            if err := r.Render(w, "room", data); err != nil {
                http.Error(w, err.Error(), http.StatusInternalServerError)
                return
            }
            return
        }
        http.NotFound(w, req)
    })

    // Wrap with logging if provided
    var handler http.Handler = mux
    var cfg serverOpts
    for _, o := range opts { o(&cfg) }
    if cfg.logger != nil {
        handler = loggingMiddleware(cfg.logger)(handler)
    }
    return handler
}

// loggingMiddleware logs method, path, status and duration.
func loggingMiddleware(l *log.Logger) func(http.Handler) http.Handler {
    type respWriter struct {
        http.ResponseWriter
        code int
    }
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            start := time.Now()
            rw := &respWriter{ResponseWriter: w, code: http.StatusOK}
            // Wrap WriteHeader to capture status
            ww := struct{ http.ResponseWriter }{ResponseWriter: rw}
            // Replace WriteHeader via embedding trick
            rw.ResponseWriter = &writeHeaderCapturer{ResponseWriter: w, code: &rw.code}
            next.ServeHTTP(rw.ResponseWriter, r)
            dur := time.Since(start)
            l.Printf("%s %s %d %s", r.Method, r.URL.Path, rw.code, dur)
            _ = ww // avoid unused in case of build nuances
        })
    }
}

type writeHeaderCapturer struct {
    http.ResponseWriter
    code *int
}

func (w *writeHeaderCapturer) WriteHeader(statusCode int) {
    *w.code = statusCode
    w.ResponseWriter.WriteHeader(statusCode)
}
