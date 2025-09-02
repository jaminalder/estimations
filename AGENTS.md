# AGENTS.md — Estimation Poker (Go, SSR, htmx)

_Last updated: 2025-09-02_

## Project
Interactive estimation poker for Scrum teams.  
Stack: **Go 1.25**, chi, server-side rendered multipage app, htmx, Alpine.js, Bulma, **SSE**.  
No auth/login: users type a display name, create a room, share the URL.  
No persistence: **in-memory** only (ephemeral).  
Repo: https://github.com/jaminalder/estimations  
Module path: **github.com/jaminalder/estimations**  
Dev: `air` auto-reload.

## Workflow (must follow)
1) Plan: restate the feature + acceptance criteria.  
2) Think tests: pick layer (domain/app/adapter).  
3) Write the **first failing test**.  
4) Implement the **minimal** code to pass.  
5) Repeat in small steps.  
6) If requirements/tech choices are unclear → **ask**.

## Run / Dev / Test
- Dev (auto-reload): `make dev`  _(requires `air`, install: `go install github.com/air-verse/air@latest`)_
- Run once: `make run`
- Build: `make build`
- Test all: `make test`  (no caching)  
- Port: `:8080` (HTTP). SSE: `/events/{roomID}`.

## Repo map (high level)
- `cmd/server/` — main (wire router, deps, templates)
- `internal/domain/` — entities + domain services (pure Go)
- `internal/app/` — use-cases & ports (RoomRepo, IdGen, Clock, Broadcaster)
- `internal/adapters/http/` — chi routes, handlers, template rendering (SSR + htmx)
- `internal/adapters/sse/` — SSE hub, client registry, broadcasts
- `internal/adapters/memory/` — in-memory RoomRepo (no external deps)
- `web/templates/` — layouts, pages, partials (SSR; **CDNs** for CSS/JS)
- `docs/adr/` — lightweight ADRs (DDD-lite, SSR+htmx, SSE choice)

## Coding rules
- Domain stays **pure** (stdlib only).  
- Use **ports** in `internal/app/ports.go`; adapters implement them.  
- SSR first; **htmx** for partial swaps; **Alpine** for tiny local state.  
- **SSE** for realtime room events (join/leave/vote/reveal/reset).  
- Pass deps via constructors; no global state.  
- Errors: wrap with context; no panics in request path.

## Testing rules
- Domain: table-driven unit tests.  
- App/use-cases: unit tests with fakes.  
- HTTP: `httptest` for handlers (status + HTML fragment).  
- SSE: unit-test hub; one integration test from action → SSE event.  
- Each bug adds a failing test first.

## Security & safety
- No secrets or `.env` pushed to repo.  
- Escape/validate all input; **template auto-escaping ON**.  
- CSRF: ensure token on state-changing POSTs (htmx still posts forms).  
- Simple rate-limit on room creation (middleware).

## UI (CDN)
Include in `web/templates/layouts/base.tmpl.html`:
- Bulma: `https://cdn.jsdelivr.net/npm/bulma@latest/css/bulma.min.css`
- htmx: `https://unpkg.com/htmx.org@1.9.12`
- Alpine: `https://unpkg.com/alpinejs@3.x.x/dist/cdn.min.js`

## Current priorities (Now/Next/Later)
- **Now:** Design and implement the domain
- **Next:** Create/Join room, pick deck, cast vote, reveal/hide, reset (SSR + htmx).  
- **Next:** SSE broadcasts for join/leave/vote/reveal; room page updates live.  
- **Later:** Export results; persistence; ownership/admin.

## Key decisions (brief)
- DDD-lite + ports/adapters; **in-memory** repo for v1.  
- **SSR** pages, progressively enhanced with htmx; **Bulma** styling.  
- **SSE** over WebSockets (simplicity).

