package httputils

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadJSONData(t *testing.T) {
	type args struct {
		httpClient   HTTPClient
		url          string
		authHeader   string
		responseData interface{}
	}
	type testData struct {
		FieldOne int
		FieldTwo string
	}

	tests := []struct {
		name             string
		args             args
		wantResponseData interface{}
		wantErr          assert.ErrorAssertionFunc
	}{
		{
			name: "success without authorization",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(bytes.NewReader([]byte(
							`{"FieldOne": 23, "FieldTwo": "test"}`,
						))),
					}

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.On("Do", request).Return(response, nil).Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{FieldOne: 23, FieldTwo: "test"},
			wantErr:          assert.NoError,
		},
		{
			name: "success with authorization",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)
					request.Header.Add("Authorization", "Bearer token")

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(bytes.NewReader([]byte(
							`{"FieldOne": 23, "FieldTwo": "test"}`,
						))),
					}

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.On("Do", request).Return(response, nil).Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "Bearer token",
				responseData: &testData{},
			},
			wantResponseData: &testData{FieldOne: 23, FieldTwo: "test"},
			wantErr:          assert.NoError,
		},
		{
			name: "error with request creating",
			args: args{
				httpClient:   &MockHTTPClient{},
				url:          ":",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{},
			wantErr:          assert.Error,
		},
		{
			name: "error with request sending",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.
						On("Do", request).
						Return((*http.Response)(nil), iotest.ErrTimeout).
						Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{},
			wantErr:          assert.Error,
		},
		{
			name: "error with the reading of the response body",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body: ioutil.NopCloser(iotest.TimeoutReader(bytes.NewReader([]byte(
							`{"FieldOne": 23, "FieldTwo": "test"}`,
						)))),
					}

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.On("Do", request).Return(response, nil).Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{},
			wantErr:          assert.Error,
		},
		{
			name: "error with the response status",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)

					response := &http.Response{
						StatusCode: http.StatusInternalServerError,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("error"))),
					}

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.On("Do", request).Return(response, nil).Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{},
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.EqualError(t, err, "request was failed: 500 error")
			},
		},
		{
			name: "error with the unmarshalling of the response body",
			args: args{
				httpClient: func() HTTPClient {
					request, err := http.NewRequest(http.MethodGet, "http://example.com/", nil)
					require.NoError(t, err)

					response := &http.Response{
						StatusCode: http.StatusOK,
						Body:       ioutil.NopCloser(bytes.NewReader([]byte("incorrect"))),
					}

					httpClient := &MockHTTPClient{}
					httpClient.InnerMock.On("Do", request).Return(response, nil).Times(1)

					return httpClient
				}(),
				url:          "http://example.com/",
				authHeader:   "",
				responseData: &testData{},
			},
			wantResponseData: &testData{},
			wantErr:          assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := LoadJSONData(
				tt.args.httpClient,
				tt.args.url,
				tt.args.authHeader,
				tt.args.responseData,
			)

			tt.args.httpClient.(*MockHTTPClient).InnerMock.AssertExpectations(t)
			assert.Equal(t, tt.wantResponseData, tt.args.responseData)
			tt.wantErr(t, err)
		})
	}
}
