//+build USE_DB_AUTH_STORER

// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"

	"github.com/go-sql-driver/mysql"
)

// UserAct represents a row from 'app_mvp_dating.user_acts'.
type UserAct struct {
	UserID       uint           `json:"user_id"`        // user_id
	TargetUserID uint           `json:"target_user_id"` // target_user_id
	Likes        sql.NullBool   `json:"likes"`          // likes
	VisitDate    mysql.NullTime `json:"visit_date"`     // visit_date
	SeenCount    sql.NullInt64  `json:"seen_count"`     // seen_count
	Deleted      sql.NullBool   `json:"deleted"`        // deleted

	// xo fields
	_exists bool
}

// Exists determines if the UserAct exists in the database.
func (ua *UserAct) Exists() bool {
	return ua._exists
}

// Deleted provides information if the UserAct has been deleted from the database.
func (ua *UserAct) IsDeleted() bool {
	return ua.Deleted.Valid && ua.Deleted.Bool == true
}

// Insert inserts the UserAct to the database.
func (ua *UserAct) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if ua._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key must be provided
	const sqlstr = `INSERT INTO app_mvp_dating.user_acts (` +
		`user_id, target_user_id, likes, visit_date, seen_count, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, ua.UserID, ua.TargetUserID, ua.Likes, ua.VisitDate, ua.SeenCount, ua.Deleted)
	_, err = db.Exec(sqlstr, ua.UserID, ua.TargetUserID, ua.Likes, ua.VisitDate, ua.SeenCount, ua.Deleted)
	if err != nil {
		return err
	}

	// set existence
	ua._exists = true

	return nil
}

// Update updates the UserAct in the database.
func (ua *UserAct) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ua._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if ua.Deleted.Valid && ua.Deleted.Bool == true {
		return errors.New("update failed: marked for deletion")
	}

	// sql query with composite primary key
	const sqlstr = `UPDATE app_mvp_dating.user_acts SET ` +
		`likes = ?, visit_date = ?, seen_count = ?, deleted = ?` +
		` WHERE user_id = ? AND target_user_id = ?`

	// run query
	XOLog(sqlstr, ua.Likes, ua.VisitDate, ua.SeenCount, ua.Deleted, ua.UserID, ua.TargetUserID)
	_, err = db.Exec(sqlstr, ua.Likes, ua.VisitDate, ua.SeenCount, ua.Deleted, ua.UserID, ua.TargetUserID)
	return err
}

// Save saves the UserAct to the database.
func (ua *UserAct) Save(db XODB) error {
	if ua.Exists() {
		return ua.Update(db)
	}

	return ua.Insert(db)
}

// Delete deletes the UserAct from the database.
func (ua *UserAct) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !ua._exists {
		return nil
	}

	// if deleted, bail
	if ua.Deleted.Valid && ua.Deleted.Bool == true {
		return nil
	}

	// sql query with composite primary key
	const sqlstr = `UPDATE app_mvp_dating.user_acts SET deleted = ? WHERE user_id = ? AND target_user_id = ?`

	// run query
	XOLog(sqlstr, true, ua.UserID, ua.TargetUserID)
	_, err = db.Exec(sqlstr, true, ua.UserID, ua.TargetUserID)
	if err != nil {
		return err
	}

	// set deleted
	ua.Deleted = sql.NullBool{Bool: true, Valid: true}

	return nil
}

// UserByUserID returns the User associated with the UserAct's UserID (user_id).
//
// Generated from foreign key 'fk_user_likes_1'.
func (ua *UserAct) UserByUserID(db XODB) (*User, error) {
	return UserByUserID(db, ua.UserID)
}

// UserByTargetUserID returns the User associated with the UserAct's TargetUserID (target_user_id).
//
// Generated from foreign key 'fk_user_likes_2'.
func (ua *UserAct) UserByTargetUserID(db XODB) (*User, error) {
	return UserByUserID(db, ua.TargetUserID)
}

// UserActsByTargetUserID retrieves a row from 'app_mvp_dating.user_acts' as a UserAct.
//
// Generated from index 'fk_user_likes_2_idx'.
func UserActsByTargetUserID(db XODB, targetUserID uint) ([]*UserAct, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, target_user_id, likes, visit_date, seen_count, deleted ` +
		`FROM app_mvp_dating.user_acts ` +
		`WHERE target_user_id = ?`

	// run query
	XOLog(sqlstr, targetUserID)
	q, err := db.Query(sqlstr, targetUserID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*UserAct{}
	for q.Next() {
		ua := UserAct{
			_exists: true,
		}

		// scan
		err = q.Scan(&ua.UserID, &ua.TargetUserID, &ua.Likes, &ua.VisitDate, &ua.SeenCount, &ua.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &ua)
	}

	return res, nil
}

// UserActByTargetUserID retrieves a row from 'app_mvp_dating.user_acts' as a UserAct.
//
// Generated from index 'user_acts_target_user_id_pkey'.
func UserActByTargetUserID(db XODB, targetUserID uint) (*UserAct, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, target_user_id, likes, visit_date, seen_count, deleted ` +
		`FROM app_mvp_dating.user_acts ` +
		`WHERE target_user_id = ?`

	// run query
	XOLog(sqlstr, targetUserID)
	ua := UserAct{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, targetUserID).Scan(&ua.UserID, &ua.TargetUserID, &ua.Likes, &ua.VisitDate, &ua.SeenCount, &ua.Deleted)
	if err != nil {
		return nil, err
	}

	return &ua, nil
}