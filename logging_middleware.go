package httputils

import (
	"fmt"
	"net/http"
	"time"
)

// LoggingMiddleware ...
func LoggingMiddleware(
	handler http.Handler,
	logger Logger,
	clock func() time.Time,
) http.Handler {
	return http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		startTime := clock()
		handler.ServeHTTP(writer, request)

		elapsedTime := clock().Sub(startTime)
		message := fmt.Sprintf("%s %s %s", request.Method, request.URL, elapsedTime)
		logger.Print(message)
	})
}
