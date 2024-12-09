package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
)

type Resp struct {
	reader *bufio.Reader
}

func NewResp(r io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(r),
	}
}

func (r *Resp) readLine() ([]byte, int, error) {
	n := 0
	line := make([]byte, 0, r.reader.Buffered())
	// line := []byte{}
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, fmt.Errorf("can't read line from buffer: %w", err)
		}
		n += 1
		line = append(line, b)
		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readLength() (int, int, error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, fmt.Errorf("can't read line for parsing integer: %w", err)
	}
	integer, err := strconv.Atoi(string(line))
	if err != nil {
		return 0, 0, fmt.Errorf("can't parse integer: %w", err)
	}
	return integer, n, nil
}

func (r *Resp) readArray() (Value, error) {
	v := Value{}
	v.typ = ARRAY
	r.reader.ReadByte()
	length, _, err := r.readLength()
	if err != nil {
		return Value{}, fmt.Errorf("can't read length of array: %w", err)
	}
	v.array = make([]Value, 0, length)
	for i := 0; i < length; i++ {
		val, err := r.Read()
		if err != nil {
			return Value{}, fmt.Errorf("can't read value: %w", err)
		}
		v.array = append(v.array, val)
	}
	return v, nil
}

func (r *Resp) readBulk() (Value, error) {
	v := Value{}
	v.typ = BULK
	r.reader.ReadByte()
	length, _, err := r.readLength()
	if err != nil {
		return Value{}, fmt.Errorf("can't read length of bulk string: %w", err)
	}
	if length == -1 {
		v.str = ""
		return v, nil
	}
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, fmt.Errorf("can't read bulk value: %w", err)
	}
	v.str = string(line)
	return v, nil
}

func (r *Resp) readString() (Value, error) {
	v := Value{}
	v.typ = STRING
	r.reader.ReadByte()
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, fmt.Errorf("can't read line of string: %w", err)
	}
	v.str = string(line)
	return v, nil
}

func (r *Resp) readError() (Value, error) {
	v := Value{}
	v.typ = ERROR
	r.reader.ReadByte()
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, fmt.Errorf("can't read line of string: %w", err)
	}
	v.str = string(line)
	return v, nil
}

func (r *Resp) readInteger() (Value, error) {
	v := Value{}
	v.typ = INTEGER
	r.reader.ReadByte()
	line, _, err := r.readLine()
	if err != nil {
		return Value{}, fmt.Errorf("can't read line of integer: %w", err)
	}
	v.num, _ = strconv.Atoi(string(line))
	return v, nil
}

func (r *Resp) Read() (Value, error) {
	typ, err := r.reader.Peek(1)
	if err != nil {
		return Value{}, fmt.Errorf("can't read: %w", err)
	}
	switch typ[0] {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	case STRING:
		return r.readString()
	case INTEGER:
		return r.readInteger()
	case ERROR:
		return r.readError()
	default:
		return Value{}, fmt.Errorf("can't read: unexpected type: %v", string(typ))
	}
}
