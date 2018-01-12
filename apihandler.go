//=============================================================
// type ApiController
//=============================================================
package main

import (
	"becouple/appvendor"
	"becouple/models"
	"becouple/models/xodb"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwtPkg "github.com/dgrijalva/jwt-go"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
)

const (
	// expireIn in 10 minutes (in miliseconds unit)
	expireIn int64 = 10 * 60 * 1000 //TODO may change in production
	// expDuration in duration unit
	expDuration time.Duration = time.Duration(expireIn) * time.Millisecond
)

//APIController controller for API routes
type APIController struct {
	app       *BeCoupleApp
	validator map[string]func(r *http.Request) []error
}

//NewAPIController creates new API controller
func NewAPIController(app *BeCoupleApp) *APIController {
	api := new(APIController)
	api.app = app
	api.validator = make(map[string]func(r *http.Request) []error)

	delegate := validator.New()

	// make validators
	api.validator["/auth"] = func(r *http.Request) []error {
		var errs []error

		key := strings.TrimSpace(r.FormValue(appvendor.PropPrimaryID))
		password := r.FormValue(appvendor.PropPassword)

		if err := delegate.Var(key, "email"); err != nil {
			logrus.WithError(err).Errorln("validate primaryID")
			errs = append(errs, err)
		}

		//TODO [production] sync with web client
		if err := delegate.Var(password, "min=4,max=16"); err != nil {
			logrus.WithError(err).Errorln("validate password")
			errs = append(errs, err)
		}

		return errs
	}

	api.validator["/register"] = func(r *http.Request) []error {
		errs := api.validator["/auth"](r)

		fullname := strings.TrimSpace(r.FormValue(appvendor.PropFullName))

		if err := delegate.Var(fullname, "min=1,max=45"); err != nil {
			logrus.WithError(err).Errorln("validate full name")
			errs = append(errs, err)
		}

		return errs
	}

	api.validator["/confirm"] = func(r *http.Request) []error {
		errs := api.validator["/auth"](r)

		cnfToken := r.FormValue(appvendor.PropConfirmToken)

		if err := delegate.Var(cnfToken, "len=6"); err != nil {
			logrus.WithError(err).Errorln("validate confirm token")
			errs = append(errs, err)
		}

		return errs
	}

	return api
}

func generateAccessToken(key string) (string, error) {
	expireAt := time.Now().Add(expDuration).Unix()
	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	mcl := jwtPkg.MapClaims{
		"Id":  key,
		"exp": expireAt,
		//TODO [production] may change these when go live
	}

	token := jwtPkg.NewWithClaims(appJwtSigningMethod, mcl)

	// Sign and get the complete encoded token as a string using the secret
	return token.SignedString([]byte(appJwtSecret)) // must convert to []byte, otherwise we get error 'key is invalid'
}

func generateRefreshToken(key string) (string, error) {
	// Create a new never-expired refresh_token
	mcl := jwtPkg.MapClaims{
		"Id": key,
	}

	token := jwtPkg.NewWithClaims(appJwtSigningMethod, mcl)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(appJwtSecret)) // must convert to []byte, otherwise we get error 'key is invalid'
	return tokenString, err
}

// route '/api/register'
func (api *APIController) register(w http.ResponseWriter, r *http.Request) {

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
	}

	// validate input
	if errs := api.validator["/register"](r); len(errs) > 0 {
		response.Err = appvendor.ConcateErrorWith(errs, "\n")
		json.NewEncoder(w).Encode(response)
		return
	}

	key := r.FormValue(appvendor.PropPrimaryID)
	password := r.FormValue(appvendor.PropPassword)
	fullname := r.FormValue(appvendor.PropFullName)

	// get user from store and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil && err != authboss.ErrUserNotFound {
		// unknown error, prevent further register
		appvendor.InternalServerError(w, err.Error())
		return
	} else if obj != nil {
		// user already exists
		response.Err = authboss.ErrUserFound.Error()
		response.ErrCode = appvendor.ErrorAccountAlreadyInUsed
		json.NewEncoder(w).Encode(response)
		return
	}

	// process registration
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	attr, err := authboss.AttributesFromRequest(r) // Attributes from overriden forms
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	attr[appvendor.PropEmail] = key
	attr[appvendor.PropPassword] = string(hashedPass)
	attr[appvendor.PropFullName] = fullname

	// insert user into store
	if err := api.app.Storer.Create(key, attr); err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	// fire authbosss afterEvent (to update user after registration step, ex: generate confirm token)
	ctx := api.app.Ab.InitContext(w, r)
	ctx.User = attr
	if err := api.app.Ab.Callbacks.FireAfter(authboss.EventRegister, ctx); err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	// user register successful
	response.Success = true
	response.ErrCode = 0
	json.NewEncoder(w).Encode(response)
	return
}

