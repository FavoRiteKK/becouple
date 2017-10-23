//+build USE_DB_AUTH_STORER

package appvendor

import (
	"database/sql"
	"errors"
	_ "github.com/go-sql-driver/mysql"
	"time"
)

type DBManager interface {
	Connect() error
	HasConn() (bool, error)
	Insert(email string, password string, fullname string) (sql.Result, error)
	GetUserByEmail(email string) (*sql.Row, error)
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

func (mgr *manager) Insert(email string, password string, fullname string) (sql.Result, error) {
	if ok, err := mgr.HasConn(); !ok {
		return nil, err
	}

	stmt, err := mgr.db.Prepare("INSERT IGNORE INTO `user` (`email`, `password`, `fullname`) VALUES(?, ?, ?)")
	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec("", password, fullname)
	return result, err
}

func (mgr *manager) GetUserByEmail(email string) (*sql.Row, error) {
	if ok, err := mgr.HasConn(); !ok {
		return nil, err
	}

	stmt, err := mgr.db.Prepare("SELECT `user_id`, `email`, `password` FROM `user` WHERE email = ? LIMIT 1")
	if stmt != nil {
		defer stmt.Close()
	}

	if err != nil {
		return nil, err
	}

	row := stmt.QueryRow(email)
	return row, nil
}
