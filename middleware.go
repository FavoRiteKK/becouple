package main

import (
	"fmt"
	"net/http"

	"becouple/appvendor"
	"becouple/models"
	"encoding/json"
	jwtPkg "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/justinas/nosurf"
	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
	"net/http/httptest"
	"regexp"
	"strings"
)

type authProtector struct {
	f  http.HandlerFunc
	ab *authboss.Authboss
}

func authProtect(f http.HandlerFunc, ab *authboss.Authboss) authProtector {
	return authProtector{f, ab}
}

func (ap authProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := ap.ab.CurrentUser(w, r); err != nil {
		logrus.WithError(err).Errorln("error fetching current user")
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		logrus.WithField("path", r.URL.Path).Errorln("Redirecting unauthorized user from path")
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		ap.f(w, r)
	}
}

func nosurfing(exemptedRegex interface{}) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		surfing := nosurf.New(h)

		if exemptedRegex != nil {
			surfing.ExemptRegexp(exemptedRegex)
		}

		surfing.SetFailureHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logrus.WithField("csrf_token", nosurf.Reason(r)).Errorln("Failed to validate XSRF Token")
			w.WriteHeader(http.StatusBadRequest)
		}))

		return surfing
	}
}

func logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		session, err := appvendor.SessionStore.Get(r, appvendor.SessionCookieName)
		if err == nil {
			fmt.Print("Session: ")
			first := true
			for k, v := range session.Values {
				if first {
					first = false
				} else {
					fmt.Print(", ")
				}
				fmt.Printf("%s = %v", k, v)
			}
			fmt.Println()
		}
		//fmt.Print("Database: ")
		//for _, u := range database.Users {
		//	fmt.Printf("%#v\n", u)
		//}

		// switch out response writer for a recorder
		// for all subsequent handlers
		rec := httptest.NewRecorder()
		next.ServeHTTP(rec, r)

		// log result, after request went through middleware chains
		logrus.WithFields(logrus.Fields{
			"code": rec.Code,
			"body": rec.Body.String(),
		}).Infoln("Http response")

		// copy everything from response recorder
		// to actual response writer
		for k, v := range rec.HeaderMap {
			w.Header()[k] = v
		}
		w.WriteHeader(rec.Code)
		rec.Body.WriteTo(w)

	})
}

// detect no such resource exists
func noresourceMiddleware(router *mux.Router) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cnf := NewNoResourceHandler()
		cnf.next = next
		cnf.router = router

		// if this middleware is used, then router's notfoundhandler must be nil
		router.NotFoundHandler = nil
		return cnf
	}
}

func jwtMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		jwtAuth := NewJwtAuth()

		// regex to check if route contains '/api/'
		// if it does, the middleware would check jwt token
		regex := regexp.MustCompile("/api/")

		jwtAuth.token = jwtPkg.New(appJwtSigningMethod)
		jwtAuth.next = next
		jwtAuth.except = func(r *http.Request) bool {
			path := r.URL.Path

			// exempt this exactly path
			if strings.Contains("/api/auth;/api/register;", path) {
				return true
			}

			return !regex.MatchString(path)
		}
		return jwtAuth
	}
}

func confirmingMiddleware() func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		cnf := NewConfirmingHandler()
		cnf.next = next

		cnf.except = func(r *http.Request) bool {
			path := r.URL.Path

			// exempt this exactly path
			if strings.Contains("/api/confirm", path) {
				return true
			}

			return false
		}
		return cnf
	}
}

//////////////////////////////////////////////////
// jwtAuth middleware
//////////////////////////////////////////////////
var (
	appJwtSigningMethod = jwtPkg.SigningMethodHS256
	appJwtSecret        = "qweasd123"
)

type JwtAuth struct {
	token *jwtPkg.Token
	next  http.Handler
	// method to parse request, return true if request should be skipped for jwt token validation
	except func(r *http.Request) bool
}

func NewJwtAuth() *JwtAuth {
	jwtToken := &JwtAuth{}

	return jwtToken
}

func (jwt *JwtAuth) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if except function return true for the request, then skip check jwt token
	if jwt.except(r) {
		jwt.next.ServeHTTP(w, r)
		return
	}

	// parse jwt token
	token, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwtPkg.Token) (interface{}, error) {
		b := []byte(appJwtSecret)
		return b, nil
	})

	if vErr, ok := err.(*jwtPkg.ValidationError); ok {
		if vErr.Errors&jwtPkg.ValidationErrorExpired != 0 {
			json.NewEncoder(w).Encode(models.ServerResponse{
				Success: false,
				ErrCode: appvendor.ErrorTokenExpired,
				Err:     err.Error(),
			})
			return
		}
	} else if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// with token received, have request contains token information
	claims, ok := token.Claims.(jwtPkg.MapClaims)
	if !ok {
		appvendor.InternalServerError(w, "token's claims should be type MapClaims")
	}

	key, ok := claims["Id"].(string)
	if !ok {
		appvendor.InternalServerError(w, "the claims should have attribute 'Id' of string")
	} else if key == "" {
		appvendor.InternalServerError(w, "the claims should have attribute 'Id' of string")
	}

	r.Header.Set(appvendor.PropEmail, key)

	// check if token has error 'account not confirmed'
	if _, ok := claims[appvendor.PropJwtError]; ok {
		r.Header.Set(appvendor.PropJwtError, string(appvendor.ErrorAccountNotConfirmed))
	}

	// serve next
	jwt.next.ServeHTTP(w, r)
}

//////////////////////////////////////////////////
// confirmation middleware
//////////////////////////////////////////////////

type ConfirmingHandler struct {
	next http.Handler
	// method to parse request, return true if request should be skipped for jwt token validation
	except func(r *http.Request) bool
}

func NewConfirmingHandler() *ConfirmingHandler {
	return &ConfirmingHandler{}
}

func (cnf *ConfirmingHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// if except function return true for the request, then skip check jwt token
	if cnf.except(r) {
		cnf.next.ServeHTTP(w, r)
		return
	}

	// default response to error
	response := models.ServerResponse{
		Success: false,
		ErrCode: appvendor.ErrorGeneral,
	}

	// if previous middleware (jwt) has error 'account not confirmed', reject further request
	if r.Header.Get(appvendor.PropJwtError) == string(appvendor.ErrorAccountNotConfirmed) {
		response.ErrCode = appvendor.ErrorAccountNotConfirmed
		response.Err = "Account not confirmed"
		json.NewEncoder(w).Encode(response)
		return
	}

	// serve next
	cnf.next.ServeHTTP(w, r)
}

//////////////////////////////////////////////////
// No resource middleware
//////////////////////////////////////////////////

type NoResourceHandler struct {
	next http.Handler
	// method to parse request, return true if request should be skipped for jwt token validation
	router *mux.Router
}

func NewNoResourceHandler() *NoResourceHandler {
	return &NoResourceHandler{}
}

func (nores *NoResourceHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var match mux.RouteMatch

	if nores.router.Match(r, &match) == false {
		http.NotFoundHandler().ServeHTTP(w, r)
		return
	}

	// serve next
	nores.next.ServeHTTP(w, r)
}
