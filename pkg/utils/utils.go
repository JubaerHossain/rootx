package utils

import (
	"encoding/json"
	"time"
)

type CustomDateFormat time.Time

// MarshalJSON implements the json.Marshaler interface
func (c CustomDateFormat) MarshalJSON() ([]byte, error) {
	t := time.Time(c)
	formatted := t.Format("2006-01-02") // Format the time as "YYYY-MM-DD"
	return json.Marshal(formatted)
}
