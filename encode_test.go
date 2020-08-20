package json_test

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"math"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/goccy/go-json"
)

type recursiveT struct {
	A *recursiveT `json:"a,omitempty"`
	B *recursiveU `json:"b,omitempty"`
	C *recursiveU `json:"c,omitempty"`
	D string      `json:"d,omitempty"`
}

type recursiveU struct {
	T *recursiveT `json:"t,omitempty"`
}

func Test_Marshal(t *testing.T) {
	t.Run("int", func(t *testing.T) {
		bytes, err := json.Marshal(-10)
		assertErr(t, err)
		assertEq(t, "int", `-10`, string(bytes))
	})
	t.Run("int8", func(t *testing.T) {
		bytes, err := json.Marshal(int8(-11))
		assertErr(t, err)
		assertEq(t, "int8", `-11`, string(bytes))
	})
	t.Run("int16", func(t *testing.T) {
		bytes, err := json.Marshal(int16(-12))
		assertErr(t, err)
		assertEq(t, "int16", `-12`, string(bytes))
	})
	t.Run("int32", func(t *testing.T) {
		bytes, err := json.Marshal(int32(-13))
		assertErr(t, err)
		assertEq(t, "int32", `-13`, string(bytes))
	})
	t.Run("int64", func(t *testing.T) {
		bytes, err := json.Marshal(int64(-14))
		assertErr(t, err)
		assertEq(t, "int64", `-14`, string(bytes))
	})
	t.Run("uint", func(t *testing.T) {
		bytes, err := json.Marshal(uint(10))
		assertErr(t, err)
		assertEq(t, "uint", `10`, string(bytes))
	})
	t.Run("uint8", func(t *testing.T) {
		bytes, err := json.Marshal(uint8(11))
		assertErr(t, err)
		assertEq(t, "uint8", `11`, string(bytes))
	})
	t.Run("uint16", func(t *testing.T) {
		bytes, err := json.Marshal(uint16(12))
		assertErr(t, err)
		assertEq(t, "uint16", `12`, string(bytes))
	})
	t.Run("uint32", func(t *testing.T) {
		bytes, err := json.Marshal(uint32(13))
		assertErr(t, err)
		assertEq(t, "uint32", `13`, string(bytes))
	})
	t.Run("uint64", func(t *testing.T) {
		bytes, err := json.Marshal(uint64(14))
		assertErr(t, err)
		assertEq(t, "uint64", `14`, string(bytes))
	})
	t.Run("float32", func(t *testing.T) {
		bytes, err := json.Marshal(float32(3.14))
		assertErr(t, err)
		assertEq(t, "float32", `3.14`, string(bytes))
	})
	t.Run("float64", func(t *testing.T) {
		bytes, err := json.Marshal(float64(3.14))
		assertErr(t, err)
		assertEq(t, "float64", `3.14`, string(bytes))
	})
	t.Run("bool", func(t *testing.T) {
		bytes, err := json.Marshal(true)
		assertErr(t, err)
		assertEq(t, "bool", `true`, string(bytes))
	})
	t.Run("string", func(t *testing.T) {
		bytes, err := json.Marshal("hello world")
		assertErr(t, err)
		assertEq(t, "string", `"hello world"`, string(bytes))
	})
	t.Run("struct", func(t *testing.T) {
		bytes, err := json.Marshal(struct {
			A int    `json:"a"`
			B uint   `json:"b"`
			C string `json:"c"`
			D int    `json:"-"`  // ignore field
			a int    `json:"aa"` // private field
		}{
			A: -1,
			B: 1,
			C: "hello world",
		})
		assertErr(t, err)
		assertEq(t, "struct", `{"a":-1,"b":1,"c":"hello world"}`, string(bytes))
		t.Run("null", func(t *testing.T) {
			type T struct {
				A *struct{} `json:"a"`
			}
			var v T
			bytes, err := json.Marshal(&v)
			assertErr(t, err)
			assertEq(t, "struct", `{"a":null}`, string(bytes))
		})
		t.Run("recursive", func(t *testing.T) {
			bytes, err := json.Marshal(recursiveT{
				A: &recursiveT{
					B: &recursiveU{
						T: &recursiveT{
							D: "hello",
						},
					},
					C: &recursiveU{
						T: &recursiveT{
							D: "world",
						},
					},
				},
			})
			assertErr(t, err)
			assertEq(t, "recursive", `{"a":{"b":{"t":{"d":"hello"}},"c":{"t":{"d":"world"}}}}`, string(bytes))
		})
		t.Run("embedded", func(t *testing.T) {
			type T struct {
				A string `json:"a"`
			}
			type U struct {
				*T
				B string `json:"b"`
			}
			type T2 struct {
				A string `json:"a,omitempty"`
			}
			type U2 struct {
				*T2
				B string `json:"b,omitempty"`
			}
			t.Run("exists field", func(t *testing.T) {
				bytes, err := json.Marshal(&U{
					T: &T{
						A: "aaa",
					},
					B: "bbb",
				})
				assertErr(t, err)
				assertEq(t, "embedded", `{"a":"aaa","b":"bbb"}`, string(bytes))
				t.Run("omitempty", func(t *testing.T) {
					bytes, err := json.Marshal(&U2{
						T2: &T2{
							A: "aaa",
						},
						B: "bbb",
					})
					assertErr(t, err)
					assertEq(t, "embedded", `{"a":"aaa","b":"bbb"}`, string(bytes))
				})
			})
			t.Run("none field", func(t *testing.T) {
				bytes, err := json.Marshal(&U{
					B: "bbb",
				})
				assertErr(t, err)
				assertEq(t, "embedded", `{"b":"bbb"}`, string(bytes))
				t.Run("omitempty", func(t *testing.T) {
					bytes, err := json.Marshal(&U2{
						B: "bbb",
					})
					assertErr(t, err)
					assertEq(t, "embedded", `{"b":"bbb"}`, string(bytes))
				})
			})
		})
		t.Run("omitempty", func(t *testing.T) {
			type T struct {
				A int                    `json:",omitempty"`
				B int8                   `json:",omitempty"`
				C int16                  `json:",omitempty"`
				D int32                  `json:",omitempty"`
				E int64                  `json:",omitempty"`
				F uint                   `json:",omitempty"`
				G uint8                  `json:",omitempty"`
				H uint16                 `json:",omitempty"`
				I uint32                 `json:",omitempty"`
				J uint64                 `json:",omitempty"`
				K float32                `json:",omitempty"`
				L float64                `json:",omitempty"`
				O string                 `json:",omitempty"`
				P bool                   `json:",omitempty"`
				Q []int                  `json:",omitempty"`
				R map[string]interface{} `json:",omitempty"`
				S *struct{}              `json:",omitempty"`
				T int                    `json:"t,omitempty"`
			}
			var v T
			v.T = 1
			bytes, err := json.Marshal(&v)
			assertErr(t, err)
			assertEq(t, "struct", `{"t":1}`, string(bytes))
			t.Run("int", func(t *testing.T) {
				var v struct {
					A int `json:"a,omitempty"`
					B int `json:"b"`
				}
				v.B = 1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "int", `{"b":1}`, string(bytes))
			})
			t.Run("int8", func(t *testing.T) {
				var v struct {
					A int  `json:"a,omitempty"`
					B int8 `json:"b"`
				}
				v.B = 1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "int8", `{"b":1}`, string(bytes))
			})
			t.Run("int16", func(t *testing.T) {
				var v struct {
					A int   `json:"a,omitempty"`
					B int16 `json:"b"`
				}
				v.B = 1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "int16", `{"b":1}`, string(bytes))
			})
			t.Run("int32", func(t *testing.T) {
				var v struct {
					A int   `json:"a,omitempty"`
					B int32 `json:"b"`
				}
				v.B = 1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "int32", `{"b":1}`, string(bytes))
			})
			t.Run("int64", func(t *testing.T) {
				var v struct {
					A int   `json:"a,omitempty"`
					B int64 `json:"b"`
				}
				v.B = 1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "int64", `{"b":1}`, string(bytes))
			})
			t.Run("string", func(t *testing.T) {
				var v struct {
					A int    `json:"a,omitempty"`
					B string `json:"b"`
				}
				v.B = "b"
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "string", `{"b":"b"}`, string(bytes))
			})
			t.Run("float32", func(t *testing.T) {
				var v struct {
					A int     `json:"a,omitempty"`
					B float32 `json:"b"`
				}
				v.B = 1.1
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "float32", `{"b":1.1}`, string(bytes))
			})
			t.Run("float64", func(t *testing.T) {
				var v struct {
					A int     `json:"a,omitempty"`
					B float64 `json:"b"`
				}
				v.B = 3.14
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "float64", `{"b":3.14}`, string(bytes))
			})
			t.Run("slice", func(t *testing.T) {
				var v struct {
					A int   `json:"a,omitempty"`
					B []int `json:"b"`
				}
				v.B = []int{1, 2, 3}
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "slice", `{"b":[1,2,3]}`, string(bytes))
			})
			t.Run("array", func(t *testing.T) {
				var v struct {
					A int    `json:"a,omitempty"`
					B [2]int `json:"b"`
				}
				v.B = [2]int{1, 2}
				bytes, err := json.Marshal(&v)
				assertErr(t, err)
				assertEq(t, "array", `{"b":[1,2]}`, string(bytes))
			})
			t.Run("map", func(t *testing.T) {
				v := new(struct {
					A int                    `json:"a,omitempty"`
					B map[string]interface{} `json:"b"`
				})
				v.B = map[string]interface{}{"c": 1}
				bytes, err := json.Marshal(v)
				assertErr(t, err)
				assertEq(t, "array", `{"b":{"c":1}}`, string(bytes))
			})
		})
		t.Run("head_omitempty", func(t *testing.T) {
			type T struct {
				A *struct{} `json:"a,omitempty"`
			}
			var v T
			bytes, err := json.Marshal(&v)
			assertErr(t, err)
			assertEq(t, "struct", `{}`, string(bytes))
		})
		t.Run("pointer_head_omitempty", func(t *testing.T) {
			type V struct{}
			type U struct {
				B *V `json:"b,omitempty"`
			}
			type T struct {
				A *U `json:"a"`
			}
			bytes, err := json.Marshal(&T{A: &U{}})
			assertErr(t, err)
			assertEq(t, "struct", `{"a":{}}`, string(bytes))
		})
		t.Run("head_int_omitempty", func(t *testing.T) {
			type T struct {
				A int `json:"a,omitempty"`
			}
			var v T
			bytes, err := json.Marshal(&v)
			assertErr(t, err)
			assertEq(t, "struct", `{}`, string(bytes))
		})
	})
	t.Run("slice", func(t *testing.T) {
		t.Run("[]int", func(t *testing.T) {
			bytes, err := json.Marshal([]int{1, 2, 3, 4})
			assertErr(t, err)
			assertEq(t, "[]int", `[1,2,3,4]`, string(bytes))
		})
		t.Run("[]interface{}", func(t *testing.T) {
			bytes, err := json.Marshal([]interface{}{1, 2.1, "hello"})
			assertErr(t, err)
			assertEq(t, "[]interface{}", `[1,2.1,"hello"]`, string(bytes))
		})
	})

	t.Run("array", func(t *testing.T) {
		bytes, err := json.Marshal([4]int{1, 2, 3, 4})
		assertErr(t, err)
		assertEq(t, "array", `[1,2,3,4]`, string(bytes))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("map[string]int", func(t *testing.T) {
			bytes, err := json.Marshal(map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			})
			assertErr(t, err)
			assertEq(t, "map", len(`{"a":1,"b":2,"c":3,"d":4}`), len(string(bytes)))
		})
		t.Run("map[string]interface{}", func(t *testing.T) {
			type T struct {
				A int
			}
			v := map[string]interface{}{
				"a": 1,
				"b": 2.1,
				"c": &T{
					A: 10,
				},
				"d": 4,
			}
			bytes, err := json.Marshal(v)
			assertErr(t, err)
			assertEq(t, "map[string]interface{}", len(`{"a":1,"b":2.1,"c":{"A":10},"d":4}`), len(string(bytes)))
		})
	})
}

