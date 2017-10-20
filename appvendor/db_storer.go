//+build USE_DB_AUTH_STORER

package appvendor

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"gopkg.in/authboss.v1"
	"log"
	"time"
)

type AuthUser struct {
	ID   int
	Name string

	// Auth
	Email    string
	Password string

	// OAuth2
	Oauth2Uid      string
	Oauth2Provider string
	Oauth2Token    string
	Oauth2Refresh  string
	Oauth2Expiry   time.Time

	// Confirm
	ConfirmToken string
	Confirmed    bool

	// Lock
	AttemptNumber int64
	AttemptTime   time.Time
	Locked        time.Time

	// Recover
	RecoverToken       string
	RecoverTokenExpiry time.Time

	// Remember is in another table
}

type AuthStorer struct {
	dbHelper DBManager
}

func NewAuthStorer() *AuthStorer {
	return &AuthStorer{DBHelper}
}

func (s AuthStorer) Create(key string, attr authboss.Attributes) error {
	var user AuthUser
	if err := attr.Bind(&user, true); err != nil {
        log.Println(err.Error())
		return err
	}

	//TODO get user's fullname somehow

	// save to db
	result, err := s.dbHelper.Insert(user.Email, user.Password, "Anonymous")
	if err != nil {
        log.Println("Error insert query: ", err.Error())
		return err
	}

	fmt.Println("==========> Create user result:")
	spew.Dump(result)

	return nil
}

func (s AuthStorer) Put(key string, attr authboss.Attributes) error {
	return s.Create(key, attr)
}

func (s AuthStorer) Get(key string) (result interface{}, err error) {
	row, err := s.dbHelper.GetUserByEmail(key)
	if err != nil {
		log.Println("Error select query", err.Error())
		return nil, err
	}

	var user AuthUser

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Println(err.Error())
		return nil, authboss.ErrUserNotFound
	}

	return &user, nil

}

func (s AuthStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s AuthStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	return nil, nil
}

func (s AuthStorer) AddToken(key, token string) error {
	return nil
}

func (s AuthStorer) DelTokens(key string) error {
	return nil
}

func (s AuthStorer) UseToken(givenKey, token string) error {
	return authboss.ErrTokenNotFound
}

func (s AuthStorer) ConfirmUser(tok string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}

func (s AuthStorer) RecoverUser(rec string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}
