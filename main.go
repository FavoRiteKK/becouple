package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/authboss.v1"
	_ "gopkg.in/authboss.v1/auth"
	_ "gopkg.in/authboss.v1/confirm"
	_ "gopkg.in/authboss.v1/lock"
	aboauth "gopkg.in/authboss.v1/oauth2"
	_ "gopkg.in/authboss.v1/recover"
	_ "gopkg.in/authboss.v1/register"
	_ "gopkg.in/authboss.v1/remember"

	"github.com/aarondl/tpl"
	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"becouple/appvendor"
	//"net/smtp"
	"regexp"
	"strings"
)

var funcs = template.FuncMap{
	"formatDate": func(date time.Time) string {
		return date.Format("2006/01/02 03:04pm")
	},
	"yield": func() string { return "" },
}

var (
	ab            = authboss.New()
	database      = NewMemStorer()
	templates     = tpl.Must(tpl.Load("views", "views/partials", "layout.html.tpl", funcs))
	schemaDec     = schema.NewDecoder()
	//smtpGMailPass = "qweasd1234"
)

func setupAuthboss(addr string) {
	ab.Storer = database
	ab.OAuth2Storer = database
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views"
	ab.RootURL = addr

	ab.LayoutDataMaker = layoutData

	ab.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"google": authboss.OAuth2Provider{
			OAuth2Config: &oauth2.Config{
				ClientID:     ``,
				ClientSecret: ``,
				Scopes:       []string{`profile`, `email`},
				Endpoint:     google.Endpoint,
			},
			Callback: aboauth.Google,
		},
	}

	b, err := ioutil.ReadFile(filepath.Join("views", "layout.html.tpl"))
	if err != nil {
		panic(err)
	}
	ab.Layout = template.Must(template.New("layout").Funcs(funcs).Parse(string(b)))

	ab.XSRFName = "csrf_token"
	ab.XSRFMaker = func(_ http.ResponseWriter, r *http.Request) string {
		return nosurf.Token(r)
	}

	ab.CookieStoreMaker = NewCookieStorer
	ab.SessionStoreMaker = NewSessionStorer

	ab.EmailFrom = "khiemnv@rikkeisoft.com"

	//TODO change to SMTPMailer in production
	ab.Mailer = authboss.LogMailer(os.Stdout)
	//ab.Mailer = authboss.SMTPMailer("smtp.gmail.com:587",
	//	smtp.PlainAuth("", ab.EmailFrom, smtpGMailPass, "smtp.gmail.com"))

	// TODO may change these when go production
	ab.Policies = []authboss.Validator{
		authboss.Rules{
			FieldName:       "email",
			Required:        true,
			MustMatch:       regexp.MustCompile(`^\S+@\S+$`),
			MatchError:      "Not an email address",
			AllowWhitespace: false,
		},
		authboss.Rules{
			FieldName:       "password",
			Required:        true,
			MinLength:       8,
			MaxLength:       16,
			AllowWhitespace: false,
		},
	}

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

func setupRouter() *mux.Router {
	// Set up our router
	schemaDec.IgnoreUnknownKeys(true)
	router := mux.NewRouter()
	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Web Routes
	webRouter.PathPrefix("/auth").Handler(ab.NewRouter())

	webRouter.Handle("/blogs/new", authProtect(newblog)).Methods("GET")
	webRouter.Handle("/blogs/{id}/edit", authProtect(edit)).Methods("GET")
	webRouter.HandleFunc("/blogs", index).Methods("GET")
	webRouter.HandleFunc("/", index).Methods("GET")

	webRouter.Handle("/blogs/{id}/edit", authProtect(update)).Methods("POST")
	webRouter.Handle("/blogs/new", authProtect(create)).Methods("POST")

	// This should actually be a DELETE but I can't be bothered to make a proper
	// destroy link using javascript atm.
	webRouter.Handle("/blogs/{id}/destroy", authProtect(destroy)).Methods("POST")

	webRouter.HandleFunc("/test", func(writer http.ResponseWriter, r *http.Request) {
		appvendor.Mgr.GetAllUser()
	}).Methods("GET")

	// Api Routes
	apiRouter.HandleFunc("/auth", authenticate).Methods("POST")
	apiRouter.HandleFunc("/logout", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside /api/logout?")

	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "No such resource exists")
	})

	return router
}

