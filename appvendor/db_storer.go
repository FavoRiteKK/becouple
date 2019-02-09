package appvendor

import (
	"becouple/models/xodb"
	"database/sql"

	"github.com/sirupsen/logrus"
	"github.com/volatiletech/authboss"
)

//IDBStorer db storer logic interface
type IDBStorer interface {
	SaveCredential(refreshToken string, userID uint, deviceName string) error
	SavePhoto(uri string, userID uint64) error
	DeleteUser(user *xodb.User) error
	DeletePermanently(user *xodb.User) error
	GetCredentialByRefreshToken(refreshToken string, deviceName string) (*xodb.Credential, error)
	GetUserByID(userID uint) (*xodb.User, error)
}

//AuthStorer authentication storer
type AuthStorer struct {
	dbHelper IDBManager
}

//NewAuthStorer returns instance of AuthStorer
func NewAuthStorer() *AuthStorer {
	dbHelper := new(manager)
	// connect database first
	if err := dbHelper.Connect(); err != nil {
		logrus.WithError(err).Errorln("error connecting database")
	}

	// then instantiate our store object
	return &AuthStorer{dbHelper}
}

//Create creates/insert an user to db
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

//Put saves/update user
func (s AuthStorer) Put(key string, attr authboss.Attributes) error {
	user, err := s.dbHelper.GetUserByEmail(key)
	if err != nil {
		logrus.WithError(err).Errorln("cannot save/update user")
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

//Get returns user by email parameter
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

//PutOAuth creates/inserts user from uid and provider
func (s AuthStorer) PutOAuth(uid, provider string, attr authboss.Attributes) error {
	return s.Create(uid+provider, attr)
}

//GetOAuth returns user by uid and provider
func (s AuthStorer) GetOAuth(uid, provider string) (result interface{}, err error) {
	return nil, nil
}

//AddToken adds token to db
func (s AuthStorer) AddToken(key, token string) error {
	return nil
}

//DelTokens deletes token from db
func (s AuthStorer) DelTokens(key string) error {
	return nil
}

//UseToken uses token
func (s AuthStorer) UseToken(givenKey, token string) error {
	return authboss.ErrTokenNotFound
}

//ConfirmUser marks the user as confirmed (ex: his email is valid)
func (s AuthStorer) ConfirmUser(tok string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}

//RecoverUser recovers user's password (ex: send an email with recovery code)
func (s AuthStorer) RecoverUser(rec string) (result interface{}, err error) {
	return nil, authboss.ErrUserNotFound
}

//DeleteUser marks user as deleted in db
func (s AuthStorer) DeleteUser(user *xodb.User) error {
	return s.dbHelper.DeleteUser(user)
}

//DeletePermanently permanently removes user form db
func (s AuthStorer) DeletePermanently(user *xodb.User) error {
	return s.dbHelper.DeletePermanently(user)
}

//SaveCredential saves credential with refresh token and key and device name
func (s AuthStorer) SaveCredential(refreshToken string, userID uint, deviceName string) error {

	credential := &xodb.Credential{
		RefreshToken: refreshToken,
		UserID:       userID,
		DeviceName:   deviceName,
	}

	return s.dbHelper.SaveCredential(credential)
}

//GetCredentialByRefreshToken returns credential by refresh token and device name
func (s AuthStorer) GetCredentialByRefreshToken(refreshToken string, deviceName string) (*xodb.Credential, error) {
	cred, err := s.dbHelper.GetCredentialByRefreshToken(refreshToken, deviceName)
	if err != nil {
		logrus.WithError(err).Errorln("error with get credential method")
	}
	return cred, err
}

// SavePhoto saves user photo to db
func (s AuthStorer) SavePhoto(uri string, userID uint) error {
	userPhoto := &xodb.UserPhoto{
		PhotoURI: uri,
		UserID:   userID,
	}

	return s.dbHelper.SaveUserPhoto(userPhoto)
}