type marshalJSON struct{}

func (*marshalJSON) MarshalJSON() ([]byte, error) {
	return []byte(`1`), nil
}

func Test_MarshalJSON(t *testing.T) {
	t.Run("*struct", func(t *testing.T) {
		bytes, err := json.Marshal(&marshalJSON{})
		assertErr(t, err)
		assertEq(t, "MarshalJSON", "1", string(bytes))
	})
	t.Run("time", func(t *testing.T) {
		bytes, err := json.Marshal(time.Time{})
		assertErr(t, err)
		assertEq(t, "MarshalJSON", `"0001-01-01T00:00:00Z"`, string(bytes))
	})
}

func Test_MarshalIndent(t *testing.T) {
	prefix := "-"
	indent := "\t"
	t.Run("struct", func(t *testing.T) {
		bytes, err := json.MarshalIndent(struct {
			A int    `json:"a"`
			B uint   `json:"b"`
			C string `json:"c"`
			D int    `json:"-"`  // ignore field
			a int    `json:"aa"` // private field
		}{
			A: -1,
			B: 1,
			C: "hello world",
		}, prefix, indent)
		assertErr(t, err)
		result := "{\n-\t\"a\": -1,\n-\t\"b\": 1,\n-\t\"c\": \"hello world\"\n-}"
		assertEq(t, "struct", result, string(bytes))
	})
	t.Run("slice", func(t *testing.T) {
		t.Run("[]int", func(t *testing.T) {
			bytes, err := json.MarshalIndent([]int{1, 2, 3, 4}, prefix, indent)
			assertErr(t, err)
			result := "[\n-\t1,\n-\t2,\n-\t3,\n-\t4\n-]"
			assertEq(t, "[]int", result, string(bytes))
		})
		t.Run("[]interface{}", func(t *testing.T) {
			bytes, err := json.MarshalIndent([]interface{}{1, 2.1, "hello"}, prefix, indent)
			assertErr(t, err)
			result := "[\n-\t1,\n-\t2.1,\n-\t\"hello\"\n-]"
			assertEq(t, "[]interface{}", result, string(bytes))
		})
	})

	t.Run("array", func(t *testing.T) {
		bytes, err := json.MarshalIndent([4]int{1, 2, 3, 4}, prefix, indent)
		assertErr(t, err)
		result := "[\n-\t1,\n-\t2,\n-\t3,\n-\t4\n-]"
		assertEq(t, "array", result, string(bytes))
	})
	t.Run("map", func(t *testing.T) {
		t.Run("map[string]int", func(t *testing.T) {
			bytes, err := json.MarshalIndent(map[string]int{
				"a": 1,
				"b": 2,
				"c": 3,
				"d": 4,
			}, prefix, indent)
			assertErr(t, err)
			result := "{\n-\t\"a\": 1,\n-\t\"b\": 2,\n-\t\"c\": 3,\n-\t\"d\": 4\n-}"
			assertEq(t, "map", len(result), len(string(bytes)))
		})
		t.Run("map[string]interface{}", func(t *testing.T) {
			type T struct {
				E int
				F int
			}
			v := map[string]interface{}{
				"a": 1,
				"b": 2.1,
				"c": &T{
					E: 10,
					F: 11,
				},
				"d": 4,
			}
			bytes, err := json.MarshalIndent(v, prefix, indent)
			assertErr(t, err)
			result := "{\n-\t\"a\": 1,\n-\t\"b\": 2.1,\n-\t\"c\": {\n-\t\t\"E\": 10,\n-\t\t\"F\": 11\n-\t},\n-\t\"d\": 4\n-}"
			assertEq(t, "map[string]interface{}", len(result), len(string(bytes)))
		})
	})
}

