package httputils

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"

	"github.com/irenicaa/go-http-utils/models"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
}

// ErrKeyIsMissed ...
var ErrKeyIsMissed = errors.New("key is missed")

// ...
var (
	IDPattern   = regexp.MustCompile(`/\d+`)
	DatePattern = regexp.MustCompile(`/\d{4}-\d{2}-\d{2}`)
)

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

// GetJSONData ...
func GetJSONData(reader io.Reader, data interface{}) error {
	dataAsJSON, err := ioutil.ReadAll(reader)
	if err != nil {
		return fmt.Errorf("unable to read the JSON data: %s", err)
	}

	if err := json.Unmarshal(dataAsJSON, data); err != nil {
		return fmt.Errorf("unable to unmarshal the JSON data: %s", err)
	}

	return nil
}

// HandleError ...
func HandleError(
	writer http.ResponseWriter,
	logger Logger,
	status int,
	format string,
	arguments ...interface{},
) {
	message := fmt.Sprintf(format, arguments...)
	logger.Print(message)

	writer.WriteHeader(status)
	writer.Write([]byte(message))
}

// HandleJSON ...
func HandleJSON(writer http.ResponseWriter, logger Logger, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		status, message :=
			http.StatusInternalServerError, "unable to marshal the data: %v"
		HandleError(writer, logger, status, message, err)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(dataBytes)
}
