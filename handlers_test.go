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
	"github.com/volatiletech/authboss"
	"becouple/appvendor"
)

func TestApiRegister(t *testing.T) {

	w := httptest.NewRecorder()
	vals := url.Values{}

	// test data: already exist
	email := "zeratul@heroes.com"
	pass := "qwe123"
	vals.Set("primaryID", email)
	vals.Set("password", pass)

	r, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.Router.ServeHTTP(w, r)

	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse:", result)
	}

	if w.Code != http.StatusOK {
		t.Error("It should be http 200:", w.Code)
	}

	if result.Success == true {
		t.Error("It should be failed (email exists):", result.Success)
	}

	if result.Err != authboss.ErrUserFound.Error() {
		t.Error("It should be error 'user found':", result.Err)
	}

	// test data: new account
	email = "qwe@gmail.com"
	vals.Set("primaryID", email)
	vals.Set("password", pass)

	r, _ = http.NewRequest("POST", "/api/register", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	app.Router.ServeHTTP(w, r)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse:", result)
	}

	if result.Success != true && result.ErrCode != 0 {
		t.Error("Should return success, but got ", result)
	}

	obj, _ := app.Storer.Get(email)
	user, _ := obj.(*appvendor.AuthUser)
	if user == nil {
		t.Error("The user should be saved properly in the store.")
	}

	if user != nil {
		if user.Confirmed != false {
			t.Error("The new user's confirmation should be false")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
			t.Error("The new user's password and input password not match")
		}
	}
}

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
