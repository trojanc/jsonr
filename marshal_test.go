package jsonr

import (
	"encoding/json"
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

func Test_MarshalAndUnmarshal(t *testing.T) {
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
			name: "nil",
			args: args{},
			want: "null",
		},
		{
			name: "Empty TestStruct",
			args: args{
				v: TestStruct{},
			},
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{}}",
		},
		{
			name: "int",
			args: args{
				v: int(1),
			},
			want: "{\"_t\":\"int\",\"v\":1}",
		},
		// repeat the above for each primitive type
		{
			name: "int8",
			args: args{
				v: int8(2),
			},
			want: "{\"_t\":\"int8\",\"v\":2}",
		},
		{
			name: "int16",
			args: args{
				v: int16(3),
			},
			want: "{\"_t\":\"int16\",\"v\":3}",
		},
		{
			name: "int32",
			args: args{
				v: int32(4),
			},
			want: "{\"_t\":\"int32\",\"v\":4}",
		},
		{
			name: "int64",
			args: args{
				v: int64(5),
			},
			want: "{\"_t\":\"int64\",\"v\":5}",
		},
		{
			name: "uint",
			args: args{
				v: uint(6),
			},
			want: "{\"_t\":\"uint\",\"v\":6}",
		},
		{
			name: "uint8",
			args: args{
				v: uint8(7),
			},
			want: "{\"_t\":\"uint8\",\"v\":7}",
		},
		{
			name: "uint16",
			args: args{
				v: uint16(8),
			},
			want: "{\"_t\":\"uint16\",\"v\":8}",
		},
		{
			name: "uint32",
			args: args{
				v: uint32(9),
			},
			want: "{\"_t\":\"uint32\",\"v\":9}",
		},
		{
			name: "uint64",
			args: args{
				v: uint64(10),
			},
			want: "{\"_t\":\"uint64\",\"v\":10}",
		},
		{
			name: "float32",
			args: args{
				v: float32(11.1),
			},
			want: "{\"_t\":\"float32\",\"v\":11.1}",
		},
		{
			name: "float64",
			args: args{
				v: float64(12.2),
			},
			want: "{\"_t\":\"float64\",\"v\":12.2}",
		},
		{
			name: "bool",
			args: args{
				v: true,
			},
			want: "{\"_t\":\"bool\",\"v\":true}",
		},
		{
			name: "byte",
			args: args{
				v: byte(13),
			},
			want: "{\"_t\":\"uint8\",\"v\":13}",
		},
		{
			name: "string",
			args: args{
				v: "test",
			},
			want: "{\"_t\":\"string\",\"v\":\"test\"}",
		},

		{
			name: "*int",
			args: args{
				v: ptr(int(1)),
			},
			want: "{\"_t\":\"*int\",\"v\":1}",
		},
		// repeat the above for each primitive type
		{
			name: "*int8",
			args: args{
				v: ptr(int8(2)),
			},
			want: "{\"_t\":\"*int8\",\"v\":2}",
		},
		{
			name: "*int16",
			args: args{
				v: ptr(int16(3)),
			},
			want: "{\"_t\":\"*int16\",\"v\":3}",
		},
		{
			name: "*int32",
			args: args{
				v: ptr(int32(4)),
			},
			want: "{\"_t\":\"*int32\",\"v\":4}",
		},
		{
			name: "*int64",
			args: args{
				v: ptr(int64(5)),
			},
			want: "{\"_t\":\"*int64\",\"v\":5}",
		},
		{
			name: "*uint",
			args: args{
				v: ptr(uint(6)),
			},
			want: "{\"_t\":\"*uint\",\"v\":6}",
		},
		{
			name: "*uint8",
			args: args{
				v: ptr(uint8(7)),
			},
			want: "{\"_t\":\"*uint8\",\"v\":7}",
		},
		{
			name: "*uint16",
			args: args{
				v: ptr(uint16(8)),
			},
			want: "{\"_t\":\"*uint16\",\"v\":8}",
		},
		{
			name: "*uint32",
			args: args{
				v: ptr(uint32(9)),
			},
			want: "{\"_t\":\"*uint32\",\"v\":9}",
		},
		{
			name: "*uint64",
			args: args{
				v: ptr(uint64(10)),
			},
			want: "{\"_t\":\"*uint64\",\"v\":10}",
		},
		{
			name: "*float32",
			args: args{
				v: ptr(float32(11.1)),
			},
			want: "{\"_t\":\"*float32\",\"v\":11.1}",
		},
		{
			name: "*float64",
			args: args{
				v: ptr(float64(12.2)),
			},
			want: "{\"_t\":\"*float64\",\"v\":12.2}",
		},
		{
			name: "*bool",
			args: args{
				v: ptr(true),
			},
			want: "{\"_t\":\"*bool\",\"v\":true}",
		},
		{
			name: "*byte",
			args: args{
				v: ptr(byte(13)),
			},
			want: "{\"_t\":\"*uint8\",\"v\":13}",
		},
		{
			name: "*string",
			args: args{
				v: ptr("test"),
			},
			want: "{\"_t\":\"*string\",\"v\":\"test\"}",
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
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"test\",\"int\":1,\"int8\":2,\"int16\":3,\"int32\":4,\"int64\":5,\"uint\":6,\"uint8\":7,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"float32\":11.1,\"float64\":12.2,\"bool\":true,\"byte\":13}}",
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
			want: "{\"_t\":\"github.com/trojanc/jsonr.TestStructPtrs\",\"v\":{\"string\":\"test\",\"int\":1,\"int8\":2,\"int16\":3,\"int32\":4,\"int64\":5,\"uint\":6,\"uint8\":7,\"uint16\":8,\"uint32\":9,\"uint64\":10,\"float32\":11.1,\"float64\":12.2,\"bool\":true,\"byte\":13}}",
		},
		{
			name: "*TestStruct empty",
			args: args{
				v: &TestStruct{},
			},
			want: "{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{}}",
		},
		{
			name: "*TestStructPtrs empty",
			args: args{
				v: &TestStructPtrs{},
			},
			want: "{\"_t\":\"*github.com/trojanc/jsonr.TestStructPtrs\",\"v\":{}}",
		},
		{
			name: "*[]TestStruct empty",
			args: args{
				v: &[]TestStruct{},
			},
			want: "{\"_t\":\"*[]github.com/trojanc/jsonr.TestStruct\",\"v\":[]}",
		},
		{
			name: "[]TestStruct empty",
			args: args{
				v: []TestStruct{},
			},
			want: "{\"_t\":\"[]github.com/trojanc/jsonr.TestStruct\",\"v\":[]}",
		},
		{
			name: "[]TestStruct",
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
			want: "{\"_t\":\"[]github.com/trojanc/jsonr.TestStruct\",\"v\":[{\"string\":\"a\"},{\"string\":\"b\"}]}",
		},
		{
			name: "[]any with TestStruct",
			args: args{
				v: []any{
					TestStruct{
						String: "a",
					},
					TestStruct{
						String: "b",
					},
				},
			},
			want: "{\"_t\":\"[]interface\",\"v\":[{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"a\"}},{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"b\"}}]}",
		},
		{
			name: "[]any with *TestStruct",
			args: args{
				v: []any{
					&TestStruct{
						String: "a",
					},
					&TestStruct{
						String: "b",
					},
				},
			},
			want: "{\"_t\":\"[]interface\",\"v\":[{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"a\"}},{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"b\"}}]}",
		},
		{
			name: "[]any with *TestStruct",
			args: args{
				v: []any{
					&TestStruct{
						String: "a",
					},
					TestStruct{
						String: "b",
					},
				},
			},
			want: "{\"_t\":\"[]interface\",\"v\":[{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"a\"}},{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"b\"}}]}",
		},
		{
			name: "[]*TestStruct",
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
			want: "{\"_t\":\"[]*github.com/trojanc/jsonr.TestStruct\",\"v\":[{\"string\":\"a\"},{\"string\":\"b\"}]}",
		},
		{
			name: "*[]*TestStruct",
			args: args{
				v: &[]*TestStruct{
					{
						String: "a",
					},
					{
						String: "b",
					},
				},
			},
			want: "{\"_t\":\"*[]*github.com/trojanc/jsonr.TestStruct\",\"v\":[{\"string\":\"a\"},{\"string\":\"b\"}]}",
		},
		{
			name: "map[string]TestStruct",
			args: args{
				v: map[string]TestStruct{
					"foo":  {String: "string1"},
					"john": {Int: 1},
				},
			},
			want: "{\"_t\":\"map[string]github.com/trojanc/jsonr.TestStruct\",\"v\":{\"foo\":{\"string\":\"string1\"},\"john\":{\"int\":1}}}",
		},
		{
			name: "map[string]*TestStruct",
			args: args{
				v: map[string]*TestStruct{
					"foo":  {String: "string1"},
					"john": {Int: 1},
				},
			},
			want: "{\"_t\":\"map[string]*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"foo\":{\"string\":\"string1\"},\"john\":{\"int\":1}}}",
		},
		{
			name: "Error map[*string]*TestStruct",
			args: args{
				v: map[*string]*TestStruct{
					ptr("foo"):  {String: "string1"},
					ptr("john"): {Int: 1},
				},
			},
			wantErr: errors.New("unsupported map key"),
		},
		{
			name: "map[string]any with TestStruct",
			args: args{
				v: map[string]any{
					"foo":  TestStruct{String: "string1"},
					"john": TestStruct{Int: 1},
				},
			},
			want: "{\"_t\":\"map[string]interface\",\"v\":{\"foo\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"string1\"}},\"john\":{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"int\":1}}}}",
		},
		{
			name: "map[string]any with *TestStruct",
			args: args{
				v: map[string]any{
					"foo":  &TestStruct{String: "string1"},
					"john": &TestStruct{Int: 1},
				},
			},
			want: "{\"_t\":\"map[string]interface\",\"v\":{\"foo\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"string1\"}},\"john\":{\"_t\":\"*github.com/trojanc/jsonr.TestStruct\",\"v\":{\"int\":1}}}}",
		},
		{
			name: "Error map[TestStruct]*TestStruct",
			args: args{
				v: map[TestStruct]*TestStruct{
					TestStruct{String: "a"}: {String: "string1"},
					TestStruct{String: "b"}: {Int: 1},
				},
			},
			want:    nil,
			wantErr: errors.New("unsupported map key"),
		},
		{
			name: "map[string][]TestStruct",
			args: args{
				v: map[string][]TestStruct{
					"a": {{String: "string1"}},
					"b": {{Int: 1}},
				},
			},
			want: "{\"_t\":\"map[string][]github.com/trojanc/jsonr.TestStruct\",\"v\":{\"a\":[{\"string\":\"string1\"}],\"b\":[{\"int\":1}]}}",
		},
		{ // Even though marshalling to pointers i
			name: "map[string]map[*string]string",
			args: args{
				v: map[string]map[*string]string{
					"key": {
						ptr("key1"): "value",
					},
				},
			},
			wantErr: errors.New("failed to unmarshal: json: unsupported type: map[*string]string"),
		},
		{
			name: "map[string]*[]TestStruct",
			args: args{
				v: map[string]*[]TestStruct{
					"a": {{String: "string1"}},
					"b": {{Int: 1}},
				},
			},
			want: "{\"_t\":\"map[string]*[]github.com/trojanc/jsonr.TestStruct\",\"v\":{\"a\":[{\"string\":\"string1\"}],\"b\":[{\"int\":1}]}}",
		},
		{
			name: "map[string]map[string]string",
			args: args{
				v: map[string]map[string]string{
					"1": {"a": "b"},
					"2": {"c": "d"},
				},
			},
			want: "{\"_t\":\"map[string]map[string]string\",\"v\":{\"1\":{\"a\":\"b\"},\"2\":{\"c\":\"d\"}}}",
		},
		{
			name: "[]]map[string]string",
			args: args{
				v: []map[string]string{
					{"1": "a"},
					{"2": "b"},
				},
			},
			want: "{\"_t\":\"[]map[string]string\",\"v\":[{\"1\":\"a\"},{\"2\":\"b\"}]}",
		},
		{
			name: "map[string]any with nils",
			args: args{
				v: map[string]any{
					"1": nil,
					"2": nil,
				},
			},
			want: "{\"_t\":\"map[string]interface\",\"v\":{\"1\":null,\"2\":null}}",
		},
		{
			name: "[]any with nils",
			args: args{
				v: []any{
					nil,
					"2",
				},
			},
			want: "{\"_t\":\"[]interface\",\"v\":[null,{\"_t\":\"string\",\"v\":\"2\"}]}",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Marshal(tt.args.v)
			assert.Equal(t, tt.wantErr, err)

			if err == nil {
				if got != nil {
					fmt.Println(string(got))
				}
				value := string(got)
				assert.Equal(t, tt.want, value)

				obj, err := Unmarshal(got,
					RegisterType(TestStruct{}),
					RegisterType(TestStructPtrs{}),
				)
				assert.NoError(t, err)
				fmt.Println(obj)
				assert.Equal(t, tt.args.v, obj)
			}
		})
	}
}