type StringTag struct {
	BoolStr    bool        `json:",string"`
	IntStr     int64       `json:",string"`
	UintptrStr uintptr     `json:",string"`
	StrStr     string      `json:",string"`
	NumberStr  json.Number `json:",string"`
}

func TestRoundtripStringTag(t *testing.T) {
	tests := []struct {
		name string
		in   StringTag
		want string // empty to just test that we roundtrip
	}{
		{
			name: "AllTypes",
			in: StringTag{
				BoolStr:    true,
				IntStr:     42,
				UintptrStr: 44,
				StrStr:     "xzbit",
				NumberStr:  "46",
			},
			want: `{
				"BoolStr": "true",
				"IntStr": "42",
				"UintptrStr": "44",
				"StrStr": "\"xzbit\"",
				"NumberStr": "46"
			}`,
		},
		{
			// See golang.org/issues/38173.
			name: "StringDoubleEscapes",
			in: StringTag{
				StrStr:    "\b\f\n\r\t\"\\",
				NumberStr: "0", // just to satisfy the roundtrip
			},
			want: `{
				"BoolStr": "false",
				"IntStr": "0",
				"UintptrStr": "0",
				"StrStr": "\"\\u0008\\u000c\\n\\r\\t\\\"\\\\\"",
				"NumberStr": "0"
			}`,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Indent with a tab prefix to make the multi-line string
			// literals in the table nicer to read.
			got, err := json.MarshalIndent(&test.in, "\t\t\t", "\t")
			if err != nil {
				t.Fatal(err)
			}
			if got := string(got); got != test.want {
				t.Fatalf(" got: %s\nwant: %s\n", got, test.want)
			}

			// Verify that it round-trips.
			var s2 StringTag
			if err := json.Unmarshal(got, &s2); err != nil {
				t.Fatalf("Decode: %v", err)
			}
			if !reflect.DeepEqual(test.in, s2) {
				t.Fatalf("decode didn't match.\nsource: %#v\nEncoded as:\n%s\ndecode: %#v", test.in, string(got), s2)
			}
		})
	}
}

