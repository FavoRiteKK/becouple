package main

import (
	"becouple/appvendor"
	"becouple/models"
	"fmt"
	"github.com/aarondl/tpl"
	jwtPkg "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/json"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/justinas/nosurf"
	"gopkg.in/authboss.v1"
	"log"
	"net/http"
	"strconv"
	"time"
)

type AppController struct {
	app        *BeCoupleApp
	decoder    *schema.Decoder
	templates  tpl.Templates
	CsrfEnable bool
}

func NewController(app *BeCoupleApp) *AppController {
	ctrl := new(AppController)

	ctrl.app = app
	ctrl.decoder = schema.NewDecoder()
	ctrl.decoder.IgnoreUnknownKeys(true)

	ctrl.templates = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))

	ctrl.CsrfEnable = true

	return ctrl
}

// route '/', '/blogs'
func (ctrl *AppController) index(w http.ResponseWriter, r *http.Request) {
	data := ctrl.layoutData(w, r).MergeKV("posts", blogs)
	ctrl.mustRender(w, r, "index", data)
}

// route '/blogs/new
func (ctrl *AppController) newblog(w http.ResponseWriter, r *http.Request) {
	data := ctrl.layoutData(w, r).MergeKV("post", Blog{})
	ctrl.mustRender(w, r, "new", data)
}

var nextID = len(blogs) + 1

// route /blogs/new
func (ctrl *AppController) create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	// TODO: Validation

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
func (ctrl *AppController) edit(w http.ResponseWriter, r *http.Request) {
	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	data := ctrl.layoutData(w, r).MergeKV("post", blogs.Get(id))
	ctrl.mustRender(w, r, "edit", data)
}

// route '/blogs/{id}/edit'
func (ctrl *AppController) update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	// TODO: Validation

	var b = blogs.Get(id)
	if badRequest(w, ctrl.decoder.Decode(b, r.PostForm)) {
		return
	}

	b.Date = time.Now()

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/blogs/{id}/destroy'
func (ctrl *AppController) destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := ctrl.blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/api/auth'
func (ctrl *AppController) authenticate(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("primaryID")
	password := r.FormValue("password")

	log.Printf("pID: %v, pass: %v", key, password)

	// Create a new token object, specifying signing method and the claims
	// you would like it to contain.
	token := jwtPkg.NewWithClaims(appJwtSigningMethod, jwtPkg.MapClaims{
		"Id":  "Christopher",
		"exp": time.Now().Add(time.Hour * 1).Unix(),
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString([]byte(appJwtSecret)) // must convert to []byte, otherwise we get error 'key is invalid'

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ar := models.AuthResponse{
		Jwt:            tokenString,
		ServerResponse: &models.ServerResponse{Success: true},
	}

	jwtAuth, err := json.Marshal(ar)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jwtAuth)
}

func (ctrl *AppController) blogID(w http.ResponseWriter, r *http.Request) (int, bool) {
	vars := mux.Vars(r)
	str := vars["id"]

	id, err := strconv.Atoi(str)
	if err != nil {
		log.Println("Error parsing blog id:", err)
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	if id <= 0 {
		http.Redirect(w, r, "/", http.StatusFound)
		return 0, false
	}

	return id, true
}

func (ctrl *AppController) mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
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

func (ctrl *AppController) layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ctrl.app.ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*appvendor.AuthUser).Name
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: ctrl.app.ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ctrl.app.ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}
