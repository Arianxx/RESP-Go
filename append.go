package resp

import "strconv"

func appendPredix(b []byte, c byte, n int) []byte {
	b = append(b, c)
	b = strconv.AppendInt(b, int64(n), 10)
	return append(b, '\r', '\n')
}

func appendBulkLength(b []byte, n int) []byte {
	return appendPredix(b, '$', n)
}

func appendArrayLength(b []byte, n int) []byte {
	return appendPredix(b, '*', n)
}

func appendLineEnding(b []byte) []byte {
	return append(b, '\r', '\n')
}

func AppendSimpleString(b, c []byte) []byte {
	return appendLineEnding(append(append(b, '+'), c...))
}

func AppendError(b, c []byte) []byte {
	return appendLineEnding(append(append(b, '-'), c...))
}

func AppendInteger(b []byte, n int) []byte {
	return appendLineEnding(append(append(b, ':'), []byte(strconv.Itoa(n))...))
}

func AppendBulkString(b, c []byte) []byte {
	return appendLineEnding(append(appendBulkLength(b, len(c)), c...))
}

func AppendArray(b []byte, c ...[]byte) []byte {
	b = appendArrayLength(b, len(c))
	for _, d := range c {
		b = append(b, d...)
	}
	return appendLineEnding(b)
}