// byte slices are special even if they're renamed types.
type renamedByte byte
type renamedByteSlice []byte
type renamedRenamedByteSlice []renamedByte

func TestEncodeRenamedByteSlice(t *testing.T) {
	s := renamedByteSlice("abc")
	result, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	expect := `"YWJj"`
	if string(result) != expect {
		t.Errorf(" got %s want %s", result, expect)
	}
	r := renamedRenamedByteSlice("abc")
	result, err = json.Marshal(r)
	if err != nil {
		t.Fatal(err)
	}
	if string(result) != expect {
		t.Errorf(" got %s want %s", result, expect)
	}
}

func TestMarshalRawMessageValue(t *testing.T) {
	type (
		T1 struct {
			M json.RawMessage `json:",omitempty"`
		}
		T2 struct {
			M *json.RawMessage `json:",omitempty"`
		}
	)

	var (
		rawNil   = json.RawMessage(nil)
		rawEmpty = json.RawMessage([]byte{})
		rawText  = json.RawMessage([]byte(`"foo"`))
	)

	tests := []struct {
		in   interface{}
		want string
		ok   bool
	}{
		// Test with nil RawMessage.
		{rawNil, "null", true},
		{&rawNil, "null", true},
		{[]interface{}{rawNil}, "[null]", true},
		{&[]interface{}{rawNil}, "[null]", true},
		{[]interface{}{&rawNil}, "[null]", true},
		{&[]interface{}{&rawNil}, "[null]", true},
		{struct{ M json.RawMessage }{rawNil}, `{"M":null}`, true},
		{&struct{ M json.RawMessage }{rawNil}, `{"M":null}`, true},
		{struct{ M *json.RawMessage }{&rawNil}, `{"M":null}`, true},
		{&struct{ M *json.RawMessage }{&rawNil}, `{"M":null}`, true},
		{map[string]interface{}{"M": rawNil}, `{"M":null}`, true},
		{&map[string]interface{}{"M": rawNil}, `{"M":null}`, true},
		{map[string]interface{}{"M": &rawNil}, `{"M":null}`, true},
		{&map[string]interface{}{"M": &rawNil}, `{"M":null}`, true},
		{T1{rawNil}, "{}", true},
		{T2{&rawNil}, `{"M":null}`, true},
		{&T1{rawNil}, "{}", true},
		{&T2{&rawNil}, `{"M":null}`, true},

		// Test with empty, but non-nil, RawMessage.
		{rawEmpty, "", false},
		{&rawEmpty, "", false},
		{[]interface{}{rawEmpty}, "", false},
		{&[]interface{}{rawEmpty}, "", false},
		{[]interface{}{&rawEmpty}, "", false},
		{&[]interface{}{&rawEmpty}, "", false},
		{struct{ X json.RawMessage }{rawEmpty}, "", false},
		{&struct{ X json.RawMessage }{rawEmpty}, "", false},
		{struct{ X *json.RawMessage }{&rawEmpty}, "", false},
		{&struct{ X *json.RawMessage }{&rawEmpty}, "", false},
		{map[string]interface{}{"nil": rawEmpty}, "", false},
		{&map[string]interface{}{"nil": rawEmpty}, "", false},
		{map[string]interface{}{"nil": &rawEmpty}, "", false},
		{&map[string]interface{}{"nil": &rawEmpty}, "", false},

		{T1{rawEmpty}, "{}", true},
		{T2{&rawEmpty}, "", false},
		{&T1{rawEmpty}, "{}", true},
		{&T2{&rawEmpty}, "", false},

		// Test with RawMessage with some text.
		//
		// The tests below marked with Issue6458 used to generate "ImZvbyI=" instead "foo".
		// This behavior was intentionally changed in Go 1.8.
		// See https://golang.org/issues/14493#issuecomment-255857318
		{rawText, `"foo"`, true}, // Issue6458
		{&rawText, `"foo"`, true},
		{[]interface{}{rawText}, `["foo"]`, true},  // Issue6458
		{&[]interface{}{rawText}, `["foo"]`, true}, // Issue6458
		{[]interface{}{&rawText}, `["foo"]`, true},
		{&[]interface{}{&rawText}, `["foo"]`, true},
		{struct{ M json.RawMessage }{rawText}, `{"M":"foo"}`, true}, // Issue6458
		{&struct{ M json.RawMessage }{rawText}, `{"M":"foo"}`, true},
		{struct{ M *json.RawMessage }{&rawText}, `{"M":"foo"}`, true},
		{&struct{ M *json.RawMessage }{&rawText}, `{"M":"foo"}`, true},
		{map[string]interface{}{"M": rawText}, `{"M":"foo"}`, true},  // Issue6458
		{&map[string]interface{}{"M": rawText}, `{"M":"foo"}`, true}, // Issue6458
		{map[string]interface{}{"M": &rawText}, `{"M":"foo"}`, true},
		{&map[string]interface{}{"M": &rawText}, `{"M":"foo"}`, true},
		{T1{rawText}, `{"M":"foo"}`, true}, // Issue6458
		{T2{&rawText}, `{"M":"foo"}`, true},
		{&T1{rawText}, `{"M":"foo"}`, true},
		{&T2{&rawText}, `{"M":"foo"}`, true},
	}

	for i, tt := range tests {
		b, err := json.Marshal(tt.in)
		if ok := (err == nil); ok != tt.ok {
			if err != nil {
				t.Errorf("test %d, unexpected failure: %v", i, err)
			} else {
				t.Errorf("test %d, unexpected success", i)
			}
		}
		if got := string(b); got != tt.want {
			t.Errorf("test %d, Marshal(%#v) = %q, want %q", i, tt.in, got, tt.want)
		}
	}
}

