package main

import (
	"becouple/appvendor"
	"becouple/models"
	"becouple/models/xodb"
	"encoding/json"
	"fmt"
	"github.com/aarondl/tpl"
	jwtPkg "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/justinas/nosurf"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//=============================================================
// type WebController
//=============================================================

type WebController struct {
	app        *BeCoupleApp
	decoder    *schema.Decoder
	templates  tpl.Templates
	CsrfEnable bool
}

func NewWebController(app *BeCoupleApp) *WebController {
	ctrl := new(WebController)

	ctrl.app = app
	ctrl.decoder = schema.NewDecoder()
	ctrl.decoder.IgnoreUnknownKeys(true)

	ctrl.templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))

	ctrl.CsrfEnable = true

	return ctrl
}

// route '/', '/blogs'
func (ctrl *WebController) index(w http.ResponseWriter, r *http.Request) {
	data := ctrl.layoutData(w, r).MergeKV("posts", blogs)
	ctrl.mustRender(w, r, "index", data)
}

// route '/blogs/new
func (ctrl *WebController) newblog(w http.ResponseWriter, r *http.Request) {
	data := ctrl.layoutData(w, r).MergeKV("post", Blog{})
	ctrl.mustRender(w, r, "new", data)
}

var nextID = len(blogs) + 1

// route /blogs/new
func (ctrl *WebController) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	var b Blog
	if badRequest(w, ctrl.decoder.Decode(&b, r.PostForm)) {
		return
	}

	b.ID = nextID
	nextID++
	b.Date = time.Now()
	b.AuthorID = "Zeratul"

	blogs = append(blogs, b)

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/blogs/{id}/edit'
func (ctrl *WebController) edit(w http.ResponseWriter, r *http.Request) {
	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	data := ctrl.layoutData(w, r).MergeKV("post", blogs.Get(id))
	ctrl.mustRender(w, r, "edit", data)
}

// route '/blogs/{id}/edit'
func (ctrl *WebController) update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	var b = blogs.Get(id)
	if badRequest(w, ctrl.decoder.Decode(b, r.PostForm)) {
		return
	}

	b.Date = time.Now()

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/blogs/{id}/destroy'
func (ctrl *WebController) destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (ctrl *WebController) blogID(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	str := vars["id"]

	id, err := strconv.Atoi(str)
	if err != nil {
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	if id <= 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	return id, true
}

func (ctrl *WebController) mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	if ctrl.CsrfEnable {
		data.MergeKV("csrf_token", nosurf.Token(r))
	}

	err := ctrl.templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering templates:", err)
}

func (ctrl *WebController) layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ctrl.app.Ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*xodb.User).Fullname
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: ctrl.app.Ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ctrl.app.Ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

//=============================================================
// type ApiController
//=============================================================

type APIController struct {
	app       *BeCoupleApp
	validator map[string]func(r *http.Request) []error
}

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
			logrus.WithError(err).Errorln("validate email")
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
		var errs []error

		email := strings.TrimSpace(r.Header.Get(appvendor.PropEmail))

		if err := delegate.Var(email, "email"); err != nil {
			logrus.WithError(err).Errorln("validate email")
			errs = append(errs, err)
		}

		cnfToken := r.FormValue(appvendor.PropConfirmToken)

		if err := delegate.Var(cnfToken, "len=6"); err != nil {
			logrus.WithError(err).Errorln("validate confirm token")
			errs = append(errs, err)
		}

		return errs
	}

	return api
}

// route '/api/register'
func (api *APIController) register(w http.ResponseWriter, r *http.Request) {

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
	}

	w.Header().Set("Content-Type", "application/json")

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

// route '/api/confirm' params: 'Header JWT', 'confirm_token'
func (api *APIController) confirm(w http.ResponseWriter, r *http.Request) {

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
	}

	w.Header().Set("Content-Type", "application/json")

	// validate input
	if errs := api.validator["/confirm"](r); errs != nil {
		response.Err = appvendor.ConcateErrorWith(errs, "\n")
		json.NewEncoder(w).Encode(response)
		return
	}

	email := r.Header.Get(appvendor.PropEmail)
	cnfToken := r.FormValue(appvendor.PropConfirmToken)

	// get user from storer and check if exists
	obj, err := api.app.Storer.Get(email)
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

	// compare confirm token
	if user.Confirmed == true || user.ConfirmToken != strings.ToUpper(cnfToken) {
		response.ErrCode = appvendor.ErrorAccountCannotConfirm
		response.Err = "Cannot confirm the issuer"
		json.NewEncoder(w).Encode(response)
		return
	}

	// update user
	attr := authboss.Attributes{}
	attr[appvendor.PropConfirmToken] = ""
	attr[appvendor.PropConfirmed] = true

	if err := api.app.Storer.Put(email, attr); err != nil {
		appvendor.InternalServerError(w, "Cannot update attributes of confirmed user")
		return
	}

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

	w.Header().Set("Content-Type", "application/json")

	// validate input
	if errs := api.validator["/auth"](r); errs != nil {
		response.Err = appvendor.ConcateErrorWith(errs, "\n")
		json.NewEncoder(w).Encode(response)
		return
	}

	key := r.FormValue(appvendor.PropPrimaryID)
	password := r.FormValue(appvendor.PropPassword)

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

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	mcl := jwtPkg.MapClaims{
		"Id":  key,
		"exp": time.Now().Add(time.Minute * 10).Unix(),
		//TODO [production] may change these when go live
	}

	if user.Confirmed == false {
		mcl[appvendor.PropJwtError] = appvendor.ErrorAccountNotConfirmed
	}

	token := jwtPkg.NewWithClaims(appJwtSigningMethod, mcl)

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(appJwtSecret)) // must convert to []byte, otherwise we get error 'key is invalid'

	if err != nil {
		appvendor.InternalServerError(w, err.Error())
		return
	}

	response.Data[appvendor.JFieldToken] = tokenString

	// if user is not confirmed
	if user.Confirmed == false {
		response.ErrCode = appvendor.ErrorAccountNotConfirmed
		response.Err = "Account not confirmed"
		json.NewEncoder(w).Encode(response)
		return
	}

	response.Success = true
	response.ErrCode = 0

	json.NewEncoder(w).Encode(response)
}
