package main

import (
	"becouple/models"
	jwtPkg "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"time"
)

// route '/', '/blogs'
func index(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("posts", blogs)
	mustRender(w, r, "index", data)
}

// route '/blogs/new
func newblog(w http.ResponseWriter, r *http.Request) {
	data := layoutData(w, r).MergeKV("post", Blog{})
	mustRender(w, r, "new", data)
}

var nextID = len(blogs) + 1

// route /blogs/new
func create(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	// TODO: Validation

	var b Blog
	if badRequest(w, beApp.schemaDec.Decode(&b, r.PostForm)) {
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
func edit(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	data := layoutData(w, r).MergeKV("post", blogs.Get(id))
	mustRender(w, r, "edit", data)
}

// route '/blogs/{id}/edit'
func update(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if badRequest(w, err) {
		return
	}

	id, ok := blogID(w, r)
	if !ok {
		return
	}

	// TODO: Validation

	var b = blogs.Get(id)
	if badRequest(w, beApp.schemaDec.Decode(b, r.PostForm)) {
		return
	}

	b.Date = time.Now()

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/blogs/{id}/destroy'
func destroy(w http.ResponseWriter, r *http.Request) {
	id, ok := blogID(w, r)
	if !ok {
		return
	}

	blogs.Delete(id)

	http.Redirect(w, r, "/", http.StatusFound)
}

// route '/api/auth'
func authenticate(w http.ResponseWriter, r *http.Request) {
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

func blogID(w http.ResponseWriter, r *http.Request) (int, bool) {
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