func TestCompareGoJSON(t *testing.T) {
	obj := TestStruct{
		String: "tester",
	}

	goStr, err := json.Marshal(obj)
	assert.NoError(t, err)
	assert.Equal(t, "{\"string\":\"tester\"}", fmt.Sprintf("%v", string(goStr)))

	// Try using native go json.Unmarshal() using type any
	var data any
	err = json.Unmarshal(goStr, data)
	assert.Error(t, err) // Expect error, json.Unmarshal() cannot read into any

	// Try using native go json.Unmarshal() using type map[string]any
	mapData := make(map[string]any)
	err = json.Unmarshal(goStr, &mapData)
	assert.NoError(t, err) // Expect error, json.Unmarshal() cannot read into any
	assert.Equal(t, "map[string:tester]", fmt.Sprintf("%v", mapData))

	// Now lets use jsonr
	test, err := Marshal(obj)
	assert.NoError(t, err)
	assert.Equal(t, "{\"_t\":\"github.com/trojanc/jsonr.TestStruct\",\"v\":{\"string\":\"tester\"}}", string(test))

	data, err = Unmarshal(test, RegisterType(TestStruct{}))
	assert.NoError(t, err)
	assert.Equal(t, "{tester 0 0 0 0 0 0 0 0 0 0 0 0 false 0}", fmt.Sprintf("%v", data))
	assert.Equal(t, obj, data)
}

func ptr[T any](v T) *T {
	return &v
}
