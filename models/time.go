package models

import (
	"strings"
	"time"
)

const format = "2006-01-02"

type CivilTime time.Time

func (c *CivilTime) UnmarshalJSON(b []byte) error {
	value := strings.Trim(string(b), `"`)
	if value == "" || value == "null" {
		return nil
	}

	t, err := time.Parse(format, value)
	if err != nil {
		return err
	}

	*c = CivilTime(t)
	return nil
}

func (c CivilTime) MarshalJSON() ([]byte, error) {
	return []byte(`"` + time.Time(c).Format(format) + `"`), nil
}
