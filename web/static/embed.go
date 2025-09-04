package static

import "embed"

// FS embeds static assets like CSS.
//
//go:embed app.css app.js
var FS embed.FS
