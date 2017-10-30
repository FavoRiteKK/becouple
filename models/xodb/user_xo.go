// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"errors"

	"github.com/go-sql-driver/mysql"
)

// User represents a row from 'app_mvp_dating.user'.
type User struct {
	UserID             uint           `json:"user_id"`              // user_id
	Email              string         `json:"email"`                // email
	Password           string         `json:"password"`             // password
	Fullname           string         `json:"fullname"`             // fullname
	Nickname           string         `json:"nickname"`             // nickname
	AvatarURI          string         `json:"avatar_uri"`           // avatar_uri
	PhoneNumber        string         `json:"phone_number"`         // phone_number
	Gender             NullGender     `json:"gender"`               // gender
	DateOfBirth        mysql.NullTime `json:"date_of_birth"`        // date_of_birth
	Job                string         `json:"job"`                  // job
	LivingAt           string         `json:"living_at"`            // living_at
	HomeTown           string         `json:"home_town"`            // home_town
	WorkingAt          string         `json:"working_at"`           // working_at
	ShortAbout         string         `json:"short_about"`          // short_about
	Height             uint8          `json:"height"`               // height
	Weight             uint8          `json:"weight"`               // weight
	Status             NullStatus     `json:"status"`               // status
	Oauth2UID          string         `json:"oauth2_uid"`           // oauth2_uid
	Oauth2Provider     string         `json:"oauth2_provider"`      // oauth2_provider
	Oauth2Token        string         `json:"oauth2_token"`         // oauth2_token
	Oauth2Refresh      string         `json:"oauth2_refresh"`       // oauth2_refresh
	Oauth2Expiry       mysql.NullTime `json:"oauth2_expiry"`        // oauth2_expiry
	ConfirmToken       string         `json:"confirm_token"`        // confirm_token
	Confirmed          bool           `json:"confirmed"`            // confirmed
	AttemptNumber      uint8          `json:"attempt_number"`       // attempt_number
	AttemptTime        mysql.NullTime `json:"attempt_time"`         // attempt_time
	Locked             mysql.NullTime `json:"locked"`               // locked
	RecoverToken       string         `json:"recover_token"`        // recover_token
	RecoverTokenExpiry mysql.NullTime `json:"recover_token_expiry"` // recover_token_expiry
	Deleted            bool           `json:"deleted"`              // deleted
	// xo fields
	_exists bool
}

// if field type is enum, the object can't be insert with default value 0, so we must set default value here
// to create new object (mostly for inserting)
func NewLegalUser() *User {
	return &User{}
}

// Exists determines if the User exists in the database.
func (u *User) Exists() bool {
	return u._exists
}

// Deleted provides information if the User has been deleted from the database.
func (u *User) IsDeleted() bool {
	return u.Deleted == true
}

// Insert inserts the User to the database.
func (u *User) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if u._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO app_mvp_dating.user (` +
		`email, password, fullname, nickname, avatar_uri, phone_number, gender, date_of_birth, job, living_at, home_town, working_at, short_about, height, weight, status, oauth2_uid, oauth2_provider, oauth2_token, oauth2_refresh, oauth2_expiry, confirm_token, confirmed, attempt_number, attempt_time, locked, recover_token, recover_token_expiry, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, u.Email, u.Password, u.Fullname, u.Nickname, u.AvatarURI, u.PhoneNumber, u.Gender, u.DateOfBirth, u.Job, u.LivingAt, u.HomeTown, u.WorkingAt, u.ShortAbout, u.Height, u.Weight, u.Status, u.Oauth2UID, u.Oauth2Provider, u.Oauth2Token, u.Oauth2Refresh, u.Oauth2Expiry, u.ConfirmToken, u.Confirmed, u.AttemptNumber, u.AttemptTime, u.Locked, u.RecoverToken, u.RecoverTokenExpiry, u.Deleted)
	res, err := db.Exec(sqlstr, u.Email, u.Password, u.Fullname, u.Nickname, u.AvatarURI, u.PhoneNumber, u.Gender, u.DateOfBirth, u.Job, u.LivingAt, u.HomeTown, u.WorkingAt, u.ShortAbout, u.Height, u.Weight, u.Status, u.Oauth2UID, u.Oauth2Provider, u.Oauth2Token, u.Oauth2Refresh, u.Oauth2Expiry, u.ConfirmToken, u.Confirmed, u.AttemptNumber, u.AttemptTime, u.Locked, u.RecoverToken, u.RecoverTokenExpiry, u.Deleted)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	u.UserID = uint(id)
	u._exists = true

	return nil
}

