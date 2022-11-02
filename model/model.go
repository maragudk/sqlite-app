package model

import (
	"database/sql/driver"
	"time"

	"github.com/maragudk/errors"
)

type Article struct {
	ID      int
	Title   string
	Content string
	Created Time
	Updated Time
}

type Time struct {
	T time.Time
}

// Value satisfies driver.Valuer interface.
func (t *Time) Value() (driver.Value, error) {
	return t.T.Format(time.RFC3339Nano), nil
}

// Scan satisfies sql.Scanner interface.
func (t *Time) Scan(src any) error {
	if src == nil {
		return nil
	}

	s, ok := src.(string)
	if !ok {
		return errors.Newf("error scanning time, got %+v", src)
	}

	parsedT, err := time.Parse(time.RFC3339Nano, s)
	if err != nil {
		return err
	}

	t.T = parsedT

	return nil
}
