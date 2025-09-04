package httpadapter

import (
    "html/template"
    "io/fs"
    "net/http"
    "path/filepath"
    "strings"

    templates "github.com/jaminalder/estimations/web/templates"
    "github.com/jaminalder/estimations/internal/app"
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

// Handler bundles dependencies for request handlers.
type Handler struct {
    svc *app.Service
    r   *Renderer
}

// NewServer wires routes using chi and returns an http.Handler.
func NewServer(svc *app.Service, r *Renderer, opts ...Option) http.Handler {
    h := &Handler{svc: svc, r: r}
    return newRouter(h, opts...)
}
