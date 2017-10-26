//+build USE_DB_AUTH_STORER

// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql"
	"errors"
)

// ConversationMessage represents a row from 'app_mvp_dating.conversation_message'.
type ConversationMessage struct {
	ConversationMessageID uint           `json:"conversation_message_id"` // conversation_message_id
	UserID                uint           `json:"user_id"`                 // user_id
	ConversationID        uint           `json:"conversation_id"`         // conversation_id
	Message               sql.NullString `json:"message"`                 // message
	IP                    sql.NullString `json:"ip"`                      // ip
	Deleted               sql.NullBool   `json:"deleted"`                 // deleted

	// xo fields
	_exists bool
}

// Exists determines if the ConversationMessage exists in the database.
func (cm *ConversationMessage) Exists() bool {
	return cm._exists
}

// Deleted provides information if the ConversationMessage has been deleted from the database.
func (cm *ConversationMessage) IsDeleted() bool {
	return cm.Deleted.Valid && cm.Deleted.Bool == true
}

// Insert inserts the ConversationMessage to the database.
func (cm *ConversationMessage) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if cm._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO app_mvp_dating.conversation_message (` +
		`user_id, conversation_id, message, ip, deleted` +
		`) VALUES (` +
		`?, ?, ?, ?, ?` +
		`)`

	// run query
	XOLog(sqlstr, cm.UserID, cm.ConversationID, cm.Message, cm.IP, cm.Deleted)
	res, err := db.Exec(sqlstr, cm.UserID, cm.ConversationID, cm.Message, cm.IP, cm.Deleted)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	cm.ConversationMessageID = uint(id)
	cm._exists = true

	return nil
}

// Update updates the ConversationMessage in the database.
func (cm *ConversationMessage) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cm._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if cm.Deleted.Valid && cm.Deleted.Bool == true {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.conversation_message SET ` +
		`user_id = ?, conversation_id = ?, message = ?, ip = ?, deleted = ?` +
		` WHERE conversation_message_id = ?`

	// run query
	XOLog(sqlstr, cm.UserID, cm.ConversationID, cm.Message, cm.IP, cm.Deleted, cm.ConversationMessageID)
	_, err = db.Exec(sqlstr, cm.UserID, cm.ConversationID, cm.Message, cm.IP, cm.Deleted, cm.ConversationMessageID)
	return err
}

// Save saves the ConversationMessage to the database.
func (cm *ConversationMessage) Save(db XODB) error {
	if cm.Exists() {
		return cm.Update(db)
	}

	return cm.Insert(db)
}

// Delete deletes the ConversationMessage from the database.
func (cm *ConversationMessage) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !cm._exists {
		return nil
	}

	// if deleted, bail
	if cm.Deleted.Valid && cm.Deleted.Bool == true {
		return nil
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.conversation_message SET deleted = ? WHERE conversation_message_id = ?`

	// run query
	XOLog(sqlstr, true, cm.ConversationMessageID)
	_, err = db.Exec(sqlstr, true, cm.ConversationMessageID)
	if err != nil {
		return err
	}

	// set deleted
	cm.Deleted = sql.NullBool{Bool: true, Valid: true}

	return nil
}

// Conversation returns the Conversation associated with the ConversationMessage's ConversationID (conversation_id).
//
// Generated from foreign key 'fk_conversation_message_1'.
func (cm *ConversationMessage) Conversation(db XODB) (*Conversation, error) {
	return ConversationByConversationID(db, cm.ConversationID)
}

// User returns the User associated with the ConversationMessage's UserID (user_id).
//
// Generated from foreign key 'fk_conversation_message_2'.
func (cm *ConversationMessage) User(db XODB) (*User, error) {
	return UserByUserID(db, cm.UserID)
}

// ConversationMessageByConversationMessageID retrieves a row from 'app_mvp_dating.conversation_message' as a ConversationMessage.
//
// Generated from index 'conversation_message_conversation_message_id_pkey'.
func ConversationMessageByConversationMessageID(db XODB, conversationMessageID uint) (*ConversationMessage, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`conversation_message_id, user_id, conversation_id, message, ip, deleted ` +
		`FROM app_mvp_dating.conversation_message ` +
		`WHERE conversation_message_id = ?`

	// run query
	XOLog(sqlstr, conversationMessageID)
	cm := ConversationMessage{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, conversationMessageID).Scan(&cm.ConversationMessageID, &cm.UserID, &cm.ConversationID, &cm.Message, &cm.IP, &cm.Deleted)
	if err != nil {
		return nil, err
	}

	return &cm, nil
}

// ConversationMessagesByConversationID retrieves a row from 'app_mvp_dating.conversation_message' as a ConversationMessage.
//
// Generated from index 'fk_conversation_message_1_idx'.
func ConversationMessagesByConversationID(db XODB, conversationID uint) ([]*ConversationMessage, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`conversation_message_id, user_id, conversation_id, message, ip, deleted ` +
		`FROM app_mvp_dating.conversation_message ` +
		`WHERE conversation_id = ?`

	// run query
	XOLog(sqlstr, conversationID)
	q, err := db.Query(sqlstr, conversationID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ConversationMessage{}
	for q.Next() {
		cm := ConversationMessage{
			_exists: true,
		}

		// scan
		err = q.Scan(&cm.ConversationMessageID, &cm.UserID, &cm.ConversationID, &cm.Message, &cm.IP, &cm.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &cm)
	}

	return res, nil
}

// ConversationMessagesByUserID retrieves a row from 'app_mvp_dating.conversation_message' as a ConversationMessage.
//
// Generated from index 'fk_conversation_message_2_idx'.
func ConversationMessagesByUserID(db XODB, userID uint) ([]*ConversationMessage, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`conversation_message_id, user_id, conversation_id, message, ip, deleted ` +
		`FROM app_mvp_dating.conversation_message ` +
		`WHERE user_id = ?`

	// run query
	XOLog(sqlstr, userID)
	q, err := db.Query(sqlstr, userID)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	// load results
	res := []*ConversationMessage{}
	for q.Next() {
		cm := ConversationMessage{
			_exists: true,
		}

		// scan
		err = q.Scan(&cm.ConversationMessageID, &cm.UserID, &cm.ConversationID, &cm.Message, &cm.IP, &cm.Deleted)
		if err != nil {
			return nil, err
		}

		res = append(res, &cm)
	}

	return res, nil
}
