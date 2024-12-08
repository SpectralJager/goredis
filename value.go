package goredis

import (
	"fmt"
	"io"
	"strconv"
)

type Value struct {
	typ   byte
	str   string
	num   int
	array []Value
}

func StringValue(str string) Value {
	return Value{
		typ: STRING,
		str: str,
	}
}

func IntegerValue(num int) Value {
	return Value{
		typ: INTEGER,
		num: num,
	}
}

func ErrorValue(err error) Value {
	return Value{
		typ: ERROR,
		str: err.Error(),
	}
}

func (v Value) Integer() int {
	if v.typ == INTEGER {
		return v.num
	}
	return 0
}

func (v Value) Bulk() string {
	if v.typ == BULK {
		return v.str
	}
	return ""
}

func (v Value) Element(index int) Value {
	if v.typ != ARRAY || len(v.array) <= index || index < 0 {
		return Value{}
	}
	return v.array[index]
}

func (v Value) Marshall() []byte {
	switch v.typ {
	case ARRAY:
		return v.marshallArray()
	case BULK:
		return v.marshallBulk()
	case STRING:
		return v.marshallString()
	case INTEGER:
		return v.marshallInteger()
	case ERROR:
		return v.marshallError()
	default:
		return []byte{}
	}
}

func (v Value) marshallInteger() []byte {
	bytes := []byte{}
	bytes = append(bytes, INTEGER)
	bytes = append(bytes, strconv.Itoa(v.num)...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallString() []byte {
	bytes := []byte{}
	bytes = append(bytes, STRING)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallBulk() []byte {
	bytes := []byte{}
	bytes = append(bytes, BULK)
	if len(v.str) != 0 {
		bytes = append(bytes, strconv.Itoa(len(v.str))...)
		bytes = append(bytes, '\r', '\n')
		bytes = append(bytes, v.str...)
	} else {
		bytes = append(bytes, []byte("-1")...)
	}
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func (v Value) marshallArray() []byte {
	bytes := []byte{}
	bytes = append(bytes, ARRAY)
	bytes = append(bytes, strconv.Itoa(len(v.array))...)
	bytes = append(bytes, '\r', '\n')
	for _, val := range v.array {
		bytes = append(bytes, val.Marshall()...)
	}
	return bytes
}

func (v Value) marshallError() []byte {
	bytes := []byte{}
	bytes = append(bytes, ERROR)
	bytes = append(bytes, v.str...)
	bytes = append(bytes, '\r', '\n')
	return bytes
}

func WriteValue(w io.Writer, val Value) error {
	bytes := val.Marshall()
	_, err := w.Write(bytes)
	if err != nil {
		return fmt.Errorf("can't write value: %w", err)
	}
	return nil
}
