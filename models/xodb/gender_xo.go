// Package xodb contains the types for schema 'app_mvp_dating'.
package xodb

// Code generated by xo. DO NOT EDIT.

import (
	"database/sql/driver"
	"errors"
)

// Gender is the 'gender' enum type from schema 'app_mvp_dating'.
type Gender uint16

const (
	// GenderX is the 'X' Gender.
	GenderX = Gender(1)

	// GenderM is the 'M' Gender.
	GenderM = Gender(2)

	// GenderF is the 'F' Gender.
	GenderF = Gender(3)
)

// String returns the string value of the Gender.
func (g Gender) String() string {
	var enumVal string

	switch g {
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
func (g Gender) MarshalText() ([]byte, error) {
	return []byte(g.String()), nil
}

// UnmarshalText unmarshals Gender from text.
func (g *Gender) UnmarshalText(text []byte) error {
	switch string(text) {
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
func (g Gender) Value() (driver.Value, error) {
	return g.String(), nil
}

// Scan satisfies the database/sql.Scanner interface for Gender.
func (g *Gender) Scan(src interface{}) error {
	buf, ok := src.([]byte)
	if !ok {
		return errors.New("invalid Gender")
	}

	return g.UnmarshalText(buf)
}
