package templates

import "embed"

// FS embeds the HTML templates for server-side rendering.
//go:embed layouts/*.tmpl.html pages/*.tmpl.html
var FS embed.FS

