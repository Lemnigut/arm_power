package service

import "time"

func parseClientTime(value string) (time.Time, error) {
	if parsed, err := time.Parse(time.RFC3339Nano, value); err == nil {
		return parsed, nil
	}

	layouts := []string{
		"2006-01-02T15:04:05.999999999",
		"2006-01-02T15:04:05.999999",
		"2006-01-02T15:04:05.999",
		"2006-01-02T15:04:05",
	}
	for _, layout := range layouts {
		if parsed, err := time.ParseInLocation(layout, value, time.Local); err == nil {
			return parsed, nil
		}
	}

	return time.Time{}, &time.ParseError{
		Layout: time.RFC3339Nano,
		Value:  value,
	}
}