// route '/api/confirm' params: 'email', 'password', 'confirm_token'
func (api *APIController) confirm(w http.ResponseWriter, r *http.Request) {

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
		Data:    make(models.Data),
	}

	// validate input
	if errs := api.validator["/confirm"](r); errs != nil {
		response.Err = appvendor.ConcateErrorWith(errs, "\n")
		json.NewEncoder(w).Encode(response)
		return
	}

	key := r.FormValue(appvendor.PropPrimaryID)
	password := r.FormValue(appvendor.PropPassword)
	cnfToken := r.FormValue(appvendor.PropConfirmToken)
	deviceName := r.FormValue(appvendor.PropDeviceName)

	// get user from storer and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil {
		response.Err = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	user, ok := obj.(*xodb.User)
	if !ok {
		appvendor.InternalServerError(w, "Storer should returns a type AuthUser")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.ErrCode = appvendor.ErrorAccountAuthorizedFailed
		response.Err = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// compare confirm token
	if user.Confirmed == true || user.ConfirmToken != strings.ToUpper(cnfToken) {
		response.ErrCode = appvendor.ErrorAccountCannotConfirm
		response.Err = "Cannot confirm the issuer or wrong confirm token"
		json.NewEncoder(w).Encode(response)
		return
	}

	// update user
	attr := authboss.Attributes{}
	attr[appvendor.PropConfirmToken] = ""
	attr[appvendor.PropConfirmed] = true

	if err := api.app.Storer.Put(key, attr); err != nil {
		appvendor.InternalServerError(w, "Cannot update attributes of confirmed user")
		return
	}

	// create access token
	tokenString, err := generateAccessToken(key)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	// create refresh token
	refreshToken, err := generateRefreshToken(key)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}
	// save refresh token
	api.app.Storer.SaveCredential(refreshToken, key, deviceName)

	response.Data[appvendor.JFieldToken] = tokenString
	response.Data[appvendor.JFieldRefreshToken] = refreshToken
	response.Data[appvendor.JFieldExpireIn] = expireIn
	response.Success = true
	response.ErrCode = 0

	json.NewEncoder(w).Encode(response)
}

// route '/api/auth' params 'primaryID', 'password'
func (api *APIController) authenticate(w http.ResponseWriter, r *http.Request) {

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
		Data:    make(models.Data),
	}

	// validate input
	if errs := api.validator["/auth"](r); errs != nil {
		response.Err = appvendor.ConcateErrorWith(errs, "\n")
		json.NewEncoder(w).Encode(response)
		return
	}

	key := r.FormValue(appvendor.PropPrimaryID)
	password := r.FormValue(appvendor.PropPassword)
	deviceName := r.FormValue(appvendor.PropDeviceName)

	// get primary from storer and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil {
		response.Err = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	user, ok := obj.(*xodb.User)
	if !ok {
		appvendor.InternalServerError(w, "Storer should returns a type AuthUser")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.ErrCode = appvendor.ErrorAccountAuthorizedFailed
		response.Err = err.Error()
		json.NewEncoder(w).Encode(response)
		return
	}

	// if user is still being locked
	if user.Locked.Valid && user.Locked.Time.After(time.Now().UTC()) {
		response.ErrCode = appvendor.ErrorAccountBeingLocked
		response.Err = "Account is still locked. Try login again later"
		json.NewEncoder(w).Encode(response)
		return
	}

	// if user is not confirmed
	if user.Confirmed == false {
		response.ErrCode = appvendor.ErrorAccountNotConfirmed
		response.Err = "Account not confirmed"
		json.NewEncoder(w).Encode(response)
		return
	}

	// create access token
	tokenString, err := generateAccessToken(key)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	// create refresh token
	refreshToken, err := generateRefreshToken(key)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	// save refresh token
	api.app.Storer.SaveCredential(refreshToken, key, deviceName)

	response.Data[appvendor.JFieldToken] = tokenString
	response.Data[appvendor.JFieldRefreshToken] = refreshToken
	response.Data[appvendor.JFieldExpireIn] = expireIn
	response.Success = true
	response.ErrCode = 0

	json.NewEncoder(w).Encode(response)
}

