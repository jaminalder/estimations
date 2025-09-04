package httpadapter

import (
    "html/template"
    "io/fs"
    "net/http"
    "path/filepath"
    "strings"

    templates "github.com/jaminalder/estimations/web/templates"
    static "github.com/jaminalder/estimations/web/static"
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

// NewServer returns an http.Handler with minimal routes.
// Import app only at the adapter layer to access services when wiring routes later.
// Keeping signature ready avoids churn when adding endpoints.
// Accepts svc but this minimal index route doesn't use it yet.
func NewServer(_ any, r *Renderer) http.Handler {
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
    return mux
}
