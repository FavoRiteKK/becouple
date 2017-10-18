package appvendor

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
    "time"
)

type Manager interface {
	//AddArticle(article *article.Article) error
	GetAllUser() error
	// Add other methods
}

type manager struct {
	db *sql.DB
}

var Mgr Manager

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

	Mgr = &manager{db: db}
}

//func (mgr *manager) AddArticle(article *article.Article) (err error) {
//    mgr.db.Create(article)
//    if errs := mgr.db.GetErrors(); len(errs) > 0 {
//        err = errs[0]
//    }
//    return
//}

func (mgr *manager) GetAllUser() error {
	result, err := mgr.db.Exec("SELECT * FROM `user`")
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	log.Println(result.RowsAffected())
	return nil
}
