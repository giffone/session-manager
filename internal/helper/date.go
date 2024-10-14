package helper

import (
	"errors"
	"time"
)

func ValidateFromTo(from, to string) (time.Time, time.Time, error) {
	fromT := time.Unix(0, 0)
	toT := time.Unix(0, 0)

	if from == "" && to == "" {
		fromT = time.Now().Truncate(24 * time.Hour)
		toT = fromT.Add(24 * time.Hour)
		return fromT, toT, nil
	}

	t, err := ParseDate(from)
	if err != nil {
		return fromT, toT, err
	}
	fromT = t.Truncate(24 * time.Hour)

	if to == "" {
		toT = time.Now().Truncate(24 * time.Hour).Add(24 * time.Hour)
	} else {
		t, err := ParseDate(to)
		if err != nil {
			return fromT, toT, err
		}
		toT = t.Truncate(24 * time.Hour)
	}

	return fromT, toT, nil

}

func ParseDate(s string) (t time.Time, err error) {
	t, err = time.Parse(time.RFC3339, s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.DateTime, s)
	if err == nil {
		return t, nil
	}
	t, err = time.Parse(time.DateOnly, s)
	if err == nil {
		return t, nil
	}

	return t, errors.New(`date format must be '2006-01-02T15:04:05Z07:00' or '2006-01-02 15:04:05' or '2006-01-02'`)
}
