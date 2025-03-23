package jsonr

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestUnmarshal(t *testing.T) {
	type args struct {
		data    []byte
		options []UnmarshalOption
	}
	tests := []struct {
		name    string
		args    args
		want    any
		wantErr assert.ErrorAssertionFunc
		errStr  string
	}{
		{
			name: "Option giving error",
			args: args{
				options: []UnmarshalOption{
					func(*unmarshalOptions) error {
						return errors.New("oops")
					},
				},
			},
			wantErr: assert.Error,
			errStr:  "could not apply option: oops",
		},
		{
			name: "Register invalid type",
			args: args{
				options: []UnmarshalOption{
					RegisterType(func() {}),
				},
			},
			wantErr: assert.Error,
			errStr:  "could not apply option: only instance of structs should be used",
		},
		{
			name: "Broken data",
			args: args{
				data: []byte(`{`),
			},
			wantErr: assert.Error,
			errStr:  "unexpected end of JSON input",
		},
		{
			name: "Broken slice",
			args: args{
				data: []byte(`{"_t":"[]string","v":{}}`),
			},
			wantErr: assert.Error,
			errStr:  "error unmarshalling slice: json: cannot unmarshal object into Go value of type []string",
		},
		{
			name: "Broken slice value",
			args: args{
				data: []byte(`{"_t":"[]interface","v":[{"_t":"string", "v":234}]}`),
			},
			wantErr: assert.Error,
			errStr:  "json: cannot unmarshal number into Go value of type string",
		},
		{
			name: "Broken map",
			args: args{
				data: []byte(`{"_t":"map[string]string","v":[]}`),
			},
			wantErr: assert.Error,
			errStr:  "error unmarshalling map: json: cannot unmarshal array into Go value of type map[string]string",
		},
		{
			name: "Broken map value",
			args: args{
				data: []byte(`{"_t":"map[string]interface","v":{"a":{"_t":"string", "v":234}}}`),
			},
			wantErr: assert.Error,
			errStr:  "json: cannot unmarshal number into Go value of type string",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Unmarshal(tt.args.data, tt.args.options...)
			tt.wantErr(t, err)
			if len(tt.errStr) > 0 {
				assert.Equal(t, tt.errStr, err.Error())
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
