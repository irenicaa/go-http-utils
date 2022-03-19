package httputils

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strconv"

	"github.com/irenicaa/go-http-utils/models"
)

// ...
var (
	IDPattern   = regexp.MustCompile(`/\d+`)
	DatePattern = regexp.MustCompile(`/\d{4}-\d{2}-\d{2}`)
)

// ErrKeyIsMissed ...
var ErrKeyIsMissed = errors.New("key is missed")

// GetIDFromURL ...
func GetIDFromURL(request *http.Request) (int, error) {
	idAsStr := IDPattern.FindString(request.URL.Path)
	if idAsStr == "" {
		return 0, errors.New("unable to find an ID")
	}

	id, err := strconv.Atoi(idAsStr[1:])
	if err != nil {
		return 0, fmt.Errorf("unable to parse the ID: %s", err)
	}

	return id, nil
}

// GetDateFromURL ...
func GetDateFromURL(request *http.Request) (models.Date, error) {
	dateAsStr := DatePattern.FindString(request.URL.Path)
	if dateAsStr == "" {
		return models.Date{}, errors.New("unable to find a date")
	}

	date, err := models.ParseDate(dateAsStr[1:])
	if err != nil {
		return models.Date{}, fmt.Errorf("unable to parse the date: %s", err)
	}

	return date, nil
}

// GetIntFormValue ...
func GetIntFormValue(
	request *http.Request,
	key string,
	min int,
	max int,
) (int, error) {
	value := request.FormValue(key)
	if value == "" {
		return 0, ErrKeyIsMissed
	}

	valueAsInt, err := strconv.Atoi(value)
	if err != nil {
		return 0, fmt.Errorf("value is incorrect: %v", err)
	}
	if valueAsInt < min {
		return 0, errors.New("value too less")
	}
	if valueAsInt > max {
		return 0, errors.New("value too greater")
	}

	return valueAsInt, nil
}

// GetDateFormValue ...
func GetDateFormValue(request *http.Request, key string) (models.Date, error) {
	value := request.FormValue(key)
	if value == "" {
		return models.Date{}, ErrKeyIsMissed
	}

	parsedDate, err := models.ParseDate(value)
	if err != nil {
		return models.Date{}, fmt.Errorf("unable to parse the date: %v", err)
	}

	return parsedDate, nil
}
