package httpadapter

import (
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"regexp"
	"strings"
	"testing"

	"github.com/jaminalder/estimations/internal/adapters/idgen"
	"github.com/jaminalder/estimations/internal/adapters/memory"
	"github.com/jaminalder/estimations/internal/app"
)

func newTestServer(t *testing.T, logOut io.Writer) http.Handler {
	t.Helper()
	r, err := NewRenderer()
	if err != nil {
		t.Fatalf("renderer: %v", err)
	}
	repo := memory.NewRoomRepo()
	ids := idgen.NewRandom(10, 8)
	svc := &app.Service{Rooms: repo, Ids: ids}
	return NewServer(svc, r, WithLogger(log.New(logOut, "", 0)))
}

func TestLanding_ShowsCreateForm(t *testing.T) {
	var logs strings.Builder
	srv := newTestServer(t, &logs)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != 200 {
		t.Fatalf("status: got %d want 200", rec.Code)
	}
	body := rec.Body.String()
	if !strings.Contains(body, "<form") || !strings.Contains(body, "action=\"/rooms\"") || !strings.Contains(body, "name=\"title\"") {
		t.Fatalf("landing should contain form posting to /rooms with title field; got body: %q", body)
	}
}

func TestCreateRoom_RedirectsToLobby_AndLogs(t *testing.T) {
	var logs strings.Builder
	srv := newTestServer(t, &logs)

	// Post to create a room
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rooms", strings.NewReader("title=My+Session"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusSeeOther {
		t.Fatalf("status: got %d want 303", rec.Code)
	}
	loc := rec.Header().Get("Location")
	if loc == "" {
		t.Fatalf("expected Location header redirecting to lobby")
	}
	// Expect /rooms/{id}/lobby with id format from idgen (~16 chars, URL-safe)
	re := regexp.MustCompile(`^/rooms/[A-Za-z0-9_-]{11,20}/lobby$`)
	if !re.MatchString(loc) {
		t.Fatalf("unexpected redirect Location: %q", loc)
	}
	// Expect a log line with method, path, status
	if !strings.Contains(logs.String(), "POST /rooms 303") {
		t.Fatalf("logs should contain request line, got: %q", logs.String())
	}

	// Follow to lobby
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", loc, nil)
	srv.ServeHTTP(rec2, req2)
	if rec2.Code != 200 {
		t.Fatalf("lobby status: got %d want 200", rec2.Code)
	}
	body := rec2.Body.String()
	if !strings.Contains(body, "Your name") || !strings.Contains(body, "Enter Room") {
		t.Fatalf("lobby should contain name prompt, got: %q", body)
	}
}

func TestLobby_UnknownRoom_404_AndLogs(t *testing.T) {
	var logs strings.Builder
	srv := newTestServer(t, &logs)

	rec := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/rooms/unknown/lobby", nil)
	srv.ServeHTTP(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Fatalf("status: got %d want 404", rec.Code)
	}
	if !strings.Contains(logs.String(), "GET /rooms/unknown/lobby 404") {
		t.Fatalf("logs should contain 404 entry, got: %q", logs.String())
	}
}

func TestLobby_Join_RedirectsToRoom_AndLogs(t *testing.T) {
	var logs strings.Builder
	srv := newTestServer(t, &logs)

	// Create a room first
	rec := httptest.NewRecorder()
	req := httptest.NewRequest("POST", "/rooms", strings.NewReader("title=My+Session"))
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec, req)
	if rec.Code != http.StatusSeeOther {
		t.Fatalf("create status: %d", rec.Code)
	}
	loc := rec.Header().Get("Location")

	// Lobby should include form posting to /rooms/{id}/join
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest("GET", loc, nil)
	srv.ServeHTTP(rec2, req2)
	if rec2.Code != 200 {
		t.Fatalf("lobby status: %d", rec2.Code)
	}
	if !regexp.MustCompile(`action="/rooms/[A-Za-z0-9_-]{11,20}/join"`).MatchString(rec2.Body.String()) {
		t.Fatalf("lobby should contain join form action, got: %q", rec2.Body.String())
	}

	// Post join
	rec3 := httptest.NewRecorder()
	joinURL := strings.Replace(loc, "/lobby", "/join", 1)
	req3 := httptest.NewRequest("POST", joinURL, strings.NewReader("name=Alice"))
	req3.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	srv.ServeHTTP(rec3, req3)
	if rec3.Code != http.StatusSeeOther {
		t.Fatalf("join status: got %d want 303", rec3.Code)
	}
	roomURL := rec3.Header().Get("Location")
	if !regexp.MustCompile(`^/rooms/[A-Za-z0-9_-]{11,20}$`).MatchString(roomURL) {
		t.Fatalf("unexpected room redirect: %q", roomURL)
	}
	if !strings.Contains(logs.String(), "POST "+joinURL+" 303") {
		// Allow either absolute path or exact match; fallback check for method and 303
		if !strings.Contains(logs.String(), "POST /rooms/") || !strings.Contains(logs.String(), " 303 ") {
			t.Fatalf("logs should include join POST, got: %q", logs.String())
		}
	}

	// Follow to room page
	rec4 := httptest.NewRecorder()
	req4 := httptest.NewRequest("GET", roomURL, nil)
	srv.ServeHTTP(rec4, req4)
	if rec4.Code != 200 {
		t.Fatalf("room status: %d", rec4.Code)
	}
	body := rec4.Body.String()
	if !strings.Contains(body, "Select Your Estimate") || !strings.Contains(body, "Alice") {
		t.Fatalf("room page content missing participant/deck, got: %q", body)
	}
}
