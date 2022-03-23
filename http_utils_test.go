package httputils

import (
	"bytes"
	"io"
	"testing"
	"testing/iotest"

	"github.com/stretchr/testify/assert"
)

func TestGetJSONData(t *testing.T) {
	type testData struct {
		FieldOne int
		FieldTwo string
	}
	type args struct {
		reader io.Reader
		data   interface{}
	}

	tests := []struct {
		name     string
		args     args
		wantData interface{}
		wantErr  assert.ErrorAssertionFunc
	}{
		{
			name: "success",
			args: args{
				reader: bytes.NewReader([]byte(`{"FieldOne": 23, "FieldTwo": "test"}`)),
				data:   &testData{},
			},
			wantData: &testData{FieldOne: 23, FieldTwo: "test"},
			wantErr:  assert.NoError,
		},
		{
			name: "error on reading",
			args: args{
				reader: iotest.TimeoutReader(bytes.NewReader(
					[]byte(`{"FieldOne": 23, "FieldTwo": "test"}`),
				)),
				data: &testData{},
			},
			wantData: &testData{},
			wantErr:  assert.Error,
		},
		{
			name: "error on unmarshalling",
			args: args{
				reader: bytes.NewReader([]byte("incorrect")),
				data:   &testData{},
			},
			wantData: &testData{},
			wantErr:  assert.Error,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := GetJSONData(tt.args.reader, tt.args.data)

			assert.Equal(t, tt.wantData, tt.args.data)
			tt.wantErr(t, err)
		})
	}
}
