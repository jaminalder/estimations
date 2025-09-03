package httpadapter

import (
    "html/template"
    "net/http"

    templates "github.com/jaminalder/estimations/web/templates"
)

// Renderer holds parsed templates and renders pages.
type Renderer struct {
    tpl *template.Template
}

// NewRenderer parses templates from the given glob patterns.
func NewRenderer() (*Renderer, error) {
    // Parse from embedded FS for testability and single binary deploy.
    tpl, err := template.ParseFS(templates.FS,
        "layouts/*.tmpl.html",
        "pages/*.tmpl.html",
    )
    if err != nil {
        return nil, err
    }
    return &Renderer{tpl: tpl}, nil
}

// NewServer returns an http.Handler with minimal routes.
// Import app only at the adapter layer to access services when wiring routes later.
// Keeping signature ready avoids churn when adding endpoints.
// Accepts svc but this minimal index route doesn't use it yet.
func NewServer(_ any, r *Renderer) http.Handler {
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
        w.Header().Set("Content-Type", "text/html; charset=utf-8")
        // Render the index page which extends the base layout.
        if err := r.tpl.ExecuteTemplate(w, "index", nil); err != nil {
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return
        }
    })
    return mux
}
