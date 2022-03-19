package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCORSMiddleware(t *testing.T) {
	type args struct {
		handler http.Handler
		request *http.Request
	}

	tests := []struct {
		name         string
		args         args
		wantResponse *http.Response
	}{
		{
			name: "without the CORS headers",
			args: args{
				handler: http.HandlerFunc(func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					writer.Write([]byte("Hello, world!"))
				}),
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/test",
					nil,
				),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Type":                 []string{"text/plain; charset=utf-8"},
					"Access-Control-Allow-Origin":  []string{""},
					"Access-Control-Allow-Methods": []string{""},
					"Access-Control-Allow-Headers": []string{""},
				},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("Hello, world!"))),
				ContentLength: -1,
			},
		},
		{
			name: "with the CORS headers",
			args: args{
				handler: http.HandlerFunc(func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					writer.Write([]byte("Hello, world!"))
				}),
				request: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodGet,
						"http://example.com/test",
						nil,
					)
					request.Header.Set("Origin", "http://example.com")
					request.Header.Set("Access-Control-Request-Method", http.MethodPost)
					request.Header.Set("Access-Control-Request-Headers", "Content-Type")

					return request
				}(),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Content-Type":                 []string{"text/plain; charset=utf-8"},
					"Access-Control-Allow-Origin":  []string{"http://example.com"},
					"Access-Control-Allow-Methods": []string{http.MethodPost},
					"Access-Control-Allow-Headers": []string{"Content-Type"},
				},
				Body:          ioutil.NopCloser(bytes.NewReader([]byte("Hello, world!"))),
				ContentLength: -1,
			},
		},
		{
			name: "with the OPTIONS HTTP method",
			args: args{
				handler: http.HandlerFunc(func(
					writer http.ResponseWriter,
					request *http.Request,
				) {
					writer.Write([]byte("Hello, world!"))
				}),
				request: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodOptions,
						"http://example.com/test",
						nil,
					)
					request.Header.Set("Origin", "http://example.com")
					request.Header.Set("Access-Control-Request-Method", http.MethodPost)
					request.Header.Set("Access-Control-Request-Headers", "Content-Type")

					return request
				}(),
			},
			wantResponse: &http.Response{
				Status: strconv.Itoa(http.StatusOK) + " " +
					http.StatusText(http.StatusOK),
				StatusCode: http.StatusOK,
				Proto:      "HTTP/1.1",
				ProtoMajor: 1,
				ProtoMinor: 1,
				Header: http.Header{
					"Access-Control-Allow-Origin":  []string{"http://example.com"},
					"Access-Control-Allow-Methods": []string{http.MethodPost},
					"Access-Control-Allow-Headers": []string{"Content-Type"},
				},
				Body:          ioutil.NopCloser(bytes.NewReader(nil)),
				ContentLength: -1,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			responseRecorder := httptest.NewRecorder()
			wrappedHandler := CORSMiddleware(tt.args.handler)
			wrappedHandler.ServeHTTP(responseRecorder, tt.args.request)

			assert.Equal(t, tt.wantResponse, responseRecorder.Result())
		})
	}
}
