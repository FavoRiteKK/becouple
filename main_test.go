package main_test

import (
	. "becouple"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var app *BeCoupleApp

func setup() {
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}
	addr := "localhost:" + port
	app = NewApp(addr)

	app.Ab.XSRFMaker = func(_ http.ResponseWriter, _ *http.Request) (token string) {
		return "unused"
	}

	// disable csrf while testing
	app.Ctrl.CsrfEnable = false
}

func shutdown() {

}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestGetIndex(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/", nil)

	app.Router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("It should have written a 200: ", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("It should have wrote a response.")
	}

	if str := w.Body.String(); !strings.Contains(str, "Blogs - Index") {
		t.Error("It should have rendered 'Blog - Index': ", str)
	}
}

func TestGetRegisterIndex(t *testing.T) {
	w := httptest.NewRecorder()
	r, _ := http.NewRequest("GET", "/auth/register", nil)

	app.Router.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Error("It should have written a 200: ", w.Code)
	}

	if w.Body.Len() == 0 {
		t.Error("It should have wrote a response.")
	}

	if str := w.Body.String(); !strings.Contains(str, "<form") {
		t.Error("It should have rendered register form: ", str)
	}
}