type marshalerError struct{}

func (*marshalerError) MarshalJSON() ([]byte, error) {
	return nil, errors.New("unexpected error")
}

func Test_MarshalerError(t *testing.T) {
	var v marshalerError
	_, err := json.Marshal(&v)
	expect := `json: error calling MarshalJSON for type *json_test.marshalerError: unexpected error`
	assertEq(t, "marshaler error", expect, fmt.Sprint(err))
}

// Ref has Marshaler and Unmarshaler methods with pointer receiver.
type Ref int

func (*Ref) MarshalJSON() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *Ref) UnmarshalJSON([]byte) error {
	*r = 12
	return nil
}

// Val has Marshaler methods with value receiver.
type Val int

func (Val) MarshalJSON() ([]byte, error) {
	return []byte(`"val"`), nil
}

// RefText has Marshaler and Unmarshaler methods with pointer receiver.
type RefText int

func (*RefText) MarshalText() ([]byte, error) {
	return []byte(`"ref"`), nil
}

func (r *RefText) UnmarshalText([]byte) error {
	*r = 13
	return nil
}

// ValText has Marshaler methods with value receiver.
type ValText int

func (ValText) MarshalText() ([]byte, error) {
	return []byte(`"val"`), nil
}

func TestRefValMarshal(t *testing.T) {
	var s = struct {
		R0 Ref
		R1 *Ref
		R2 RefText
		R3 *RefText
		V0 Val
		V1 *Val
		V2 ValText
		V3 *ValText
	}{
		R0: 12,
		R1: new(Ref),
		R2: 14,
		R3: new(RefText),
		V0: 13,
		V1: new(Val),
		V2: 15,
		V3: new(ValText),
	}
	const want = `{"R0":"ref","R1":"ref","R2":"\"ref\"","R3":"\"ref\"","V0":"val","V1":"val","V2":"\"val\"","V3":"\"val\""}`
	b, err := json.Marshal(&s)
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// C implements Marshaler and returns unescaped JSON.
type C int

func (C) MarshalJSON() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

// CText implements Marshaler and returns unescaped text.
type CText int

func (CText) MarshalText() ([]byte, error) {
	return []byte(`"<&>"`), nil
}

func TestMarshalerEscaping(t *testing.T) {
	var c C
	want := `"\u003c\u0026\u003e"`
	b, err := json.Marshal(c)
	if err != nil {
		t.Fatalf("Marshal(c): %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Marshal(c) = %#q, want %#q", got, want)
	}

	var ct CText
	want = `"\"\u003c\u0026\u003e\""`
	b, err = json.Marshal(ct)
	if err != nil {
		t.Fatalf("Marshal(ct): %v", err)
	}
	if got := string(b); got != want {
		t.Errorf("Marshal(ct) = %#q, want %#q", got, want)
	}
}

type marshalPanic struct{}

func (marshalPanic) MarshalJSON() ([]byte, error) { panic(0xdead) }

func TestMarshalPanic(t *testing.T) {
	defer func() {
		if got := recover(); !reflect.DeepEqual(got, 0xdead) {
			t.Errorf("panic() = (%T)(%v), want 0xdead", got, got)
		}
	}()
	json.Marshal(&marshalPanic{})
	t.Error("Marshal should have panicked")
}

func TestMarshalUncommonFieldNames(t *testing.T) {
	v := struct {
		A0, À, Aβ int
	}{}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatal("Marshal:", err)
	}
	want := `{"A0":0,"À":0,"Aβ":0}`
	got := string(b)
	if got != want {
		t.Fatalf("Marshal: got %s want %s", got, want)
	}
}

