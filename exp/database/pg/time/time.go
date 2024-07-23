package time

import (
	"database/sql/driver"
	"fmt"
	"time"
)

const TimeFormat = "15:04:05"

// Time type.
type Time time.Time

func NewTime(hour, min, sec int) Time {
	v := time.Date(0, time.January, 1, hour, min, sec, 0, time.UTC)
	return Time(v)
}

// Scan implementation.
func (v *Time) Scan(src interface{}) error {
	switch src := src.(type) {
	case []byte:
		return v.UnmarshalText(string(src))
	case string:
		return v.UnmarshalText(src)
	case time.Time:
		*v = Time(src)
	case nil:
		*v = Time{}
	default:
		return fmt.Errorf("cannot sql.Scan() Time from: %#v", v)
	}
	return nil
}

// Value implementation.
func (v Time) Value() (driver.Value, error) {
	return driver.Value(time.Time(v).Format(TimeFormat)), nil
}

// UnmarshalText parses time from the text
func (v *Time) UnmarshalText(value string) error {
	dd, err := time.Parse(TimeFormat, value)
	if err != nil {
		return err
	}
	*v = Time(dd)
	return nil
}
