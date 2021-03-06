package main

import (
	"becouple/appvendor"
	"becouple/models/xodb"
	"encoding/base64"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	// "github.com/betacraft/yaag/middleware"
	// "github.com/betacraft/yaag/yaag"

	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/justinas/alice"
	"github.com/justinas/nosurf"
	"github.com/onrik/logrus/filename"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
	aboauth "github.com/volatiletech/authboss/oauth2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// BeCoupleApp the main application
type BeCoupleApp struct {
	WebCtrl *WebController
	APICtrl *APIController
	Router  *mux.Router
	Ab      *authboss.Authboss
	Storer  *appvendor.AuthStorer
}

// NewApp creates new main application
func NewApp(authbossRootURL string) *BeCoupleApp {
	app := &BeCoupleApp{}

	app.BeforeSetup()
	app.Storer = appvendor.NewAuthStorer()

	app.SetupControllers()
	app.SetupAuthBoss(authbossRootURL)
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

// BeforeSetup : entrance before app's components setup
func (app *BeCoupleApp) BeforeSetup() {
	// log filename and line number
	logrus.AddHook(filename.NewHook())
	xodb.XOLog = func(query string, params ...interface{}) {
		logrus.WithFields(logrus.Fields{
			"query":  query,
			"params": params,
		}).Infoln("XOLog")
	}
}

// SetupControllers setups web and api controllers
func (app *BeCoupleApp) SetupControllers() {
	web := NewWebController(app)
	app.WebCtrl = web

	api := NewAPIController(app)
	app.APICtrl = api
}

// SetupAuthBoss setups authboss module
func (app *BeCoupleApp) SetupAuthBoss(rootURL string) {
	ab := authboss.New()
	app.Ab = ab

	ab.Storer = app.Storer
	ab.OAuth2Storer = app.Storer
	ab.MountPath = "/auth"
	ab.ViewsPath = "ab_views"
	ab.RootURL = rootURL

	ab.LayoutDataMaker = app.WebCtrl.layoutData

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

	ab.EmailFrom = "vankhiem583@gmail.com"

	//TODO [production] change to SMTPMailer in
	ab.Mailer = authboss.LogMailer(os.Stdout)
	//ab.Mailer = authboss.SMTPMailer("smtp.gmail.com:587",
	//	smtp.PlainAuth("", ab.EmailFrom, smtpGMailPass, "smtp.gmail.com"))

	// TODO [production] may change these when go
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
			MinLength:       4,
			MaxLength:       16,
			AllowWhitespace: false,
		},
	}

	if err := ab.Init(); err != nil {
		log.Fatal(err)
	}
}

// SetupRouter setups routes and paths
func (app *BeCoupleApp) SetupRouter() {
	// Set up our router
	router := mux.NewRouter()
	app.Router = router

	webRouter := router.PathPrefix("/").Subrouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	// Web Routes
	webRouter.PathPrefix("/auth").Handler(app.Ab.NewRouter())

	webRouter.Handle("/blogs/new", authProtect(app.WebCtrl.newblog, app.Ab)).Methods("GET")
	webRouter.Handle("/blogs/{id}/edit", authProtect(app.WebCtrl.edit, app.Ab)).Methods("GET")
	webRouter.HandleFunc("/blogs", app.WebCtrl.index).Methods("GET")
	webRouter.HandleFunc("/", app.WebCtrl.index).Methods("GET")

	webRouter.Handle("/blogs/{id}/edit", authProtect(app.WebCtrl.update, app.Ab)).Methods("POST")
	webRouter.Handle("/blogs/new", authProtect(app.WebCtrl.create, app.Ab)).Methods("POST")

	// This should actually be a DELETE but I can't be bothered to make a proper
	// destroy link using javascript atm.
	webRouter.Handle("/blogs/{id}/destroy", authProtect(app.WebCtrl.destroy, app.Ab)).Methods("POST")

	webRouter.HandleFunc("/test", func(writer http.ResponseWriter, r *http.Request) {
		//log.Println(appvendor.DBHelper.GetUserByEmail("qwe@gmail.com"))
	}).Methods("GET")

	// Api Routes
	apiRouter.HandleFunc("/register",
		WrapAPIResponseHeader(app.APICtrl.register)).Methods("POST")

	apiRouter.HandleFunc("/confirm",
		WrapAPIResponseHeader(app.APICtrl.confirm)).Methods("POST")

	apiRouter.HandleFunc("/auth",
		WrapAPIResponseHeader(app.APICtrl.authenticate)).Methods("POST")

	apiRouter.HandleFunc("/user/profile",
		WrapAPIResponseHeader(app.APICtrl.getProfile)).Methods("GET")

	apiRouter.HandleFunc("/user/personalInfo",
		WrapAPIResponseHeader(app.APICtrl.editPersonalInfo)).Methods("POST")

	apiRouter.HandleFunc("/user/basicInfo",
		WrapAPIResponseHeader(app.APICtrl.editBasicInfo)).Methods("POST")

	apiRouter.HandleFunc("/upload",
		WrapAPIResponseHeader(app.APICtrl.upload)).Methods("POST")

	apiRouter.HandleFunc("/logout",
		WrapAPIResponseHeader(app.APICtrl.logout)).Methods("POST")

	apiRouter.HandleFunc("/refreshToken",
		WrapAPIResponseHeader(app.APICtrl.refreshToken)).Methods("POST")

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		io.WriteString(w, "No such resource exists")
	})
}

// SetupMiddleware setups middelwares
func (app *BeCoupleApp) SetupMiddleware() http.Handler {
	// Set up our middleware chain
	// also, remove csrf validator for any route path that contains /api/
	stack := alice.New(logger,
		noresourceMiddleware(app.Router),
		nosurfing("/api/"),
		jwtMiddleware(),
		app.Ab.ExpireMiddleware).Then(app.Router)

	// TODO [PRODUCTION] remove yaag
	//yaag.Init(&yaag.Config{On: true, DocTitle: "Gorilla Mux", DocPath: "design/doc/apidoc.html"})
	//return middleware.Handle(stack)
	return stack
}

//SetupClientStore setups client store
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