func TestMarshalerError(t *testing.T) {
	s := "test variable"
	st := reflect.TypeOf(s)
	errText := "json: test error"

	tests := []struct {
		err  *json.MarshalerError
		want string
	}{
		{
			json.NewMarshalerError(st, fmt.Errorf(errText), ""),
			"json: error calling MarshalJSON for type " + st.String() + ": " + errText,
		},
		{
			json.NewMarshalerError(st, fmt.Errorf(errText), "TestMarshalerError"),
			"json: error calling TestMarshalerError for type " + st.String() + ": " + errText,
		},
	}

	for i, tt := range tests {
		got := tt.err.Error()
		if got != tt.want {
			t.Errorf("MarshalerError test %d, got: %s, want: %s", i, got, tt.want)
		}
	}
}

type unmarshalerText struct {
	A, B string
}

// needed for re-marshaling tests
func (u unmarshalerText) MarshalText() ([]byte, error) {
	return []byte(u.A + ":" + u.B), nil
}

func (u *unmarshalerText) UnmarshalText(b []byte) error {
	pos := bytes.IndexByte(b, ':')
	if pos == -1 {
		return errors.New("missing separator")
	}
	u.A, u.B = string(b[:pos]), string(b[pos+1:])
	return nil
}

func TestTextMarshalerMapKeysAreSorted(t *testing.T) {
	b, err := json.Marshal(map[unmarshalerText]int{
		{"x", "y"}: 1,
		{"y", "x"}: 2,
		{"a", "z"}: 3,
		{"z", "a"}: 4,
	})
	if err != nil {
		t.Fatalf("Failed to Marshal text.Marshaler: %v", err)
	}
	const want = `{"a:z":3,"x:y":1,"y:x":2,"z:a":4}`
	if len(string(b)) != len(want) {
		t.Errorf("Marshal map with text.Marshaler keys: got %#q, want %#q", b, want)
	}
}

// https://golang.org/issue/33675
func TestNilMarshalerTextMapKey(t *testing.T) {
	v := map[*unmarshalerText]int{
		(*unmarshalerText)(nil): 1,
		{"A", "B"}:              2,
	}
	b, err := json.Marshal(v)
	if err != nil {
		t.Fatalf("Failed to Marshal *text.Marshaler: %v", err)
	}
	const want = `{"":1,"A:B":2}`
	if len(string(b)) != len(want) {
		t.Errorf("Marshal map with *text.Marshaler keys: got %#q, want %#q", b, want)
	}
}

