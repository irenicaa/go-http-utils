package middlewares

import (
	"net/http"
)

// CORSMiddleware ...
func CORSMiddleware(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		writer.Header().Set(
			"Access-Control-Allow-Origin",
			request.Header.Get("Origin"),
		)
		writer.Header().Set(
			"Access-Control-Allow-Methods",
			request.Header.Get("Access-Control-Request-Method"),
		)
		writer.Header().Set(
			"Access-Control-Allow-Headers",
			request.Header.Get("Access-Control-Request-Headers"),
		)
		if request.Method == http.MethodOptions {
			return
		}

		handler.ServeHTTP(writer, request)
	})
}
