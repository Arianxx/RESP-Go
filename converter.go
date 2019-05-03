package resp

import "fmt"

var (
	errInteger      = &errProtocol{"failed to parse integer"}
	errStringLength = &errProtocol{"failed to parse bulk string length"}
	errStringEnding = &errProtocol{"failed to parse the string ending symbol"}
	errArrayLength  = &errProtocol{"failed to parse array length"}
	errArrayEnding  = &errProtocol{"failed to parse the array ending symbol"}
)

// errProtocol used to represents the error about incorrect data that does not conform to the protocol
type errProtocol struct {
	msg string
}

func (p *errProtocol) Error() string {
	return p.msg
}

type Converter interface {
	Parse(raw []byte) (*Command, error, []byte)
}

type SimpleString struct {
	c *Command
}

var simpleStringSign = byte('+')

func NewSimpleString() *SimpleString {
	return &SimpleString{NewCommand()}
}

func NewSimpleStringConverter() Converter {
	return NewSimpleString()
}

func (s *SimpleString) Parse(raw []byte) (cmd *Command, err error, surplus []byte) {
	endIndex := -1
	for i, c := range raw {
		if c == '\r' {
			if (i == 0 && len(s.c.Raw) != 0 && s.c.Raw[len(s.c.Raw)-1] == '\n') ||
				raw[i-1] == '\n' {
				endIndex = i
			}
			break
		}
	}
	if endIndex == -1 {
		s.c.Raw = append(s.c.Raw, raw...)
		return nil, nil, []byte{}
	}

	s.c.Raw = append(s.c.Raw, raw[:endIndex+1]...)
	s.c.Args = append(s.c.Args, s.c.Raw[1:len(s.c.Raw)-2])
	return s.c, nil, raw[endIndex+1:]
}

type RespError struct {
	SimpleString
}

var respErrorSign = byte('-')

func NewRespError() *RespError {
	return &RespError{*NewSimpleString()}
}

func NewRespErrorConverter() Converter {
	return NewRespError()
}

type Integer struct {
	c *Command
}

var integerSign = byte(':')

func NewInteger() *Integer {
	return &Integer{NewCommand()}
}

func NewIntegerConverter() Converter {
	return NewInteger()
}

func (t *Integer) Parse(raw []byte) (*Command, error, []byte) {
	endIndex := -1
	for i, c := range raw {
		if len(t.c.Raw) == 0 && i == 0 && c == integerSign {
			continue
		}

		if (len(t.c.Raw) == 0 && i == 1 && c == '-') ||
			(len(t.c.Raw) == 1 && i == 0 && c == '-') {
			continue
		}

		if c < '0' || c > '9' {
			return nil, errInteger, raw
		}

		if c == '\r' {
			if (i == 0 && len(t.c.Raw) != 0 && t.c.Raw[len(t.c.Raw)-1] == '\n') ||
				raw[i-1] == '\n' {
				endIndex = i
			}
			break
		}
	}
	if endIndex == -1 {
		t.c.Raw = append(t.c.Raw, raw...)
		return nil, nil, []byte{}
	}

	t.c.Raw = append(t.c.Raw, raw[:endIndex+1]...)
	t.c.Args = append(t.c.Args, t.c.Raw[1:len(t.c.Raw)-2])
	return t.c, nil, raw[endIndex+1:]
}

type BulkString struct {
	length    int
	gotLength bool
	c         *Command
}

var bulkStringSign = byte('$')

func NewBulkString() *BulkString {
	b := BulkString{0, false, NewCommand()}
	b.c.Args = append(b.c.Args, []byte{})
	return &b
}

func NewBulkStringConverter() Converter {
	return NewBulkString()
}

func (b *BulkString) Parse(raw []byte) (*Command, error, []byte) {
	if !b.gotLength {
		return b.parseToGetLength(raw)
	} else {
		return b.parseToGetString(raw)
	}
}

func (b *BulkString) parseToGetLength(raw []byte) (*Command, error, []byte) {
	for i, c := range raw {
		if len(b.c.Raw) == 0 && i == 0 && c == bulkStringSign {
			continue
		}

		if c == '-' {
			if (i == 0 && len(b.c.Raw) == 1) || (i == 1 && len(b.c.Raw) == 0) {
				b.length = -1
				continue
			}
		}

		if c == '\r' {
			if i+1 >= len(raw) || raw[i+1] != '\n' {
				return nil, errStringLength, raw
			}
			b.gotLength = true
			b.c.Raw = append(b.c.Raw, raw[:i+2]...)
			return b.Parse(raw[i+2:])
		}

		if c < '0' || c > '9' {
			return nil, errStringLength, raw
		}

		if b.length == -1 {
			continue
		}

		b.length = 10*b.length + int(c) - '0'
	}
	b.c.Raw = append(b.c.Raw, raw...)
	return nil, nil, []byte{}
}

