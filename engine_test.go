package vodka

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouteParams(t *testing.T) {
	app := DefaultRouter()

	app.GET("/users/:id", func(c *Context) {
		c.String(200, c.Param("id"))
	})

	req := httptest.NewRequest("GET", "/users/67", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", w.Code)
	}

	if w.Body.String() != "67" {
		t.Fatalf("expected 67, got %s", w.Body.String())
	}
}
