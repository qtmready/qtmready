package defs

import (
	"encoding/json"
	"strings"
	"time"
)

type (
	// Timestamp is hack around github's funky use of timestamps.
	Timestamp time.Time
)

func (t *Timestamp) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	switch v := raw.(type) {
	case float64:
		*t = Timestamp(time.Unix(int64(v), 0))
	case string:
		if strings.HasSuffix(v, "Z") {
			t_, err := time.Parse("2006-01-02T15:04:05Z", v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		} else {
			t_, err := time.Parse(time.RFC3339, v)
			if err != nil {
				return err
			}

			*t = Timestamp(t_)
		}
	}

	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	t_ := time.Time(t)
	return json.Marshal(t_.Format(time.RFC3339))
}

func (t Timestamp) Time() time.Time {
	return time.Time(t)
}
