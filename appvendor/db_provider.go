package appvendor

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"time"
    "log"
)

type DBManager interface {
	// Add other methods
	Insert(email string, password string, fullname string) (sql.Result, error)
	GetUserByEmail(email string) (*sql.Row, error)
}

type manager struct {
	db *sql.DB
}

var DBHelper DBManager

func init() {
	// open global database handler
	db, err := sql.Open("mysql", "root:qweasdzxc@123@/app_mvp_dating")
	if err != nil {
		log.Fatal(err.Error()) // Just for example purpose. You should use proper error handling instead of panic
	}

	// Open doesn't open a connection. Validate DSN data:
	err = db.Ping()
	if err != nil {
		log.Fatal(err.Error()) // proper error handling instead of panic in your app
	}

	// set 5 minutes long lived connection
	db.SetConnMaxLifetime(time.Minute * 5)
	db.SetMaxIdleConns(0)
	db.SetMaxOpenConns(5)

	DBHelper = &manager{db: db}
}

func (mgr *manager) Insert(email string, password string, fullname string) (sql.Result, error) {
    stmt, err := mgr.db.Prepare("INSERT IGNORE INTO `user` (`password`, `fullname`) VALUES(?, ?, ?)")
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
