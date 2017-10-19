//+build USE_DB_AUTH_STORER

package appvendor

import (
	"gopkg.in/authboss.v1"
	"log"
)

type AuthStorer struct {
	dbHelper DBManager
}

func NewAuthStorer() *AuthStorer {
	return &AuthStorer{DBHelper}
}

func (s AuthStorer) Create(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	return nil
}

func (s AuthStorer) Put(key string, attr authboss.Attributes) error {
	return s.Create(key, attr)
}

func (s AuthStorer) Get(key string) (result interface{}, err error) {
	row, err := s.dbHelper.GetUserByEmail(key)
	if err != nil {
		log.Print(err.Error())
		return nil, err
	}

	var user User

	err = row.Scan(&user.ID, &user.Email, &user.Password)
	if err != nil {
		log.Print(err.Error())
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
