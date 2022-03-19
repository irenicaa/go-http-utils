package httputils

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/irenicaa/go-http-utils/models"
	"github.com/stretchr/testify/assert"
)

func TestGetIDFromURL(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/23",
					nil,
				),
			},
			want:    23,
			wantErr: assert.NoError,
		},
		{
			name: "error on finding",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name: "error on parsing",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/99999999999999999999999999",
					nil,
				),
			},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetIDFromURL(tt.args.request)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestGetDateFromURL(t *testing.T) {
	type args struct {
		request *http.Request
	}

	tests := []struct {
		name    string
		args    args
		want    models.Date
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/2006-01-02",
					nil,
				),
			},
			want:    models.Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			wantErr: assert.NoError,
		},
		{
			name: "error on finding",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos",
					nil,
				),
			},
			want:    models.Date{},
			wantErr: assert.Error,
		},
		{
			name: "error on parsing",
			args: args{
				request: httptest.NewRequest(
					http.MethodGet,
					"http://example.com/api/v1/todos/9999-99-99",
					nil,
				),
			},
			want:    models.Date{},
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDateFromURL(tt.args.request)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestGetIntFormValue(t *testing.T) {
	type args struct {
		request *http.Request
		key     string
		min     int
		max     int
	}

	tests := []struct {
		name    string
		args    args
		want    int
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=23", nil),
				key:     "key",
				min:     0,
				max:     100,
			},
			want:    23,
			wantErr: assert.NoError,
		},
		{
			name: "error with a missed key",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test", nil),
				key:     "key",
				min:     0,
				max:     100,
			},
			want: 0,
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Equal(t, ErrKeyIsMissed, err, msgAndArgs...)
			},
		},
		{
			name: "error with an incorrect key",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=value", nil),
				key:     "key",
				min:     0,
				max:     100,
			},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name: "error with a too less value",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=23", nil),
				key:     "key",
				min:     50,
				max:     100,
			},
			want:    0,
			wantErr: assert.Error,
		},
		{
			name: "error with a too greater value",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=23", nil),
				key:     "key",
				min:     0,
				max:     10,
			},
			want:    0,
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err :=
				GetIntFormValue(tt.args.request, tt.args.key, tt.args.min, tt.args.max)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestGetDateFormValue(t *testing.T) {
	type args struct {
		request *http.Request
		key     string
	}

	tests := []struct {
		name    string
		args    args
		want    models.Date
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=2006-01-02", nil),
				key:     "key",
			},
			want:    models.Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			wantErr: assert.NoError,
		},
		{
			name: "error with a missed key",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test", nil),
				key:     "key",
			},
			want: models.Date(time.Time{}),
			wantErr: func(t assert.TestingT, err error, msgAndArgs ...interface{}) bool {
				return assert.Equal(t, ErrKeyIsMissed, err, msgAndArgs...)
			},
		},
		{
			name: "error with an incorrect key",
			args: args{
				request: httptest.NewRequest(http.MethodGet, "/test?key=value", nil),
				key:     "key",
			},
			want:    models.Date(time.Time{}),
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetDateFormValue(tt.args.request, tt.args.key)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