// route '/api/user/profile' method: GET
func (api *APIController) getProfile(w http.ResponseWriter, r *http.Request) {
	key := r.Header.Get(appvendor.PropPrimaryID)

	// get user from storer and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil {
		appvendor.InternalServerError(w, "Cannot retrieve user from store")
		return
	}

	user, ok := obj.(*xodb.User)
	if !ok {
		appvendor.InternalServerError(w, "Storer should returns a type AuthUser")
		return
	}

	response := models.ServerResponse{
		Success: true,
		Data:    make(models.Data),
		ErrCode: 0,
	}

	response.Data[appvendor.JFieldUserProfile] = user

	json.NewEncoder(w).Encode(response)
}

// route '/api/user/personalInfo' params: shortAbout, livingAt, workingAt, hometown, status, weight, height
func (api *APIController) editPersonalInfo(w http.ResponseWriter, r *http.Request) {

	key := r.Header.Get(appvendor.PropPrimaryID)
	shortAbout := strings.TrimSpace(r.FormValue(appvendor.PropShortAbout))
	livingAt := strings.TrimSpace(r.FormValue(appvendor.PropLivingAt))
	workingAt := strings.TrimSpace(r.FormValue(appvendor.PropWorkingAt))
	hometown := strings.TrimSpace(r.FormValue(appvendor.PropHomeTown))
	status := r.FormValue(appvendor.PropStatus)
	weight := strings.TrimSpace(r.FormValue(appvendor.PropWeight))
	height := strings.TrimSpace(r.FormValue(appvendor.PropHeight))

	// get user from storer and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil {
		appvendor.InternalServerError(w, "Cannot retrieve user from store")
		return
	}

	_, ok := obj.(*xodb.User)
	if !ok {
		appvendor.InternalServerError(w, "Storer should returns a type AuthUser")
		return
	}

	// update user
	attr := authboss.Attributes{}
	attr[appvendor.PropShortAbout] = shortAbout
	attr[appvendor.PropLivingAt] = livingAt
	attr[appvendor.PropWorkingAt] = workingAt
	attr[appvendor.PropHomeTown] = hometown

	// correct status value
	attr[appvendor.PropStatus] = []byte(status) // must convert to []byte

	// correct weight value
	_weight, _ := strconv.ParseUint(weight, 10, 0)
	attr[appvendor.PropWeight] = uint8(_weight) // convert to match User's field type

	// correct height value
	_height, _ := strconv.ParseUint(height, 10, 0)
	attr[appvendor.PropHeight] = uint8(_height) // convert to match User's field type

	if err := api.app.Storer.Put(key, attr); err != nil {
		appvendor.InternalServerError(w, "Cannot update attributes of confirmed user")
		return
	}

	response := models.ServerResponse{
		Success: true,
		ErrCode: 0,
	}

	json.NewEncoder(w).Encode(response)
}

// route '/api/user/basicInfo' params: fullname, nickname, birthday, gender, job
func (api *APIController) editBasicInfo(w http.ResponseWriter, r *http.Request) {

	key := r.Header.Get(appvendor.PropPrimaryID)
	fullname := strings.TrimSpace(r.FormValue(appvendor.PropFullName))
	nickname := strings.TrimSpace(r.FormValue(appvendor.PropNickName))
	birthDay := strings.TrimSpace(r.FormValue(appvendor.PropDateOfBirth))
	gender := r.FormValue(appvendor.PropGender)
	job := strings.TrimSpace(r.FormValue(appvendor.PropJob))

	// get user from storer and check if exists
	obj, err := api.app.Storer.Get(key)
	if err != nil {
		appvendor.InternalServerError(w, "Cannot retrieve user from store")
		return
	}

	_, ok := obj.(*xodb.User)
	if !ok {
		appvendor.InternalServerError(w, "Storer should returns a type AuthUser")
		return
	}

	// update user
	attr := authboss.Attributes{}
	attr[appvendor.PropFullName] = fullname
	attr[appvendor.PropNickName] = nickname

	if birthDay != "" {
		attr[appvendor.PropDateOfBirth] = birthDay
	}

	attr[appvendor.PropJob] = job

	// correct status value
	attr[appvendor.PropGender] = []byte(gender) // must convert to []byte

	if err := api.app.Storer.Put(key, attr); err != nil {
		appvendor.InternalServerError(w, "Cannot update attributes of confirmed user")
		return
	}

	response := models.ServerResponse{
		Success: true,
		ErrCode: 0,
	}

	json.NewEncoder(w).Encode(response)
}

