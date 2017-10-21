package main_test

import (
	"becouple/models"
	"bytes"
	"encoding/json"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestApiAuthenticateWrong(t *testing.T) {

	w := httptest.NewRecorder()
	vals := url.Values{}

	email := "zeratul@heroes.com"
	vals.Set("primaryID", email)
	vals.Set("password", "qweasd123") // wrong password

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.Router.ServeHTTP(w, r)

	result := new(models.AuthResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type AuthResponse.")
	}

	if w.Code != http.StatusOK {
		t.Error("It should be http 200:", w.Code)
	}

	if result.Success == true {
		t.Error("It should be failed (the password parameter is wrong):", result.Success)
	}

	if result.Jwt != "" {
		t.Error("It should be empty jwt:", result.Jwt)
	}

	if result.Err != bcrypt.ErrMismatchedHashAndPassword.Error() {
		t.Error("It encodes wrong error:", result.Err)
	}

}
