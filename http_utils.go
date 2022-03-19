package httputils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
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
