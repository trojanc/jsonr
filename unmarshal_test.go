package jsonr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func Test_deserializeObject(t *testing.T) {
	type args struct {
		data       string
		targetType string
		opts       *unmarshalOptions
	}
	tests := []struct {
		name    string
		args    args
		want    interface{}
		wantErr error
	}{
		{
			name: "TestStruct",
			args: args{
				data:       "{}",
				targetType: "github.com/trojanc/jsonr.TestStruct",
				opts: &unmarshalOptions{
					typeMapping: variableTypeMapping{
						"github.com/trojanc/jsonr.TestStruct": reflect.TypeOf(TestStruct{}),
					},
				},
			},
			want: TestStruct{},
		},
		{
			name: "*TestStruct",
			args: args{
				data:       "{}",
				targetType: "*github.com/trojanc/jsonr.TestStruct",
				opts: &unmarshalOptions{
					typeMapping: variableTypeMapping{
						"github.com/trojanc/jsonr.TestStruct": reflect.TypeOf(TestStruct{}),
					},
				},
			},
			want: &TestStruct{},
		},
		{
			name: "[]TestStruct",
			args: args{
				data:       "[]",
				targetType: "[]github.com/trojanc/jsonr.TestStruct",
				opts: &unmarshalOptions{
					typeMapping: variableTypeMapping{
						"github.com/trojanc/jsonr.TestStruct": reflect.TypeOf(TestStruct{}),
					},
				},
			},
			want: []TestStruct{},
		},
		{
			name: "[2]TestStruct",
			args: args{
				data:       "[{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"a\\\"}\"},{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"b\\\"}\"}]",
				targetType: "[]github.com/trojanc/jsonr.TestStruct",
				opts: &unmarshalOptions{
					typeMapping: variableTypeMapping{
						"github.com/trojanc/jsonr.TestStruct": reflect.TypeOf(TestStruct{}),
					},
				},
			},
			want: []TestStruct{
				{
					String: "a",
				},
				{
					String: "b",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := deserializeObject(tt.args.data, tt.args.targetType, tt.args.opts)
			assert.Equal(t, tt.wantErr, err)

			if got != nil {
				if o, ok := got.(*Object); ok {
					fmt.Printf("type=  \"%s\", \nvalue= \"%s\"\n", o.Type, o.Value)
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