func (b *BulkString) parseToGetString(raw []byte) (*Command, error, []byte) {
	if len(raw) == 0 {
		return nil, nil, raw
	}

	if b.length == -1 {
		if len(raw) >= 2 && string(raw[:2]) == "\n\r" {
			// &Command{nil, nil} means resp null
			return &Command{nil, nil}, nil, raw[2:]
		} else {
			return nil, errStringEnding, raw
		}
	} else if b.length > 0 {
		for i, c := range raw {
			b.c.Raw = append(b.c.Raw, c)
			b.c.Args[0] = append(b.c.Args[0], c)

			b.length--
			if b.length == 0 {
				if i+3 >= len(raw) || string(raw[i+1:i+3]) != "\n\r" {
					return nil, errStringEnding, raw
				}
				return b.c, nil, raw[i+3:]
			}
		}
		return nil, nil, []byte{}
	}

	return nil, errStringEnding, raw
}

type Array struct {
	length    int
	gotLength bool
	inner     Converter
	c         *Command
}

var arraySign = byte('*')

func NewArray() *Array {
	return &Array{0, false, nil, NewCommand()}
}

func NewArrayConverter() Converter {
	return NewArray()
}

func (a *Array) Parse(raw []byte) (*Command, error, []byte) {
	return nil, nil, nil
}

func (b *Array) parseToGetLength(raw []byte) (*Command, error, []byte) {
	for i, c := range raw {
		if len(b.c.Raw) == 0 && i == 0 && c == arraySign {
			continue
		}

		if c == '-' {
			if (i == 0 && len(b.c.Raw) == 1) || (i == 1 && len(b.c.Raw) == 0) {
				b.length = -1
				continue
			}
		}

		if c == '\r' {
			if i+1 >= len(raw) || raw[i+1] != '\n' {
				return nil, errArrayLength, raw
			}
			b.gotLength = true
			b.c.Raw = append(b.c.Raw, raw[:i+2]...)
			return b.Parse(raw[i+2:])
		}

		if c < '0' || c > '9' {
			return nil, errArrayLength, raw
		}

		if b.length == -1 {
			continue
		}

		b.length = 10*b.length + int(c) - '0'
	}
	b.c.Raw = append(b.c.Raw, raw...)
	return nil, nil, []byte{}
}

func (b *Array) parseToGetString(raw []byte) (*Command, error, []byte) {
	if len(raw) == 0 {
		return nil, nil, raw
	}

	if b.length == -1 {
		if len(raw) >= 2 && string(raw[:2]) == "\n\r" {
			// &Command{nil, nil} means resp null
			return &Command{nil, nil}, nil, raw[2:]
		} else {
			return nil, errArrayEnding, raw
		}
	} else if b.length > 0 {
		if b.inner == nil {
			if innerConverter, ok := converters[raw[0]]; !ok {
				return nil, fmt.Errorf("unknown type symbol"), raw
			} else {
				b.inner = innerConverter()
				return b.Parse(raw)
			}
		} else {
			cmd, err, surplus := b.inner.Parse(raw)
			if err != nil {
				return nil, err, raw
			}
			if cmd != nil {
				b.length--
				b.c.Raw = append(b.c.Raw, cmd.Raw...)
				b.c.Args = append(b.c.Args, b.c.Args[0])
				b.inner = nil
			}
			return nil, nil, surplus
		}
	} else if b.length == 0 {
		if b.c.Raw[len(b.c.Raw)-1] == '\n' {
			if raw[0] == '\r' {
				b.c.Raw = append(b.c.Raw, raw[0])
				return b.c, nil, raw[1:]
			} else {
				return nil, errArrayEnding, raw
			}
		} else {
			if len(raw) > 2 {
				if string(raw[:2]) == "\n\r" {
					b.c.Raw = append(b.c.Raw, raw[:2]...)
					return b.c, nil, raw[2:]
				} else {
					return nil, errArrayEnding, raw
				}
			} else {
				b.c.Raw = append(b.c.Raw, raw...)
				return nil, nil, []byte{}
			}
		}
	}

	return nil, errStringEnding, raw
}

type ConverterConstructor func() Converter


var converters = map[byte]ConverterConstructor{
	simpleStringSign: NewSimpleStringConverter,
	respErrorSign:    NewRespErrorConverter,
	integerSign:      NewIntegerConverter,
	bulkStringSign:   NewBulkStringConverter,
	arraySign:        NewArrayConverter,
}