// Update updates the User in the database.
func (u *User) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if u.Deleted == true {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.user SET ` +
		`email = ?, password = ?, fullname = ?, nickname = ?, avatar_uri = ?, phone_number = ?, gender = ?, date_of_birth = ?, job = ?, living_at = ?, home_town = ?, working_at = ?, short_about = ?, height = ?, weight = ?, status = ?, oauth2_uid = ?, oauth2_provider = ?, oauth2_token = ?, oauth2_refresh = ?, oauth2_expiry = ?, confirm_token = ?, confirmed = ?, attempt_number = ?, attempt_time = ?, locked = ?, recover_token = ?, recover_token_expiry = ?, deleted = ?` +
		` WHERE user_id = ?`

	// run query
	XOLog(sqlstr, u.Email, u.Password, u.Fullname, u.Nickname, u.AvatarURI, u.PhoneNumber, u.Gender, u.DateOfBirth, u.Job, u.LivingAt, u.HomeTown, u.WorkingAt, u.ShortAbout, u.Height, u.Weight, u.Status, u.Oauth2UID, u.Oauth2Provider, u.Oauth2Token, u.Oauth2Refresh, u.Oauth2Expiry, u.ConfirmToken, u.Confirmed, u.AttemptNumber, u.AttemptTime, u.Locked, u.RecoverToken, u.RecoverTokenExpiry, u.Deleted, u.UserID)
	_, err = db.Exec(sqlstr, u.Email, u.Password, u.Fullname, u.Nickname, u.AvatarURI, u.PhoneNumber, u.Gender, u.DateOfBirth, u.Job, u.LivingAt, u.HomeTown, u.WorkingAt, u.ShortAbout, u.Height, u.Weight, u.Status, u.Oauth2UID, u.Oauth2Provider, u.Oauth2Token, u.Oauth2Refresh, u.Oauth2Expiry, u.ConfirmToken, u.Confirmed, u.AttemptNumber, u.AttemptTime, u.Locked, u.RecoverToken, u.RecoverTokenExpiry, u.Deleted, u.UserID)
	return err
}

// Save saves the User to the database.
func (u *User) Save(db XODB) error {
	if u.Exists() {
		return u.Update(db)
	}

	return u.Insert(db)
}

// Delete deletes the User from the database.
func (u *User) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u.Deleted == true {
		return nil
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.user SET deleted = ? WHERE user_id = ?`

	// run query
	XOLog(sqlstr, true, u.UserID)
	_, err = db.Exec(sqlstr, true, u.UserID)
	if err != nil {
		return err
	}

	// set deleted
	u.Deleted = true

	return nil
}

// Delete permanently the User from the database.
func (u *User) DeletePermanently(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !u._exists {
		return nil
	}

	// if deleted, bail
	if u.Deleted == true {
		return nil
	}

	// sql query
	const sqlstr = `DELETE FROM app_mvp_dating.user WHERE user_id = ?`

	// run query
	XOLog(sqlstr, u.UserID)
	_, err = db.Exec(sqlstr, u.UserID)
	if err != nil {
		return err
	}

	// set deleted
	u.Deleted = true

	return nil
}

