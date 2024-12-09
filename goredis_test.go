package resp

import (
	"bytes"
	"reflect"
	"strings"
	"testing"
)

func TestRespRead(t *testing.T) {
	testCases := []struct {
		input  string
		expect Value
	}{
		{
			input:  "+OK\r\n",
			expect: Value{typ: STRING, str: "OK"},
		},
		{
			input:  "-Error message\r\n",
			expect: Value{typ: ERROR, str: "Error message"},
		},
		{
			input:  "$5\r\nhello\r\n",
			expect: Value{typ: BULK, str: "hello"},
		},
		{
			input:  "*0\r\n",
			expect: Value{typ: ARRAY, array: []Value{}},
		},
		{
			input: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
			expect: Value{typ: ARRAY, array: []Value{
				{typ: BULK, str: "hello"},
				{typ: BULK, str: "world"},
			}},
		},
		{
			input: "*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n",
			expect: Value{typ: ARRAY, array: []Value{
				{typ: BULK, str: "hello"},
				{typ: BULK, str: ""},
				{typ: BULK, str: "world"},
			}},
		},
		{
			input: "*3\r\n:1\r\n:2\r\n:3\r\n",
			expect: Value{typ: ARRAY, array: []Value{
				{typ: INTEGER, num: 1},
				{typ: INTEGER, num: 2},
				{typ: INTEGER, num: 3},
			}},
		},
	}
	for _, tC := range testCases {
		t.Run(tC.input, func(t *testing.T) {
			resp := NewResp(strings.NewReader(tC.input))
			val, err := resp.Read()
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(val, tC.expect) {
				t.Fail()
			}
		})
	}
}

func TestValueWrite(t *testing.T) {
	testCases := []struct {
		input  Value
		expect string
	}{
		{
			input:  Value{typ: STRING, str: "OK"},
			expect: "+OK\r\n",
		},
		{
			input:  Value{typ: ERROR, str: "Error message"},
			expect: "-Error message\r\n",
		},
		{
			input:  Value{typ: BULK, str: "hello"},
			expect: "$5\r\nhello\r\n",
		},
		{
			input:  Value{typ: ARRAY, array: []Value{}},
			expect: "*0\r\n",
		},
		{
			input: Value{typ: ARRAY, array: []Value{
				{typ: BULK, str: "hello"},
				{typ: BULK, str: "world"},
			}},
			expect: "*2\r\n$5\r\nhello\r\n$5\r\nworld\r\n",
		},
		{
			input: Value{typ: ARRAY, array: []Value{
				{typ: BULK, str: "hello"},
				{typ: BULK, str: ""},
				{typ: BULK, str: "world"},
			}},
			expect: "*3\r\n$5\r\nhello\r\n$-1\r\n$5\r\nworld\r\n",
		},
		{
			input: Value{typ: ARRAY, array: []Value{
				{typ: INTEGER, num: 1},
				{typ: INTEGER, num: 2},
				{typ: INTEGER, num: 3},
			}},
			expect: "*3\r\n:1\r\n:2\r\n:3\r\n",
		},
	}
	for _, tC := range testCases {
		t.Run(tC.expect, func(t *testing.T) {
			var buff bytes.Buffer
			err := WriteValue(&buff, tC.input)
			if err != nil {
				t.Fatal(err)
			}
			if tC.expect != buff.String() {
				t.Fail()
			}
		})
	}
}