var re = regexp.MustCompile

// syntactic checks on form of marshaled floating point numbers.
var badFloatREs = []*regexp.Regexp{
	re(`p`),                     // no binary exponential notation
	re(`^\+`),                   // no leading + sign
	re(`^-?0[^.]`),              // no unnecessary leading zeros
	re(`^-?\.`),                 // leading zero required before decimal point
	re(`\.(e|$)`),               // no trailing decimal
	re(`\.[0-9]+0(e|$)`),        // no trailing zero in fraction
	re(`^-?(0|[0-9]{2,})\..*e`), // exponential notation must have normalized mantissa
	re(`e[0-9]`),                // positive exponent must be signed
	//re(`e[+-]0`),                // exponent must not have leading zeros
	re(`e-[1-6]$`),             // not tiny enough for exponential notation
	re(`e+(.|1.|20)$`),         // not big enough for exponential notation
	re(`^-?0\.0000000`),        // too tiny, should use exponential notation
	re(`^-?[0-9]{22}`),         // too big, should use exponential notation
	re(`[1-9][0-9]{16}[1-9]`),  // too many significant digits in integer
	re(`[1-9][0-9.]{17}[1-9]`), // too many significant digits in decimal
	// below here for float32 only
	re(`[1-9][0-9]{8}[1-9]`),  // too many significant digits in integer
	re(`[1-9][0-9.]{9}[1-9]`), // too many significant digits in decimal
}

func TestMarshalFloat(t *testing.T) {
	t.Parallel()
	nfail := 0
	test := func(f float64, bits int) {
		vf := interface{}(f)
		if bits == 32 {
			f = float64(float32(f)) // round
			vf = float32(f)
		}
		bout, err := json.Marshal(vf)
		if err != nil {
			t.Errorf("Marshal(%T(%g)): %v", vf, vf, err)
			nfail++
			return
		}
		out := string(bout)

		// result must convert back to the same float
		g, err := strconv.ParseFloat(out, bits)
		if err != nil {
			t.Errorf("Marshal(%T(%g)) = %q, cannot parse back: %v", vf, vf, out, err)
			nfail++
			return
		}
		if f != g || fmt.Sprint(f) != fmt.Sprint(g) { // fmt.Sprint handles ±0
			t.Errorf("Marshal(%T(%g)) = %q (is %g, not %g)", vf, vf, out, float32(g), vf)
			nfail++
			return
		}

		bad := badFloatREs
		if bits == 64 {
			bad = bad[:len(bad)-2]
		}
		for _, re := range bad {
			if re.MatchString(out) {
				t.Errorf("Marshal(%T(%g)) = %q, must not match /%s/", vf, vf, out, re)
				nfail++
				return
			}
		}
	}

	var (
		bigger  = math.Inf(+1)
		smaller = math.Inf(-1)
	)

	var digits = "1.2345678901234567890123"
	for i := len(digits); i >= 2; i-- {
		if testing.Short() && i < len(digits)-4 {
			break
		}
		for exp := -30; exp <= 30; exp++ {
			for _, sign := range "+-" {
				for bits := 32; bits <= 64; bits += 32 {
					s := fmt.Sprintf("%c%se%d", sign, digits[:i], exp)
					f, err := strconv.ParseFloat(s, bits)
					if err != nil {
						log.Fatal(err)
					}
					next := math.Nextafter
					if bits == 32 {
						next = func(g, h float64) float64 {
							return float64(math.Nextafter32(float32(g), float32(h)))
						}
					}
					test(f, bits)
					test(next(f, bigger), bits)
					test(next(f, smaller), bits)
					if nfail > 50 {
						t.Fatalf("stopping test early")
					}
				}
			}
		}
	}
	test(0, 64)
	test(math.Copysign(0, -1), 64)
	test(0, 32)
	test(math.Copysign(0, -1), 32)
}

var encodeStringTests = []struct {
	in  string
	out string
}{
	{"\x00", `"\u0000"`},
	{"\x01", `"\u0001"`},
	{"\x02", `"\u0002"`},
	{"\x03", `"\u0003"`},
	{"\x04", `"\u0004"`},
	{"\x05", `"\u0005"`},
	{"\x06", `"\u0006"`},
	{"\x07", `"\u0007"`},
	{"\x08", `"\u0008"`},
	{"\x09", `"\t"`},
	{"\x0a", `"\n"`},
	{"\x0b", `"\u000b"`},
	{"\x0c", `"\u000c"`},
	{"\x0d", `"\r"`},
	{"\x0e", `"\u000e"`},
	{"\x0f", `"\u000f"`},
	{"\x10", `"\u0010"`},
	{"\x11", `"\u0011"`},
	{"\x12", `"\u0012"`},
	{"\x13", `"\u0013"`},
	{"\x14", `"\u0014"`},
	{"\x15", `"\u0015"`},
	{"\x16", `"\u0016"`},
	{"\x17", `"\u0017"`},
	{"\x18", `"\u0018"`},
	{"\x19", `"\u0019"`},
	{"\x1a", `"\u001a"`},
	{"\x1b", `"\u001b"`},
	{"\x1c", `"\u001c"`},
	{"\x1d", `"\u001d"`},
	{"\x1e", `"\u001e"`},
	{"\x1f", `"\u001f"`},
}