// UserByEmail retrieves a row from 'app_mvp_dating.user' as a User.
//
// Generated from index 'email_UNIQUE'.
func UserByEmail(db XODB, email string) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password, fullname, nickname, avatar_uri, phone_number, gender, date_of_birth, job, living_at, home_town, working_at, short_about, height, weight, status, oauth2_uid, oauth2_provider, oauth2_token, oauth2_refresh, oauth2_expiry, confirm_token, confirmed+0, attempt_number, attempt_time, locked, recover_token, recover_token_expiry, deleted+0 ` +
		`FROM app_mvp_dating.user ` +
		`WHERE email = ?`

	// run query
	XOLog(sqlstr, email)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, email).Scan(&u.UserID, &u.Email, &u.Password, &u.Fullname, &u.Nickname, &u.AvatarURI, &u.PhoneNumber, &u.Gender, &u.DateOfBirth, &u.Job, &u.LivingAt, &u.HomeTown, &u.WorkingAt, &u.ShortAbout, &u.Height, &u.Weight, &u.Status, &u.Oauth2UID, &u.Oauth2Provider, &u.Oauth2Token, &u.Oauth2Refresh, &u.Oauth2Expiry, &u.ConfirmToken, &u.Confirmed, &u.AttemptNumber, &u.AttemptTime, &u.Locked, &u.RecoverToken, &u.RecoverTokenExpiry, &u.Deleted)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UserByPhoneNumber retrieves a row from 'app_mvp_dating.user' as a User.
//
// Generated from index 'phone_number_UNIQUE'.
func UserByPhoneNumber(db XODB, phoneNumber string) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password, fullname, nickname, avatar_uri, phone_number, gender, date_of_birth, job, living_at, home_town, working_at, short_about, height, weight, status, oauth2_uid, oauth2_provider, oauth2_token, oauth2_refresh, oauth2_expiry, confirm_token, confirmed, attempt_number, attempt_time, locked, recover_token, recover_token_expiry, deleted ` +
		`FROM app_mvp_dating.user ` +
		`WHERE phone_number = ?`

	// run query
	XOLog(sqlstr, phoneNumber)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, phoneNumber).Scan(&u.UserID, &u.Email, &u.Password, &u.Fullname, &u.Nickname, &u.AvatarURI, &u.PhoneNumber, &u.Gender, &u.DateOfBirth, &u.Job, &u.LivingAt, &u.HomeTown, &u.WorkingAt, &u.ShortAbout, &u.Height, &u.Weight, &u.Status, &u.Oauth2UID, &u.Oauth2Provider, &u.Oauth2Token, &u.Oauth2Refresh, &u.Oauth2Expiry, &u.ConfirmToken, &u.Confirmed, &u.AttemptNumber, &u.AttemptTime, &u.Locked, &u.RecoverToken, &u.RecoverTokenExpiry, &u.Deleted)
	if err != nil {
		return nil, err
	}

	return &u, nil
}

// UserByUserID retrieves a row from 'app_mvp_dating.user' as a User.
//
// Generated from index 'user_user_id_pkey'.
func UserByUserID(db XODB, userID uint) (*User, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`user_id, email, password, fullname, nickname, avatar_uri, phone_number, gender, date_of_birth, job, living_at, home_town, working_at, short_about, height, weight, status, oauth2_uid, oauth2_provider, oauth2_token, oauth2_refresh, oauth2_expiry, confirm_token, confirmed, attempt_number, attempt_time, locked, recover_token, recover_token_expiry, deleted ` +
		`FROM app_mvp_dating.user ` +
		`WHERE user_id = ?`

	// run query
	XOLog(sqlstr, userID)
	u := User{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, userID).Scan(&u.UserID, &u.Email, &u.Password, &u.Fullname, &u.Nickname, &u.AvatarURI, &u.PhoneNumber, &u.Gender, &u.DateOfBirth, &u.Job, &u.LivingAt, &u.HomeTown, &u.WorkingAt, &u.ShortAbout, &u.Height, &u.Weight, &u.Status, &u.Oauth2UID, &u.Oauth2Provider, &u.Oauth2Token, &u.Oauth2Refresh, &u.Oauth2Expiry, &u.ConfirmToken, &u.Confirmed, &u.AttemptNumber, &u.AttemptTime, &u.Locked, &u.RecoverToken, &u.RecoverTokenExpiry, &u.Deleted)
	if err != nil {
		return nil, err
	}

	return &u, nil
}
