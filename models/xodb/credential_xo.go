// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// Credential represents a row from 'app_mvp_dating.credentials'.
type Credential struct {
	CredID       uint   `json:"cred_id"`       // cred_id
	DeviceName   string `json:"device_name"`   // device_name
	RefreshToken string `json:"refresh_token"` // refresh_token
	UserID       uint   `json:"user_id"`       // user_id
	Deleted      bool   `json:"deleted"`       // deleted

	// xo fields
	_exists bool
}

// Exists determines if the Credential exists in the database.
func (c *Credential) Exists() bool {
	return c._exists
}

// Deleted provides information if the Credential has been deleted from the database.
func (c *Credential) IsDeleted() bool {
	return c.Deleted == true
}

// Insert inserts the Credential to the database.
func (c *Credential) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if c._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO app_mvp_dating.credentials (` +
		`device_name, refresh_token, user_id, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, c.DeviceName, c.RefreshToken, c.UserID, c.Deleted)
	res, err := db.Exec(sqlstr, c.DeviceName, c.RefreshToken, c.UserID, c.Deleted)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	c.CredID = uint(id)
	c._exists = true

	return nil
}

// Update updates the Credential in the database.
func (c *Credential) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if c.Deleted == true {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.credentials SET ` +
		`device_name = ?, refresh_token = ?, user_id = ?, deleted = ?` +
		` WHERE cred_id = ?`

	// run query
	XOLog(sqlstr, c.DeviceName, c.RefreshToken, c.UserID, c.Deleted, c.CredID)
	_, err = db.Exec(sqlstr, c.DeviceName, c.RefreshToken, c.UserID, c.Deleted, c.CredID)
	return err
}

// Save saves the Credential to the database.
func (c *Credential) Save(db XODB) error {
	if c.Exists() {
		return c.Update(db)
	}

	return c.Insert(db)
}

// Delete deletes the Credential from the database.
func (c *Credential) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !c._exists {
		return nil
	}

	// if deleted, bail
	if c.Deleted == true {
		return nil
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.credentials SET deleted = ? WHERE cred_id = ?`

	// run query
	XOLog(sqlstr, true, c.CredID)
	_, err = db.Exec(sqlstr, true, c.CredID)
	if err != nil {
		return err
	}

	// set deleted
	c.Deleted = true

	return nil
}

// CredentialByCredID retrieves a row from 'app_mvp_dating.credentials' as a Credential.
//
// Generated from index 'credentials_cred_id_pkey'.
func CredentialByCredID(db XODB, credID uint) (*Credential, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`cred_id, device_name, refresh_token, user_id, deleted ` +
		`FROM app_mvp_dating.credentials ` +
		`WHERE cred_id = ?`

	// run query
	XOLog(sqlstr, credID)
	c := Credential{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, credID).Scan(&c.CredID, &c.DeviceName, &c.RefreshToken, &c.UserID, &c.Deleted)
	if err != nil {
		return nil, err
	}

	return &c, nil
}

// CredentialByRefreshToken retrieves a row from 'app_mvp_dating.credentials' as a Credential.
func CredentialByRefreshToken(db XODB, refreshToken string, deviceName string) (*Credential, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`cred_id, device_name, refresh_token, user_id, deleted+0 ` +
		`FROM app_mvp_dating.credentials ` +
		`WHERE refresh_token = ? AND device_name = ?`

	// run query
	XOLog(sqlstr, refreshToken, deviceName)
	c := Credential{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, refreshToken, deviceName).Scan(&c.CredID, &c.DeviceName, &c.RefreshToken, &c.UserID, &c.Deleted)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
