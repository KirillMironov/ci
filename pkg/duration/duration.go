package duration

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"time"
)

type Duration time.Duration

func (d *Duration) Duration() time.Duration {
	return time.Duration(*d)
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("%q", d.Duration())), nil
}

func (d *Duration) UnmarshalJSON(b []byte) error {
	duration, err := time.ParseDuration(string(b)[1 : len(b)-1])
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}

func (d Duration) Value() (driver.Value, error) {
	return d.Duration().String(), nil
}

func (d *Duration) Scan(src any) error {
	str, ok := src.(string)
	if !ok {
		return errors.New("incompatible source type")
	}

	duration, err := time.ParseDuration(str)
	if err != nil {
		return err
	}

	*d = Duration(duration)
	return nil
}
