package nursys

import (
	"errors"
	"regexp"
	"time"
)

type Time time.Time

// MarshalJSON implements the json.Marshaler interface.
func (t Time) MarshalJSON() ([]byte, error) {
	return time.Time(t).MarshalJSON()
}

var timeRegexp = regexp.MustCompile(`"^\\d{4}-\\d{2}-\\d{2}T\\d{2}%3A\\d{2}%3A\\d{2}(?:%2E\\d+)?[A-Z]?(?:[+.-](?:08%3A\\d{2}|\\d{2}[A-Z]))?$"`)

// UnmarshalJSON implements the json.Unmarshaler interface.
func (t *Time) UnmarshalJSON(data []byte) (err error) {
	var tm time.Time
	s := string(data)
	if data[0] != '"' || data[len(data)-1] != '"' {
		return errors.New("invalid JSON string")
	}
	s = s[1 : len(s)-1]
	// Try RFC3339 first
	if err = tm.UnmarshalJSON(data); err == nil {
		*t = Time(tm)
		return nil
	}

	// Try to parse the string as non-strict RFC3339 (e.g. 2021-01-02 03:04:05Z)
	if tm, err = time.Parse("2006-01-02 15:04:05Z07:00", s); err == nil {
		*t = Time(tm)
		return nil
	}

	// Try to parse the string as a custom format without time zone (e.g. 2021-01-02T03:04:05)
	if tm, err = time.Parse("2006-01-02T15:04:05", s); err == nil {
		*t = Time(tm)
		return nil
	}

	// Try to parse the string as a custom format with time zone (e.g. 2021-01-02T03:04:05-07:00)
	if tm, err = time.Parse("2006-01-02T15:04:05-07:00", s); err == nil {
		*t = Time(tm)
		return nil
	}

	return errors.New("invalid time format")
}
