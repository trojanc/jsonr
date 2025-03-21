package jsonr

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

// Test struct with all primary go types as fields
type TestStruct struct {
	String  string  `json:"string,omitempty"`
	Int     int     `json:"int,omitempty"`
	Int8    int8    `json:"int8,omitempty"`
	Int16   int16   `json:"int16,omitempty"`
	Int32   int32   `json:"int32,omitempty"`
	Int64   int64   `json:"int64,omitempty"`
	Uint    uint    `json:"uint,omitempty"`
	Uint8   uint8   `json:"uint8,omitempty"`
	Uint16  uint16  `json:"uint16,omitempty"`
	Uint32  uint32  `json:"uint32,omitempty"`
	Uint64  uint64  `json:"uint64,omitempty"`
	Float32 float32 `json:"float32,omitempty"`
	Float64 float64 `json:"float64,omitempty"`
	Bool    bool    `json:"bool,omitempty"`
	Byte    byte    `json:"byte,omitempty"`
}
type TestStructPtrs struct {
	String  *string  `json:"string,omitempty"`
	Int     *int     `json:"int,omitempty"`
	Int8    *int8    `json:"int8,omitempty"`
	Int16   *int16   `json:"int16,omitempty"`
	Int32   *int32   `json:"int32,omitempty"`
	Int64   *int64   `json:"int64,omitempty"`
	Uint    *uint    `json:"uint,omitempty"`
	Uint8   *uint8   `json:"uint8,omitempty"`
	Uint16  *uint16  `json:"uint16,omitempty"`
	Uint32  *uint32  `json:"uint32,omitempty"`
	Uint64  *uint64  `json:"uint64,omitempty"`
	Float32 *float32 `json:"float32,omitempty"`
	Float64 *float64 `json:"float64,omitempty"`
	Bool    *bool    `json:"bool,omitempty"`
	Byte    *byte    `json:"byte,omitempty"`
}

func ptr[T any](v T) *T {
	return &v
}

func Test_newJSONRStruct(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name    string
		args    args
		want    *jsonrStruct
		wantErr bool
	}{
		{
			name: "Empty TestStruct",
			args: args{
				v: TestStruct{},
			},
			want: &jsonrStruct{
				Type:  "github.com/trojanc/jsonr.TestStruct",
				Value: "{}",
			},
		},
		{
			name: "Empty TestStructPtrs",
			args: args{
				v: TestStructPtrs{},
			},
			want: &jsonrStruct{
				Type:  "github.com/trojanc/jsonr.TestStructPtrs",
				Value: "{}",
			},
		},
		{
			name: "Fully populated TestStruct",
			args: args{
				v: TestStruct{
					String:  "test",
					Int:     1,
					Int8:    2,
					Int16:   3,
					Int32:   4,
					Int64:   5,
					Uint:    6,
					Uint8:   7,
					Uint16:  8,
					Uint32:  9,
					Uint64:  10,
					Float32: 11.1,
					Float64: 12.2,
					Bool:    true,
					Byte:    13,
				},
			},
			want: &jsonrStruct{
				Type:  "github.com/trojanc/jsonr.TestStruct",
				Value: "{\"string\":\"test\",\"int\":1,\"int8\":2,\"int16\":3,\"int32\":4,\"int64\":5,\"uint\":6,\"uint8\":7,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"float32\":11.1,\"float64\":12.2,\"bool\":true,\"byte\":13}",
			},
		},
		{
			name: "Fully populated TestStructPtr",
			args: args{
				v: TestStructPtrs{
					String:  ptr("test"),
					Int:     ptr(1),
					Int8:    ptr[int8](2),
					Int16:   ptr[int16](3),
					Int32:   ptr[int32](4),
					Int64:   ptr[int64](5),
					Uint:    ptr[uint](6),
					Uint8:   ptr[uint8](7),
					Uint16:  ptr[uint16](8),
					Uint32:  ptr[uint32](9),
					Uint64:  ptr[uint64](10),
					Float32: ptr[float32](11.1),
					Float64: ptr(12.2),
					Bool:    ptr(true),
					Byte:    ptr[byte](13),
				},
			},
			want: &jsonrStruct{
				Type:  "github.com/trojanc/jsonr.TestStructPtrs",
				Value: "{\"string\":\"test\",\"int\":1,\"int8\":2,\"int16\":3,\"int32\":4,\"int64\":5,\"uint\":6,\"uint8\":7,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"float32\":11.1,\"float64\":12.2,\"bool\":true,\"byte\":13}",
			},
		},
		{
			name: "Empty TestStruct Pointer",
			args: args{
				v: &TestStruct{},
			},
			want: &jsonrStruct{
				Type:  "*github.com/trojanc/jsonr.TestStruct",
				Value: "{}",
			},
		},
		{
			name: "Empty TestStructPtrs Pointer",
			args: args{
				v: &TestStructPtrs{},
			},
			want: &jsonrStruct{
				Type:  "*github.com/trojanc/jsonr.TestStructPtrs",
				Value: "{}",
			},
		},
		{
			name: "Pointer to slice of empty TestStruct",
			args: args{
				v: &[]TestStruct{},
			},
			want: &jsonrStruct{
				Type:  "*[]github.com/trojanc/jsonr.TestStruct",
				Value: "[]",
			},
		},
		{
			name: "Slice of empty TestStruct",
			args: args{
				v: []TestStruct{},
			},
			want: &jsonrStruct{
				Type:  "[]github.com/trojanc/jsonr.TestStruct",
				Value: "[]",
			},
		},
		{
			name: "Map of string to TestStruct",
			args: args{
				v: map[string]TestStruct{
					"foo":  {String: "string1"},
					"john": {Int: 1},
				},
			},
			want: &jsonrStruct{
				Type:  "map[string]github.com/trojanc/jsonr.TestStruct",
				Value: "{\"foo\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"string1\\\"}\"},\"john\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"int\\\":1}\"}}",
			},
		},
		{
			name: "Map of string to pointer TestStruct",
			args: args{
				v: map[string]*TestStruct{
					"foo":  {String: "string1"},
					"john": {Int: 1},
				},
			},
			want: &jsonrStruct{
				Type:  "map[string]*github.com/trojanc/jsonr.TestStruct",
				Value: "{\"foo\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"string1\\\"}\"},\"john\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"int\\\":1}\"}}",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := newJSONRStruct(tt.args.v)
			if (err != nil) != tt.wantErr {
				t.Errorf("newJSONRStruct() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			fmt.Printf("type=\"%s\", value=\"%s\"\n", got.Type, got.Value)
			assert.Equal(t, tt.want.Type, got.Type)
			assert.Equal(t, tt.want.Value, got.Value)
		})
	}
}
