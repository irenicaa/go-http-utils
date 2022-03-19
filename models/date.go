package models

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

const dateFormat = "2006-01-02"

// Date ...
type Date time.Time

// ParseDate ...
func ParseDate(data string) (Date, error) {
	parsedDate, err := time.Parse(dateFormat, data)
	if err != nil {
		return Date{}, err
	}

	return Date(parsedDate), nil
}

// UnmarshalJSON ...
func (date *Date) UnmarshalJSON(data []byte) error {
	if bytes.Equal(data, []byte("null")) {
		return nil
	}

	var formattedDate string
	if err := json.Unmarshal(data, &formattedDate); err != nil {
		return fmt.Errorf("unable to unmarshal the string: %v", err)
	}

	parsedDate, err := ParseDate(formattedDate)
	if err != nil {
		return fmt.Errorf("unable to parse the date: %v", err)
	}

	*date = parsedDate
	return nil
}

// MarshalJSON ...
func (date Date) MarshalJSON() ([]byte, error) {
	formattedDate := time.Time(date).Format(dateFormat)
	return json.Marshal(formattedDate)
}
