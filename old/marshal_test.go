package old

import (
	"errors"
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
		want    any
		wantErr error
	}{
		{
			name: "Empty TestStruct",
			args: args{
				v: TestStruct{},
			},
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{}}",
		},
		{
			name: "Empty TestStructPtrs",
			args: args{
				v: TestStructPtrs{},
			},
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStructPtrs\",\"v\":{}}",
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
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"bool\":true,\"byte\":13,\"float32\":11.1,\"float64\":12.2,\"int\":1,\"int16\":3,\"int32\":4,\"int64\":5,\"int8\":2,\"string\":\"test\",\"uint\":6,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"uint8\":7}}",
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
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStructPtrs\",\"v\":{\"bool\":true,\"byte\":13,\"float32\":11.1,\"float64\":12.2,\"int\":1,\"int16\":3,\"int32\":4,\"int64\":5,\"int8\":2,\"string\":\"test\",\"uint\":6,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"uint8\":7}}",
		},
		{
			name: "Empty TestStruct Pointer",
			args: args{
				v: &TestStruct{},
			},
			want: "{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{}}",
		},
		{
			name: "Empty TestStructPtrs Pointer",
			args: args{
				v: &TestStructPtrs{},
			},
			want: "{\"_t\":\"*github.com/trojanc/jsonr.TestStructPtrs\",\"v\":{}}",
		},
		{
			name: "Pointer to slice of empty TestStruct",
			args: args{
				v: &[]TestStruct{},
			},
			want: "{\"_t\":\"*[]github.com/trojanc/jsonr.TestStruct\",\"v\":[]}",
		},
		{
			name: "Slice of empty TestStruct",
			args: args{
				v: []TestStruct{},
			},
			want: &Object{
				Type:  "[]github.com/trojanc/jsonr.TestStruct",
				Value: "[]",
			},
		},
		{
			name: "Slice of TestStruct",
			args: args{
				v: []TestStruct{
					{
						String: "a",
					},
					{
						String: "b",
					},
				},
			},
			want: &Object{
				Type:  "[]github.com/trojanc/jsonr.TestStruct",
				Value: "[{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"a\\\"}\"},{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"b\\\"}\"}]",
			},
		},
		{
			name: "Slice of TestStruct pointers",
			args: args{
				v: []*TestStruct{
					{
						String: "a",
					},
					{
						String: "b",
					},
				},
			},
			want: &Object{
				Type:  "[]*github.com/trojanc/jsonr.TestStruct",
				Value: "[{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"a\\\"}\"},{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"b\\\"}\"}]",
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
			want: &Object{
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
			want: &Object{
				Type:  "map[string]*github.com/trojanc/jsonr.TestStruct",
				Value: "{\"foo\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"string1\\\"}\"},\"john\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"int\\\":1}\"}}",
			},
		},
		{
			name: "Map of string to any with TestStruct",
			args: args{
				v: map[string]any{
					"foo":  TestStruct{String: "string1"},
					"john": TestStruct{Int: 1},
				},
			},
			want: &Object{
				Type:  "map[string]any",
				Value: "{\"foo\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"string1\\\"}\"},\"john\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"int\\\":1}\"}}",
			},
		},
		{
			name: "Map of string to any with pointer to TestStruct",
			args: args{
				v: map[string]any{
					"foo":  &TestStruct{String: "string1"},
					"john": &TestStruct{Int: 1},
				},
			},
			want: &Object{
				Type:  "map[string]any",
				Value: "{\"foo\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"string\\\":\\\"string1\\\"}\"},\"john\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":\"{\\\"int\\\":1}\"}}",
			},
		},
		{
			name: "Map of Struct to pointer TestStruct",
			args: args{
				v: map[TestStruct]*TestStruct{
					TestStruct{String: "a"}: {String: "string1"},
					TestStruct{String: "b"}: {Int: 1},
				},
			},
			want:    nil,
			wantErr: errors.New("map keys cannot be structs"),
		},
		{
			name: "Map of string to slice of TestStruct",
			args: args{
				v: map[string][]TestStruct{
					"a": {{String: "string1"}},
					"b": {{Int: 1}},
				},
			},
			want: &Object{
				Type:  "map[string][]github.com/trojanc/jsonr.TestStruct",
				Value: "{\"a\":{\"_t\":\"[]github.com/trojanc/jsonr.TestStruct\",\"v\":\"[{\\\"_t\\\":\\\"github.com/trojanc/jsonr.TestStruct\\\",\\\"v\\\":\\\"{\\\\\\\"string\\\\\\\":\\\\\\\"string1\\\\\\\"}\\\"}]\"},\"b\":{\"_t\":\"[]github.com/trojanc/jsonr.TestStruct\",\"v\":\"[{\\\"_t\\\":\\\"github.com/trojanc/jsonr.TestStruct\\\",\\\"v\\\":\\\"{\\\\\\\"int\\\\\\\":1}\\\"}]\"}}",
			},
		},
		{
			name: "Map of string to pointer to slice of TestStruct",
			args: args{
				v: map[string]*[]TestStruct{
					"a": {{String: "string1"}},
					"b": {{Int: 1}},
				},
			},
			want: &Object{
				Type:  "map[string]*[]github.com/trojanc/jsonr.TestStruct",
				Value: "{\"a\":{\"_t\":\"*[]github.com/trojanc/jsonr.TestStruct\",\"v\":\"[{\\\"_t\\\":\\\"github.com/trojanc/jsonr.TestStruct\\\",\\\"v\\\":\\\"{\\\\\\\"string\\\\\\\":\\\\\\\"string1\\\\\\\"}\\\"}]\"},\"b\":{\"_t\":\"*[]github.com/trojanc/jsonr.TestStruct\",\"v\":\"[{\\\"_t\\\":\\\"github.com/trojanc/jsonr.TestStruct\\\",\\\"v\\\":\\\"{\\\\\\\"int\\\\\\\":1}\\\"}]\"}}",
			},
		},
		{
			name: "map[string]map[string]string",
			args: args{
				v: map[string]map[string]string{
					"1": {"a": "b"},
					"2": {"c": "d"},
				},
			},
			want: &Object{
				Type:  "map[string]map[string]string",
				Value: "{\"1\":{\"_t\":\"map[string]string\",\"v\":\"{\\\"a\\\":\\\"b\\\"}\"},\"2\":{\"_t\":\"map[string]string\",\"v\":\"{\\\"c\\\":\\\"d\\\"}\"}}",
			},
		},
		{
			name: "Slice of map[string]",
			args: args{
				v: []map[string]string{
					{"1": "a"},
					{"2": "b"},
				},
			},
			want: &Object{
				Type:  "[]map[string]string",
				Value: "[{\"_t\":\"map[string]string\",\"v\":\"{\\\"1\\\":\\\"a\\\"}\"},{\"_t\":\"map[string]string\",\"v\":\"{\\\"2\\\":\\\"b\\\"}\"}]",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			assert.Equal(t, tt.wantErr, err)

			if got != nil {
				fmt.Println(string(got))
			}
			value := string(got)
			assert.Equal(t, tt.want, value)

			obj, err := Unmarshal(got, WithUnmarshalComplexTypes(
				RegisterType(TestStruct{}),
				RegisterType(TestStructPtrs{}),
			))
			fmt.Println(obj)
			assert.Equal(t, tt.args.v, obj)
		})
	}
}

func Test(t *testing.T) {
	test := map[string][]TestStruct{
		"a": {{String: "string1"}},
		"b": {{Int: 1}},
	}
	val, err := Marshal(test, WithMarshalComplexTypes())
	assert.NoError(t, err)
	fmt.Println(string(val))

}