func TestEncodeString(t *testing.T) {
	for _, tt := range encodeStringTests {
		b, err := json.Marshal(tt.in)
		if err != nil {
			t.Errorf("Marshal(%q): %v", tt.in, err)
			continue
		}
		out := string(b)
		if out != tt.out {
			t.Errorf("Marshal(%q) = %#q, want %#q", tt.in, out, tt.out)
		}
	}
}

type jsonbyte byte

func (b jsonbyte) MarshalJSON() ([]byte, error) { return tenc(`{"JB":%d}`, b) }

type textbyte byte

func (b textbyte) MarshalText() ([]byte, error) { return tenc(`TB:%d`, b) }

type jsonint int

func (i jsonint) MarshalJSON() ([]byte, error) { return tenc(`{"JI":%d}`, i) }

type textint int

func (i textint) MarshalText() ([]byte, error) { return tenc(`TI:%d`, i) }

func tenc(format string, a ...interface{}) ([]byte, error) {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, format, a...)
	return buf.Bytes(), nil
}

// Issue 13783
func TestEncodeBytekind(t *testing.T) {
	testdata := []struct {
		data interface{}
		want string
	}{
		{byte(7), "7"},
		{jsonbyte(7), `{"JB":7}`},
		{textbyte(4), `"TB:4"`},
		{jsonint(5), `{"JI":5}`},
		{textint(1), `"TI:1"`},
		{[]byte{0, 1}, `"AAE="`},

		{[]jsonbyte{0, 1}, `[{"JB":0},{"JB":1}]`},
		{[][]jsonbyte{{0, 1}, {3}}, `[[{"JB":0},{"JB":1}],[{"JB":3}]]`},
		{[]textbyte{2, 3}, `["TB:2","TB:3"]`},

		{[]jsonint{5, 4}, `[{"JI":5},{"JI":4}]`},
		{[]textint{9, 3}, `["TI:9","TI:3"]`},
		{[]int{9, 3}, `[9,3]`},
	}
	for i, d := range testdata {
		js, err := json.Marshal(d.data)
		if err != nil {
			t.Error(err)
			continue
		}
		got, want := string(js), d.want
		if got != want {
			t.Errorf("%d: got %s, want %s", i, got, want)
		}
	}
}

// golang.org/issue/8582
func TestEncodePointerString(t *testing.T) {
	type stringPointer struct {
		N *int64 `json:"n,string"`
	}
	var n int64 = 42
	b, err := json.Marshal(stringPointer{N: &n})
	if err != nil {
		t.Fatalf("Marshal: %v", err)
	}
	if got, want := string(b), `{"n":"42"}`; got != want {
		t.Errorf("Marshal = %s, want %s", got, want)
	}
	var back stringPointer
	err = json.Unmarshal(b, &back)
	if err != nil {
		t.Fatalf("Unmarshal: %v", err)
	}
	if back.N == nil {
		t.Fatalf("Unmarshaled nil N field")
	}
	if *back.N != 42 {
		t.Fatalf("*N = %d; want 42", *back.N)
	}
}

type SamePointerNoCycle struct {
	Ptr1, Ptr2 *SamePointerNoCycle
}

var samePointerNoCycle = &SamePointerNoCycle{}

type PointerCycle struct {
	Ptr *PointerCycle
}

var pointerCycle = &PointerCycle{}

type PointerCycleIndirect struct {
	Ptrs []interface{}
}

var pointerCycleIndirect = &PointerCycleIndirect{}

func init() {
	ptr := &SamePointerNoCycle{}
	samePointerNoCycle.Ptr1 = ptr
	samePointerNoCycle.Ptr2 = ptr

	pointerCycle.Ptr = pointerCycle
	pointerCycleIndirect.Ptrs = []interface{}{pointerCycleIndirect}
}

func TestSamePointerNoCycle(t *testing.T) {
	if _, err := json.Marshal(samePointerNoCycle); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

var unsupportedValues = []interface{}{
	math.NaN(),
	math.Inf(-1),
	math.Inf(1),
	pointerCycle,
	pointerCycleIndirect,
}

func TestUnsupportedValues(t *testing.T) {
	for _, v := range unsupportedValues {
		if _, err := json.Marshal(v); err != nil {
			if _, ok := err.(*json.UnsupportedValueError); !ok {
				t.Errorf("for %v, got %T want UnsupportedValueError", v, err)
			}
		} else {
			t.Errorf("for %v, expected error", v)
		}
	}
}
