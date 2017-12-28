// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql/driver"
	"errors"
)

// Status is the 'status' enum type from schema 'app_mvp_dating'.
type NullStatus struct {
	Status uint8
	Valid  bool
}

var (
	StatusNil = NullStatus{0, false}

	// StatusSingle is the 'single' Status.
	StatusSingle = NullStatus{1, true}

	// StatusDivorce is the 'divorce' Status.
	StatusDivorce = NullStatus{2, true}

	// StatusComplicate is the 'complicate' Status.
	StatusComplicate = NullStatus{3, true}
)

// String returns the string value of the Status.
func (s NullStatus) String() string {
	var enumVal string

	switch s {
	case StatusNil:
		enumVal = ""

	case StatusSingle:
		enumVal = "single"

	case StatusDivorce:
		enumVal = "divorce"

	case StatusComplicate:
		enumVal = "complicate"
	}

	return enumVal
}

// MarshalText marshals Status into text.
func (s NullStatus) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

// UnmarshalText unmarshals Status from text.
func (s *NullStatus) UnmarshalText(text []byte) error {
	switch string(text) {
	case "":
		*s = StatusNil

	case "single":
		*s = StatusSingle

	case "divorce":
		*s = StatusDivorce

	case "complicate":
		*s = StatusComplicate

	default:
		return errors.New("invalid status")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Status.
func (s NullStatus) Value() (driver.Value, error) {
	if !s.Valid {
		return nil, nil
	}
	return s.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Status.
func (s *NullStatus) Scan(src interface{}) error {
	if src == nil {
		*s = StatusNil
		return nil
	}
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid status")
	}

	return s.UnmarshalText(buf)
}
