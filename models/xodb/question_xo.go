// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"errors"
)

// Question represents a row from 'app_mvp_dating.question'.
type Question struct {
	QuestionID uint         `json:"question_id"` // question_id
	Content    string       `json:"content"`     // content
	Deleted    bool `json:"deleted"`     // deleted

	// xo fields
	_exists bool
}

// Exists determines if the Question exists in the database.
func (q *Question) Exists() bool {
	return q._exists
}

// Deleted provides information if the Question has been deleted from the database.
func (q *Question) IsDeleted() bool {
	return q.Deleted == true
}

// Insert inserts the Question to the database.
func (q *Question) Insert(db XODB) error {
	var err error

	// if already exist, bail
	if q._exists {
		return errors.New("insert failed: already exists")
	}

	// sql insert query, primary key provided by autoincrement
	const sqlstr = `INSERT INTO app_mvp_dating.question (` +
		`content, deleted` +
		`) VALUES (` +
		`?, ?` +
		`)`

	// run query
	XOLog(sqlstr, q.Content, q.Deleted)
	res, err := db.Exec(sqlstr, q.Content, q.Deleted)
	if err != nil {
		return err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return err
	}

	// set primary key and existence
	q.QuestionID = uint(id)
	q._exists = true

	return nil
}

// Update updates the Question in the database.
func (q *Question) Update(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !q._exists {
		return errors.New("update failed: does not exist")
	}

	// if deleted, bail
	if q.Deleted == true {
		return errors.New("update failed: marked for deletion")
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.question SET ` +
		`content = ?, deleted = ?` +
		` WHERE question_id = ?`

	// run query
	XOLog(sqlstr, q.Content, q.Deleted, q.QuestionID)
	_, err = db.Exec(sqlstr, q.Content, q.Deleted, q.QuestionID)
	return err
}

// Save saves the Question to the database.
func (q *Question) Save(db XODB) error {
	if q.Exists() {
		return q.Update(db)
	}

	return q.Insert(db)
}

// Delete deletes the Question from the database.
func (q *Question) Delete(db XODB) error {
	var err error

	// if doesn't exist, bail
	if !q._exists {
		return nil
	}

	// if deleted, bail
	if q.Deleted == true {
		return nil
	}

	// sql query
	const sqlstr = `UPDATE app_mvp_dating.question SET deleted = ? WHERE question_id = ?`

	// run query
	XOLog(sqlstr, true, q.QuestionID)
	_, err = db.Exec(sqlstr, true, q.QuestionID)
	if err != nil {
		return err
	}

	// set deleted
	q.Deleted = true

	return nil
}

// QuestionByContent retrieves a row from 'app_mvp_dating.question' as a Question.
//
// Generated from index 'content_UNIQUE'.
func QuestionByContent(db XODB, content string) (*Question, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`question_id, content, deleted ` +
		`FROM app_mvp_dating.question ` +
		`WHERE content = ?`

	// run query
	XOLog(sqlstr, content)
	qVal := Question{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, content).Scan(&qVal.QuestionID, &qVal.Content, &qVal.Deleted)
	if err != nil {
		return nil, err
	}

	return &qVal, nil
}

// QuestionByQuestionID retrieves a row from 'app_mvp_dating.question' as a Question.
//
// Generated from index 'question_question_id_pkey'.
func QuestionByQuestionID(db XODB, questionID uint) (*Question, error) {
	var err error

	// sql query
	const sqlstr = `SELECT ` +
		`question_id, content, deleted ` +
		`FROM app_mvp_dating.question ` +
		`WHERE question_id = ?`

	// run query
	XOLog(sqlstr, questionID)
	qVal := Question{
		_exists: true,
	}

	err = db.QueryRow(sqlstr, questionID).Scan(&qVal.QuestionID, &qVal.Content, &qVal.Deleted)
	if err != nil {
		return nil, err
	}

	return &qVal, nil
}
