package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"

	_ "gopkg.in/authboss.v1/auth"
	_ "gopkg.in/authboss.v1/confirm"
	_ "gopkg.in/authboss.v1/lock"
	_ "gopkg.in/authboss.v1/recover"
	_ "gopkg.in/authboss.v1/register"
	_ "gopkg.in/authboss.v1/remember"
)

var (
	funcs = template.FuncMap{
		"formatDate": func(date time.Time) string {
			return date.Format("2006/01/02 03:04pm")
		},
		"yield": func() string { return "" },
	}
)

func main() {
	// set address
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}
	addr := "localhost:" + port

	// setup our app
	app := NewApp(addr)

	// debug, list routes
	//router := app.router
	//router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
	//	t, err := route.GetPathTemplate()
	//	if err != nil {
	//		return err
	//	}
	//	// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
	//	// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
	//	p, err := route.GetPathRegexp()
	//	if err != nil {
	//		return err
	//	}
	//	m, err := route.GetMethods()
	//	if err != nil {
	//		return err
	//	}
	//	fmt.Println(strings.Join(m, ","), t, p)
	//	return nil
	//})

	// Start the server
	log.Println(http.ListenAndServe(addr, app.SetupMiddleware()))
}

func badRequest(w http.ResponseWriter, err error) bool {
	if err == nil {
		return false
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusBadRequest)
	fmt.Fprintln(w, "Bad request:", err)

	return true
}