func setupStore() {
	// Initialize Sessions and Cookies
	// Typically gorilla securecookie and sessions packages require
	// highly random secret keys that are not divulged to the public.
	//
	// TODO In this example we use keys generated one time (if these keys ever become
	// compromised the gorilla libraries allow for key rotation, see gorilla docs)
	// The keys are 64-bytes as recommended for HMAC keys as per the gorilla docs.
	//
	// These values MUST be changed for any new project as these keys are already "compromised"
	// as they're in the public domain, if you do not change these your application will have a fairly
	// wide-opened security hole. You can generate your own with the code below, or using whatever method
	// you prefer:
	//
	//    func main() {
	//        fmt.Println(base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(64)))
	//    }
	//

	// We store them in base64 in the example to make it easy if we wanted to move them later to
	// a configuration environment var or file.
	cookieStoreKey, _ := base64.StdEncoding.DecodeString(`2S+t+bu22ZxFbCW0eFtwYChptomzJrjSR82AI1t3hgpHgjWRFPCHcFELqJ/Au+WCvwauz2Vgf51cpgbwY5Jnsg==`)
	sessionStoreKey, _ := base64.StdEncoding.DecodeString(`Ab5CP07McjLvEQvjmhZUyu3j7Dj2dCxDinbac89YAZXXc8RO9s/Sh8QSZwLrW0St0WazbWjFTA8kHdjXG3LXOQ==`)
	cookieStore = securecookie.New(cookieStoreKey, nil)
	sessionStore = sessions.NewCookieStore(sessionStoreKey)
}

func main() {
	setupStore()

	// set address
	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8000"
	}

	addr := "localhost:" + port
	// Initialize ab.
	setupAuthboss("http://" + addr)

	router := setupRouter()

	// Set up our middleware chain
	// also, remove csrf validator for any route path that contains /api/
	stack := alice.New(logger,
		nosurfing("/api/"),
		jwtMiddleware(),
		ab.ExpireMiddleware).Then(router)

	// debug, list routes
	router.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		t, err := route.GetPathTemplate()
		if err != nil {
			return err
		}
		// p will contain regular expression is compatible with regular expression in Perl, Python, and other languages.
		// for instance the regular expression for path '/articles/{id}' will be '^/articles/(?P<v0>[^/]+)$'
		p, err := route.GetPathRegexp()
		if err != nil {
			return err
		}
		m, err := route.GetMethods()
		if err != nil {
			return err
		}
		fmt.Println(strings.Join(m, ","), t, p)
		return nil
	})

	// Start the server
	log.Println(http.ListenAndServe(addr, stack))
}

func layoutData(w http.ResponseWriter, r *http.Request) authboss.HTMLData {
	currentUserName := ""
	userInter, err := ab.CurrentUser(w, r)
	if userInter != nil && err == nil {
		currentUserName = userInter.(*User).Name
	}

	return authboss.HTMLData{
		"loggedin":               userInter != nil,
		"username":               "",
		authboss.FlashSuccessKey: ab.FlashSuccess(w, r),
		authboss.FlashErrorKey:   ab.FlashError(w, r),
		"current_user_name":      currentUserName,
	}
}

func mustRender(w http.ResponseWriter, r *http.Request, name string, data authboss.HTMLData) {
	data.MergeKV("csrf_token", nosurf.Token(r))
	err := templates.Render(w, name, data)
	if err == nil {
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "Error occurred rendering template:", err)
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
