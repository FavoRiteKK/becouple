package main

import (
	"becouple/appvendor"
	"encoding/base64"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/authboss.v1"
	aboauth "gopkg.in/authboss.v1/oauth2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

type BeCoupleApp struct {
	Ctrl   *WebController
	Router *mux.Router
	Ab     *authboss.Authboss
	Storer *appvendor.AuthStorer
}

func NewApp(authbossRootUrl string) *BeCoupleApp {
	app := &BeCoupleApp{}

	app.Storer = appvendor.NewAuthStorer()

	app.SetupController()
	app.SetupAuthBoss(authbossRootUrl)
	app.SetupRouter()
	app.SetupClientStore()

	return app
}

//func (app *BeCoupleApp) GetController() *WebController {
//    return app.ctrl
//}
//
//// return app's router (without middlewares)
//func (app *BeCoupleApp) GetRouter() *mux.Router {
//	return app.router
//}
//
//func (app *BeCoupleApp) GetAuthBoss() *authboss.Authboss {
//    return app.ab
//}
//
//func (app *BeCoupleApp) GetStorer() *appvendor.AuthStorer {
//    return app.storer
//}

func (app *BeCoupleApp) SetupController() {
	ctrl := NewController(app)
	app.Ctrl = ctrl
}

func (app *BeCoupleApp) SetupAuthBoss(rootUrl string) {
	ab := authboss.New()
	app.Ab = ab

	ab.Storer = app.Storer
	ab.OAuth2Storer = app.Storer
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views"
	ab.RootURL = rootUrl

	ab.LayoutDataMaker = app.Ctrl.layoutData

	ab.OAuth2Providers = map[string]authboss.OAuth2Provider{
		"google": authboss.OAuth2Provider{
			OAuth2Config: &oauth2.Config{
				ClientID:     `751571472928-qfal1af5cn6ipstg8tl56rm0cncst9lv.apps.googleusercontent.com`,
				ClientSecret: `n5KWzxPao29Z1EzcCGCFmjHS`,
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

	ab.CookieStoreMaker = appvendor.NewCookieStorer
	ab.SessionStoreMaker = appvendor.NewSessionStorer

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

func (app *BeCoupleApp) SetupRouter() {
	// Set up our router
	router := mux.NewRouter()
	app.Router = router

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Web Routes
	webRouter.PathPrefix("/auth").Handler(app.Ab.NewRouter())

	webRouter.Handle("/blogs/new", authProtect(app.Ctrl.newblog, app.Ab)).Methods("GET")
	webRouter.Handle("/blogs/{id}/edit", authProtect(app.Ctrl.edit, app.Ab)).Methods("GET")
	webRouter.HandleFunc("/blogs", app.Ctrl.index).Methods("GET")
	webRouter.HandleFunc("/", app.Ctrl.index).Methods("GET")

	webRouter.Handle("/blogs/{id}/edit", authProtect(app.Ctrl.update, app.Ab)).Methods("POST")
	webRouter.Handle("/blogs/new", authProtect(app.Ctrl.create, app.Ab)).Methods("POST")

	// This should actually be a DELETE but I can't be bothered to make a proper
	// destroy link using javascript atm.
	webRouter.Handle("/blogs/{id}/destroy", authProtect(app.Ctrl.destroy, app.Ab)).Methods("POST")

	webRouter.HandleFunc("/test", func(writer http.ResponseWriter, r *http.Request) {
		log.Println(appvendor.DBHelper.GetUserByEmail("qwe@gmail.com"))
	}).Methods("GET")

	// Api Routes
	apiRouter.HandleFunc("/auth", app.Ctrl.authenticate).Methods("POST")
	apiRouter.HandleFunc("/logout", func(writer http.ResponseWriter, r *http.Request) {
		fmt.Println("Inside /api/logout?")

	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "No such resource exists")
	})
}

func (app *BeCoupleApp) SetupMiddleware() http.Handler {
	// Set up our middleware chain
	// also, remove csrf validator for any route path that contains /api/
	stack := alice.New(logger,
		nosurfing("/api/"),
		jwtMiddleware(),
		app.Ab.ExpireMiddleware).Then(app.Router)
	return stack
}

func (app *BeCoupleApp) SetupClientStore() {
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
	appvendor.CookieStore = securecookie.New(cookieStoreKey, nil)
	appvendor.SessionStore = sessions.NewCookieStore(sessionStoreKey)
}