// route '/api/refreshToken' params: refresh_token, device_name
func (api *APIController) refreshToken(w http.ResponseWriter, r *http.Request) {
	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
		Data:    make(models.Data),
	}

	refreshToken := r.FormValue(appvendor.JFieldRefreshToken)
	deviceName := r.FormValue(appvendor.PropDeviceName)

	// find credential in db
	cred, err := api.app.Storer.GetCredentialByRefreshToken(refreshToken, deviceName)
	if err != nil {
		response.ErrCode = appvendor.ErrorRefreshTokenInvalid
		response.Err = "Invalid refresh token and/or device name"
		json.NewEncoder(w).Encode(response)
		return
	}

	// refresh access token
	tokenString, err := generateAccessToken(cred.Email)
	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	response.Data[appvendor.JFieldToken] = tokenString
	response.Data[appvendor.JFieldRefreshToken] = refreshToken
	response.Data[appvendor.JFieldExpireIn] = expireIn
	response.Success = true
	response.ErrCode = 0

	json.NewEncoder(w).Encode(response)

}

// route '/api/logout'
func (api *APIController) logout(w http.ResponseWriter, r *http.Request) {
	resp := models.ServerResponse{
		Success: true,
	}

	// if request is malformed
	if key := r.Header.Get(appvendor.PropEmail); key == "" {
		resp.Success = false
		resp.ErrCode = appvendor.ErrorGeneral
		resp.Err = "Request not contain proper key (extracted from jwt middleware."

		json.NewEncoder(w).Encode(resp)
		return
	}

	// request is fine
	json.NewEncoder(w).Encode(resp)
}

const (
	_uploadDir           string = "/home/khiemnv/Pictures/_goupload/"
	_maxBytesRequestBody int64  = 100 * 1024 * 1024
)

//http://sanatgersappa.blogspot.sg/2013/03/handling-multiple-file-uploads-in-go.html
// if got error 'http: multipart handled by ParseMultipartForm', disable yaag middle in app.go
// route '/api/upload'
func (api *APIController) upload(w http.ResponseWriter, r *http.Request) {
	// limit request body to 100MB
	r.Body = http.MaxBytesReader(w, r.Body, _maxBytesRequestBody)

	//get the multipart reader for the request.
	reader, err := r.MultipartReader()

	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	logrus.Infoln("Copy part to ", _uploadDir)
	//copy each part to destination.
	count := 0

	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}

		count++
		logrus.WithField("part #", count).Infoln("Next part counter now is")

		//if part.FileName() is empty, skip this iteration.
		if part.FileName() == "" {
			logrus.Debugln("caught empty filename")
			continue
		}

		logrus.WithField("filename", part.FileName()).Infoln("Create copy of uploaded file")
		dst, err := os.Create(_uploadDir + part.FileName())
		defer dst.Close()

		if err != nil {
			appvendor.InternalServerError(w, err.Error())
			return
		}

		logrus.Infoln("Copy from part to destination")
		if _, err := io.Copy(dst, part); err != nil {
			if err.Error() == "http: request body too large" {
				// request is too large
				json.NewEncoder(w).Encode(&models.ServerResponse{
					Success: false,
					ErrCode: appvendor.ErrorRequestBodyTooLarge,
					Err:     "http: request body too large",
				})
			} else {
				appvendor.InternalServerError(w, err.Error())
			}
			return
		}
	}
	//display success message
	logrus.Infoln("Upload process successes")
}
