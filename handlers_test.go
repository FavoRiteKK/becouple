package main_test

import (
	"becouple/appvendor"
	"becouple/models"
	"becouple/models/xodb"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/volatiletech/authboss"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestApiRegisterExist(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// test data: already exist
	email := "zeratul@heroes.com"
	pass := "qwe123"
	fullName := "notimportant"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, pass)
	vals.Set(appvendor.PropFullName, fullName)

	r, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)

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
}

func TestApiRegisterNew(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// test data: new account
	email := "qwe@gmail.com"
	pass := "qwe123"
	fullName := "test master"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, pass)
	vals.Set(appvendor.PropFullName, fullName)

	r, _ := http.NewRequest("POST", "/api/register", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("POST /api/register should response with proper type ServerResponse:", result)
	}

	if result.Success != true && result.ErrCode != 0 {
		t.Error("Should return success, but got ", result)
	}

	obj, _ := app.Storer.Get(email)
	spew.Dump(obj)

	user, _ := obj.(*xodb.User)
	if user == nil {
		t.Error("The user should be saved properly in the store.")
	}

	if user != nil {
		if user.Confirmed == true {
			t.Error("The new user's confirmation should be false")
		}

		if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(pass)); err != nil {
			t.Error("The new user's password and input password not match")
		}

		if user.Fullname != fullName {
			t.Error("The new user's full name and input full name not match")
		}
	}

	// clean up user
	//if user != nil {
	//	app.Storer.DeletePermanently(user)
	//}

}

func TestApiConfirmUser(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// login user first
	email := "qwe@gmail.com"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, "qwe123")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)

	// to retrieve jwt token
	result := new(models.AuthResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type AuthResponse, compare above")
	}

	if result.Jwt == "" {
		t.Error("It should return jwt token, but got ", result.Err)
	}

	// process confirm function test
	vals = url.Values{}
	vals.Set(appvendor.PropConfirmToken, "asdzxc")

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/confirm", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Jwt))

	h.ServeHTTP(w, r)
	t.Logf("/confirm response: %v %v", w.Code, w.Body.String())

	// test result
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}

	if result.Success != true {
		t.Error("Confirm function seems malfunctioned.")
	}
}

func TestApiAuthenticateWrong(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	email := "zeratul@heroes.com"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, "qweasd123") // wrong password

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)

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

func TestApiAuthenticateSuccess(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	email := "qwe@gmail.com"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, "qwe123") // wrong password

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)

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

func TestApiLogout(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// login user first
	email := "qwe@gmail.com"
	vals.Set(appvendor.PropPrimaryID, email)
	vals.Set(appvendor.PropPassword, "qwe123")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	t.Logf("/auth response: %v, %v", w.Code, w.Body.String())

	// to retrieve jwt token
	result := new(models.AuthResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type AuthResponse, compare above")
	}

	if result.Jwt == "" {
		t.Error("It should return jwt token, but got ", result.Err)
	}

	// process logout function test
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/logout", nil)
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Jwt))

	h.ServeHTTP(w, r)
	t.Logf("/logout response: %v, %v", w.Code, w.Body.String())

	// test result
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}

	if result.Success != true {
		t.Error("Logout function seems malfunctioned")
	}
}
