package httpadapter

import (
    "net/http/httptest"
    "testing"
)

func TestIndex_RendersHello(t *testing.T) {
    r, err := NewRenderer()
    if err != nil { t.Fatalf("renderer: %v", err) }
    srv := NewServer(nil, r)

    req := httptest.NewRequest("GET", "/", nil)
    rec := httptest.NewRecorder()
    srv.ServeHTTP(rec, req)

    if rec.Code != 200 {
        t.Fatalf("status: got %d want 200", rec.Code)
    }
    body := rec.Body.String()
    if body == "" || !contains(body, "Hello") {
        t.Fatalf("body should contain 'Hello', got: %q", body)
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
