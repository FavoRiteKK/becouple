package appvendor

import (
	"becouple/models/xodb"
	"database/sql"
	"errors"
	"time"

	_ "github.com/go-sql-driver/mysql" // init mysql driver
)

//IDBManager describes db manager interface
type IDBManager interface {
	Connect() error
	HasConn() (bool, error)

	Insert(user *xodb.User) error
	GetUserByEmail(email string) (*xodb.User, error)

	SaveUser(user *xodb.User) error

	DeleteUser(user *xodb.User) error
	DeletePermanently(user *xodb.User) error

	SaveCredential(credential *xodb.Credential) error
	GetCredentialByRefreshToken(refreshToken string, deviceName string) (*xodb.Credential, error)
	SaveUserPhoto(photo *xodb.UserPhoto) error
}

type manager struct {
	db *sql.DB
}

func (mgr *manager) Connect() error {

	if ok, _ := mgr.HasConn(); ok {
		return nil
	}

	// open global database handler
	db, err := sql.Open("mysql", "root:qweasdzxc@123@/app_mvp_dating")
	if err != nil {
		return err
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		return err
	}

	// set 5 minutes long lived connection
	if db != nil {
		db.SetConnMaxLifetime(time.Minute * 5)
		db.SetMaxIdleConns(0)
		db.SetMaxOpenConns(5)
	}

	mgr.db = db
	return nil
}

func (mgr *manager) HasConn() (bool, error) {
	if mgr.db == nil {
		return false, errors.New("doesn't have database connection")
	}

	err := mgr.db.Ping()
	if err != nil {
		return false, err
	}

	return true, nil
}

func (mgr *manager) Insert(user *xodb.User) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return user.Save(mgr.db)
}

func (mgr *manager) GetUserByEmail(email string) (*xodb.User, error) {
	if ok, err := mgr.HasConn(); !ok {
		return nil, err
	}

	return xodb.UserByEmail(mgr.db, email)
}

func (mgr *manager) SaveUser(user *xodb.User) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return user.Save(mgr.db)
}

func (mgr *manager) DeleteUser(user *xodb.User) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return user.Delete(mgr.db)
}

func (mgr *manager) DeletePermanently(user *xodb.User) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return user.DeletePermanently(mgr.db)
}

func (mgr *manager) SaveCredential(credential *xodb.Credential) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return credential.Save(mgr.db)
}

func (mgr *manager) GetCredentialByRefreshToken(refreshToken string, deviceName string) (*xodb.Credential, error) {
	if ok, err := mgr.HasConn(); !ok {
		return nil, err
	}

	return xodb.CredentialByRefreshToken(mgr.db, refreshToken, deviceName)
}

func (mgr *manager) SaveUserPhoto(photo *xodb.UserPhoto) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	return photo.Save(mgr.db)
}
