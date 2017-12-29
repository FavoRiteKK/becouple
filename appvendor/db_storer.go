package appvendor

import (
	"becouple/models/xodb"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
)

type IDBStorer interface {
	SaveCredential(refreshToken, key, deviceName string) error
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

func (s AuthStorer) Create(_ string, attr authboss.Attributes) error {
	user := xodb.NewLegalUser()
	BindAuthbossUser(user, attr)

	// save to db
	err := s.dbHelper.Insert(user)
	if err != nil {
		logrus.WithError(err).Errorln("error with insert user query")
		return err
	}

	return nil
}

func (s AuthStorer) Put(key string, attr authboss.Attributes) error {
	user, err := s.dbHelper.GetUserByEmail(key)
	if err != nil {
		logrus.WithError(err).Errorln("cannot bind attribute to user")
		return err
	}

	BindAuthbossUser(user, attr)

	err = s.dbHelper.SaveUser(user)
	if err != nil {
		logrus.WithError(err).Errorln("error with save user query")
		return err
	}

	return nil
}

func (s AuthStorer) Get(key string) (result interface{}, err error) {
	user, err := s.dbHelper.GetUserByEmail(key)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, authboss.ErrUserNotFound
		}

		return nil, err
	}

	return user, nil
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

func (s AuthStorer) DeleteUser(user *xodb.User) error {
	return s.dbHelper.DeleteUser(user)
}

func (s AuthStorer) DeletePermanently(user *xodb.User) error {
	return s.dbHelper.DeletePermanently(user)
}

func (s AuthStorer) SaveCredential(refreshToken, key, deviceName string) error {

	credential := &xodb.Credential{
		RefreshToken: refreshToken,
		Email:        key,
		DeviceName:   deviceName,
	}

	return s.dbHelper.SaveCredential(credential)
}
