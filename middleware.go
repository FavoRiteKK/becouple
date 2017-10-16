package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/justinas/nosurf"
    jwtPkg "github.com/dgrijalva/jwt-go"
    "regexp"
    "github.com/dgrijalva/jwt-go/request"
    "github.com/gin-gonic/gin/json"
)

type authProtector struct {
	f http.HandlerFunc
}

func authProtect(f http.HandlerFunc) authProtector {
	return authProtector{f}
}

func (ap authProtector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if u, err := ab.CurrentUser(w, r); err != nil {
		log.Println("Error fetching current user:", err)
		w.WriteHeader(http.StatusInternalServerError)
	} else if u == nil {
		log.Println("Redirecting unauthorized user from:", r.URL.Path)
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
			log.Println("Failed to validate XSRF Token:", nosurf.Reason(r))
			w.WriteHeader(http.StatusBadRequest)
		}))

		return surfing
	}
}

func logger(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("\n%s %s %s\n", r.Method, r.URL.Path, r.Proto)
		session, err := sessionStore.Get(r, sessionCookieName)
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
		fmt.Print("Database: ")
		for _, u := range database.Users {
			fmt.Printf("%#v\n", u)
		}
		h.ServeHTTP(w, r)
	})
}

func jwtMiddleware() func(next http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        jwtAuth := NewJwtAuth(next)

        // regex to check if route contains '/api/'
        // if it does, the middleware would check jwt token
        regex := regexp.MustCompile("/api/")

        jwtAuth.token = jwtPkg.New(appJwtSigningMethod)
        jwtAuth.next = next
        jwtAuth.except = func(r *http.Request) bool {
            path := r.URL.Path

            // exempt this exactly path
            if path == "/api/auth" {
                return true
            }

            return !regex.MatchString(path)
        }
        return jwtAuth
    }
}

//////////////////////////////////////////////////
// jwtAuth middleware
//////////////////////////////////////////////////
var (
    appJwtSigningMethod = jwtPkg.SigningMethodHS256
    appJwtSecret = "qweasd123"
)

type JwtAuth struct {
    token *jwtPkg.Token
	next http.Handler
	// method to parse request, return true if request should be skipped for jwt token validation
	except func(r *http.Request) bool
}

func NewJwtAuth(next http.Handler) *JwtAuth {
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
    bearer, err := request.ParseFromRequest(r, request.OAuth2Extractor, func(token *jwtPkg.Token) (interface{}, error) {
        b := []byte(appJwtSecret)
        return b, nil
    })

    if err != nil {
        http.Error(w, err.Error(), http.StatusUnauthorized)
        return
    }

    jb, err := json.Marshal(bearer)
    if err != nil {
        http.Error(w, "Error marshal bearer", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(jb)

    jwt.next.ServeHTTP(w, r)
}

//////////////////////////////////////////////////
//////////////////////////////////////////////////
