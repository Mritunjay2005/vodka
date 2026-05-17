package vodka

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAbortStopsChain(t *testing.T) {
	calls := []string{}

	// First handler aborts
	handler1 := func(c *Context) {
		calls = append(calls, "h1")
		c.Abort()
	}

	// This should NOT run
	handler2 := func(c *Context) {
		calls = append(calls, "h2")
	}

	c := &Context{handlers: []HandlerFunc{handler1, handler2}, index: -1}
	c.Next()

	if len(calls) != 1 || calls[0] != "h1" {
		t.Errorf("got %v, want [h1]", calls)
	}
}

func TestKeys(t *testing.T) {
	app := DefaultRouter()

	app.Use(func(c *Context) {
		c.Set("M1", "First")

		c.Next()
	})

	app.Use(func(c *Context) {
		c.Set("M2", "Second")

		c.Next()
	})

	app.Use(func(c *Context) {
		c.Set("M1", "Changed")

		c.Next()
	})

	app.GET("/test", func(c *Context) {
		m1, exists := c.Get("M1")
		if !exists {
			t.Fatal("Value for M1 does not exist")
		}

		if m1.(string) != "Changed" {
			t.Fatalf("expected Changed, got %v", m1)
		}

		m2, exists := c.Get("M2")
		if !exists {
			t.Fatal("Value for M2 does not exist")
		}

		if m2.(string) != "Second" {
			t.Fatalf("expected Second, got %v", m2)
		}
	})

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)
}

func TestBindJSON(t *testing.T) {
	type User struct {
		Username string `json:"username"`
		Age      int8   `json:"age"`
	}

	app := DefaultRouter()

	app.POST("/test", func(c *Context) {
		var user User

		c.BindJSON(&user)

		if user.Username != "blufftunic" {
			t.Errorf("got %s, expected blufftunic", user.Username)
		}

		if user.Age != 20 {
			t.Errorf("got %d, expected blufftunic", user.Age)
		}

		c.JSON(200, M{
			"message": "success",
		})
	})

	body, _ := json.Marshal(M{
		"username": "blufftunic",
		"age":      20,
	})

	req := httptest.NewRequest(http.MethodPost, "/test", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	app.ServeHTTP(w, req)

	if w.Code != 200 {
		t.Errorf("Wrong Status Code: got %d, expected %d", w.Code, 200)
	}

	var response M

	json.Unmarshal(w.Body.Bytes(), &response)

	if response["message"] != "success" {
		t.Errorf("got %s, expected success", response["message"])
	}
}
