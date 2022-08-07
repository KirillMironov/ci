package duration

import (
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
