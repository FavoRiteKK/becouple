// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql/driver"
	"errors"
)

// Gender is the 'gender' enum type from schema 'app_mvp_dating'.
type NullGender struct {
	Gender uint8
	Valid  bool
}

var (
	// Gender is nil
	GenderNil = NullGender{Gender: 0, Valid: false}

	// GenderX is the 'X' Gender.
	GenderX = NullGender{Gender: 1, Valid: true}

	// GenderM is the 'M' Gender.
	GenderM = NullGender{Gender: 2, Valid: true}

	// GenderF is the 'F' Gender.
	GenderF = NullGender{Gender: 3, Valid: true}
)

// String returns the string value of the Gender.
func (g NullGender) String() string {
	var enumVal string

	switch g {
	case GenderNil:
		enumVal = "nil"

	case GenderX:
		enumVal = "X"

	case GenderM:
		enumVal = "M"

	case GenderF:
		enumVal = "F"
	}

	return enumVal
}

// MarshalText marshals Gender into text.
func (g NullGender) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

// UnmarshalText unmarshals Gender from text.
func (g *NullGender) UnmarshalText(text []byte) error {
	switch string(text) {
	case "nil":
		*g = GenderNil

	case "X":
		*g = GenderX

	case "M":
		*g = GenderM

	case "F":
		*g = GenderF

	default:
		return errors.New("invalid Gender")
	}

	return nil
}

// Value satisfies the sql/driver.Valuer interface for Gender.
func (g NullGender) Value() (driver.Value, error) {
	if !g.Valid {
		return nil, nil
	}
	return g.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Gender.
func (g *NullGender) Scan(src interface{}) error {
	if src == nil {
		*g = GenderNil
		return nil
	}
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Gender")
	}

	return g.UnmarshalText(buf)
}
