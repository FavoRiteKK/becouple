//+build USE_DB_AUTH_STORER

package appvendor

import (
	"database/sql"
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"github.com/sirupsen/logrus"
	"gopkg.in/authboss.v1"
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
	// connect database first
	if err := DBHelper.Connect(); err != nil {
		logrus.WithError(err).Errorln("error connecting database")
	}

	// then instantiate our store object
	return &AuthStorer{DBHelper}
}

func (s AuthStorer) Create(key string, attr authboss.Attributes) error {
	var user AuthUser
	if err := attr.Bind(&user, true); err != nil {
		logrus.WithError(err).Errorln("cannot bind attribute to user")
		return err
	}

	//TODO get user's fullname somehow

	// save to db
	result, err := s.dbHelper.Insert(user.Email, user.Password, "Anonymous")
	if err != nil {
		logrus.WithError(err).Errorln("error with insert user query")
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
		logrus.WithError(err).Errorln("error with get user query")
		return nil, err
	}

	var user AuthUser

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
        logrus.WithError(err).Errorln("error scanning user")

		if err == sql.ErrNoRows {
			err = authboss.ErrUserNotFound
		}

		return nil, err
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
