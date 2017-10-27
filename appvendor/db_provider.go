//+build USE_DB_AUTH_STORER

package appvendor

import (
	"becouple/models/xodb"
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DBManager interface {
	Connect() error
	HasConn() (bool, error)
	Insert(email string, password string, fullname string) error
	GetUserByEmail(email string) (*xodb.User, error)
	DeleteUser(user *xodb.User) error
	DeletePermanently(user *xodb.User) error
}

type manager struct {
	db *sql.DB
}

var DBHelper DBManager = new(manager)

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

func (mgr *manager) Insert(email string, password string, fullname string) error {
	if ok, err := mgr.HasConn(); !ok {
		return err
	}

	user := xodb.NewLegalUser()
	user.Email = email
	user.Password = password
	user.Fullname = fullname

	return user.Save(mgr.db)
}

func (mgr *manager) GetUserByEmail(email string) (*xodb.User, error) {
	if ok, err := mgr.HasConn(); !ok {
		return nil, err
	}

	return xodb.UserByEmail(mgr.db, email)
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
