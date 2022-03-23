package httputils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// Logger ...
type Logger interface {
	Print(arguments ...interface{})
}

// HandleJSON ...
func HandleJSON(writer http.ResponseWriter, logger Logger, data interface{}) {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		status, message :=
			http.StatusInternalServerError, "unable to marshal the data: %s"
		HandleError(writer, logger, status, message, err)

		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.Write(dataBytes)
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
