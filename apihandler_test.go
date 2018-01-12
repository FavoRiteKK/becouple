package main_test

import (
	"becouple/appvendor"
	"becouple/models"
	"becouple/models/xodb"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/volatiletech/authboss"
	"golang.org/x/crypto/bcrypt"
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
	pass := "qweasd"
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
	// spew.Dump(obj)

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
	vals.Set(appvendor.PropPrimaryID, "qwe@gmail.com")
	vals.Set(appvendor.PropPassword, "qweasd")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)

	// to retrieve jwt token
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse, compare above")
	}

	//if result.Data == "" {
	//	t.Error("It should return jwt token, but got ", result.Err)
	//}

	// process confirm function test
	vals = url.Values{}
	vals.Set(appvendor.PropPrimaryID, "qwe@gmail.com")
	vals.Set(appvendor.PropPassword, "qweasd")
	vals.Set(appvendor.PropConfirmToken, "VOEMEJ")

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/confirm", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	//r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Jwt))

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

	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse.")
	}

	if w.Code != http.StatusOK {
		t.Error("It should be http 200:", w.Code)
	}

	//if result.Success == true {
	//	t.Error("It should be failed (the password parameter is wrong):", result.Success)
	//}

	if result.Data[appvendor.JFieldToken] != "" {
		t.Error("It should be empty jwt:", result.Data[appvendor.JFieldToken])
	}

	//if result.Err != bcrypt.ErrMismatchedHashAndPassword.Error() {
	//	t.Error("It encodes wrong error:", result.Err)
	//}

}

func TestApiAuthenticateSuccess(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	vals.Set(appvendor.PropPrimaryID, "qwe@gmail.com")
	vals.Set(appvendor.PropPassword, "qweasd")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	t.Logf("/auth response: %v, %v", w.Code, w.Body.String())

	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse.")
	}

	if w.Code != http.StatusOK {
		t.Error("It should be http 200:", w.Code)
	}

	if result.Success == false {
		t.Error("error: success is false, expect true")
	}

	if result.Data[appvendor.JFieldToken] == nil ||
		result.Data[appvendor.JFieldToken] == "" {
		t.Error("error: token field is empty, expect JWT token")
	}

}

func TestApiGetProfile(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// login user first
	vals.Set(appvendor.PropPrimaryID, "qwe@gmail.com")
	vals.Set(appvendor.PropPassword, "qweasd")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	t.Logf("/auth response: %v, %v", w.Code, w.Body.String())

	// to retrieve jwt token
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse, compare above")
	}

	// process getProfile function test
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("GET", "/api/user/profile", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Data[appvendor.JFieldToken]))

	h.ServeHTTP(w, r)
	t.Logf("/user/getprofile response: %v, %v", w.Code, w.Body.String())

	// test result
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}

	if result.Success != true {
		t.Error("/user/profile function seems malfunctioned.")
	}

	if result.Data[appvendor.JFieldUserProfile] == nil {
		t.Errorf("/user/profile function must return a '%v' block.", appvendor.JFieldToken)
	}
}

func TestApiEditPersonalInfo(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()
	vals := url.Values{}

	// login user first
	vals.Set(appvendor.PropPrimaryID, "qwe@gmail.com")
	vals.Set(appvendor.PropPassword, "qweasd")

	r, _ := http.NewRequest("POST", "/api/auth", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	t.Logf("/auth response: %v, %v", w.Code, w.Body.String())

	// to retrieve jwt token
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse, compare above")
	}

	// process edit-personal-info function test
	vals = url.Values{}
	vals.Set(appvendor.PropShortAbout, "I'm unit testing")
	vals.Set(appvendor.PropLivingAt, "Unit test lives nowhere")
	vals.Set(appvendor.PropWorkingAt, "Unit test works nowhere")
	vals.Set(appvendor.PropHomeTown, "Unit test's hometown is nowhere")
	vals.Set(appvendor.PropStatus, "complicate")
	vals.Set(appvendor.PropWeight, "50")  // 50kg
	vals.Set(appvendor.PropHeight, "170") // 170cm

	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/user/personalInfo", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Data[appvendor.JFieldToken]))

	h.ServeHTTP(w, r)
	t.Logf("/user/editPersonalInfo response: %v, %v", w.Code, w.Body.String())

	// test result
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}

	if result.Success != true {
		t.Error("editPersonalInfo function seems malfunctioned.")
	}
}

func TestApiRefreshToken(t *testing.T) {
	h := app.SetupMiddleware()

	w := httptest.NewRecorder()

	// process refresh token function test
	vals := url.Values{}
	vals.Set(appvendor.JFieldRefreshToken, "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJJZCI6InF3ZUBnbWFpbC5jb20ifQ.RlX2mNc9TWHQ1Nhjm7Gr7sPx0fp_nqztiK2Lj1KOedk")
	vals.Set(appvendor.PropDeviceName, "Lenovo P1ma40")

	w = httptest.NewRecorder()
	r, _ := http.NewRequest("POST", "/api/refreshToken", bytes.NewBufferString(vals.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	h.ServeHTTP(w, r)
	t.Logf("/user/editPersonalInfo response: %v, %v", w.Code, w.Body.String())

	// test result
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}

	if result.Success != true {
		t.Error("editPersonalInfo function seems malfunctioned.")
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
	result := new(models.ServerResponse)
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("It should response with proper type ServerResponse, compare above")
	}

	//if result.Jwt == "" {
	//t.Error("It should return jwt token, but got ", result.Err)
	//}

	// process logout function test
	w = httptest.NewRecorder()
	r, _ = http.NewRequest("POST", "/api/logout", nil)
	r.Header.Set("Authorization", fmt.Sprintf("Bearer %s", result.Data[appvendor.JFieldToken]))

	h.ServeHTTP(w, r)
	t.Logf("/logout response: %v, %v", w.Code, w.Body.String())

	// test result
	if err := json.NewDecoder(w.Body).Decode(result); err != nil {
		t.Error("Body response is malformed, compare above.")
	}
}
