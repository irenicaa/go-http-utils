package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParseDate(t *testing.T) {
	type args struct {
		data string
	}

	tests := []struct {
		name    string
		args    args
		want    Date
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name: "success with a non-zero time",
			args: args{
				data: "2006-01-02",
			},
			want:    Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			wantErr: assert.NoError,
		},
		{
			name: "success with a zero time",
			args: args{
				data: "0001-01-01",
			},
			want:    Date(time.Time{}),
			wantErr: assert.NoError,
		},
		{
			name: "error",
			args: args{
				data: "incorrect",
			},
			want:    Date(time.Time{}),
			wantErr: assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDate(tt.args.data)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}

func TestDate_UnmarshalJSON(t *testing.T) {
	type args struct {
		data []byte
	}

	tests := []struct {
		name     string
		args     args
		wantDate Date
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success with a null",
			args: args{
				data: []byte("null"),
			},
			wantDate: Date(time.Time{}),
			wantErr:  assert.NoError,
		},
		{
			name: "success with a non-zero time",
			args: args{
				data: []byte(`"2006-01-02"`),
			},
			wantDate: Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			wantErr:  assert.NoError,
		},
		{
			name: "success with a zero time",
			args: args{
				data: []byte(`"0001-01-01"`),
			},
			wantDate: Date(time.Time{}),
			wantErr:  assert.NoError,
		},
		{
			name: "error with unmarshalling",
			args: args{
				data: []byte("23"),
			},
			wantDate: Date(time.Time{}),
			wantErr:  assert.Error,
		},
		{
			name: "error with parsing",
			args: args{
				data: []byte(`"incorrect"`),
			},
			wantDate: Date(time.Time{}),
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var date Date
			err := json.Unmarshal(tt.args.data, &date)

			assert.Equal(t, tt.wantDate, date)
			tt.wantErr(t, err)
		})
	}
}

func TestDate_MarshalJSON(t *testing.T) {
	tests := []struct {
		name    string
		date    Date
		want    []byte
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "with a non-zero time",
			date:    Date(time.Date(2006, time.January, 2, 0, 0, 0, 0, time.UTC)),
			want:    []byte(`"2006-01-02"`),
			wantErr: assert.NoError,
		},
		{
			name:    "with a zero time",
			date:    Date(time.Time{}),
			want:    []byte(`"0001-01-01"`),
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := json.Marshal(tt.date)

			assert.Equal(t, tt.want, got)
			tt.wantErr(t, err)
		})
	}
}
