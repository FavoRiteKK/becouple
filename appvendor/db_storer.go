package appvendor

import (
	"gopkg.in/authboss.v1"
	"log"
)

type MySQLStorer struct {
	dbHelper DBManager
}

func NewMySQLStorer() *MySQLStorer {
	return &MySQLStorer{DBHelper}
}

func (s MySQLStorer) Create(key string, attr authboss.Attributes) error {
	var user User
	if err := attr.Bind(&user, true); err != nil {
		return err
	}

	return nil
}

func (s MySQLStorer) Put(key string, attr authboss.Attributes) error {
	return s.Create(key, attr)
}

func (s MySQLStorer) Get(key string) (result interface{}, err error) {
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

func (s MySQLStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

func (s MySQLStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	return nil, nil
}

func (s MySQLStorer) AddToken(key, token string) error {
	return nil
}

func (s MySQLStorer) DelTokens(key string) error {
	return nil
}

func (s MySQLStorer) UseToken(givenKey, token string) error {
	return authboss.ErrTokenNotFound
}

func (s MySQLStorer) ConfirmUser(tok string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}

func (s MySQLStorer) RecoverUser(rec string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}
