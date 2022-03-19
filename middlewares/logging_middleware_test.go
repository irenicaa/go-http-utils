package middlewares

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestLoggingMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(
		writer http.ResponseWriter,
		request *http.Request,
	) {
		writer.Write([]byte("Hello, world!"))
	})

	logger := &MockLogger{}
	logger.InnerMock.
		On("Print", []interface{}{"GET http://example.com/test 2m3s"}).
		Return().
		Times(1)

	clockCount := 0
	clock := func() time.Time {
		timestamp := time.Date(2021, time.January, 15, 4, 16, 50, 1, time.UTC)
		if clockCount > 0 {
			timestamp = timestamp.Add(123 * time.Second)
		}

		clockCount++
		return timestamp
	}

	responseRecorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "http://example.com/test", nil)

	wrappedHandler := LoggingMiddleware(handler, logger, clock)
	wrappedHandler.ServeHTTP(responseRecorder, request)

	wantResponse := &http.Response{
		Status: strconv.Itoa(http.StatusOK) + " " +
			http.StatusText(http.StatusOK),
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.1",
		ProtoMajor: 1,
		ProtoMinor: 1,
		Header: http.Header{
			"Content-Type": []string{"text/plain; charset=utf-8"},
		},
		Body:          ioutil.NopCloser(bytes.NewReader([]byte("Hello, world!"))),
		ContentLength: -1,
	}

	logger.InnerMock.AssertExpectations(t)
	assert.Equal(t, wantResponse, responseRecorder.Result())
}
