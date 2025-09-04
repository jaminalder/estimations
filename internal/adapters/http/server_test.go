package httpadapter

import (
    "net/http/httptest"
    "strings"
    "testing"
)

func TestRoot_RendersLanding(t *testing.T) {
    var logs strings.Builder
    srv := newTestServer(t, &logs)

    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    srv.ServeHTTP(rec, req)

    if rec.Code != 200 {
        t.Fatalf("status: got %d want 200", rec.Code)
    }
    body := rec.Body.String()
    if body == "" || !contains(body, "Create Room") || !contains(body, "Session Title") {
        t.Fatalf("body should contain landing elements, got: %q", body)
    }
}

func TestLanding_RendersCreateRoom(t *testing.T) {
    var logs strings.Builder
    srv := newTestServer(t, &logs)

    req := httptest.NewRequest("GET", "/landing", nil)
    rec := httptest.NewRecorder()
    srv.ServeHTTP(rec, req)

    if rec.Code != 200 {
        t.Fatalf("status: got %d want 200", rec.Code)
    }
    body := rec.Body.String()
    if body == "" || !contains(body, "Create Room") || !contains(body, "Session Title") {
        t.Fatalf("body should contain landing elements, got: %q", body)
    }
}

func TestRoom_RendersMockup(t *testing.T) {
    var logs strings.Builder
    srv := newTestServer(t, &logs)

    req := httptest.NewRequest("GET", "/room", nil)
    rec := httptest.NewRecorder()
    srv.ServeHTTP(rec, req)

    if rec.Code != 200 {
        t.Fatalf("status: got %d want 200", rec.Code)
    }
    body := rec.Body.String()
    // Check key fragments from mockup
    mustContain := []string{
        "Voting in Progress",
        "Reveal Cards",
        "Reset Votes",
        "Select Your Estimate",
    }
    for _, sub := range mustContain {
        if !contains(body, sub) {
            t.Fatalf("body should contain %q, got: %q", sub, body)
        }
    }
}

func contains(s, sub string) bool { return len(s) >= len(sub) && (len(sub) == 0 || (func() bool { return indexOf(s, sub) >= 0 })()) }

func indexOf(s, sub string) int {
    n, m := len(s), len(sub)
    if m == 0 { return 0 }
    if m > n { return -1 }
    for i := 0; i <= n-m; i++ {
        if s[i:i+m] == sub { return i }
    }
    return -1
}
