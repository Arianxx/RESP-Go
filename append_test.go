package resp

import (
	"strconv"
	"testing"
)

var arrayTable = [][][]byte{
	{

		[]byte("+test\r\n"),
		[]byte(":100\r\n"),
	},
	{
		[]byte("-err\r\n"),
		[]byte("$-1\r\n"),
	},
}

func TestAppendArray(t *testing.T) {
	for _, a := range arrayTable {
		test := []byte{}
		expected := "*" + strconv.Itoa(len(a)) + "\r\n"
		for _, perMsg := range a {
			expected += string(perMsg)
		}
		expected += "\r\n"
		test = AppendArray(test, a...)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}

var bulkTable = [][]byte{
	[]byte("test"),
	[]byte("tets\n\r"),
}

func TestAppendBulkString(t *testing.T) {
	for _, bulk := range bulkTable {
		test := []byte{}
		expected := "$" + strconv.Itoa(len(bulk)) + "\r\n" + string(bulk) + "\r\n"
		test = AppendBulkString(test, bulk)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}

var errorTable = [][]byte{
	[]byte("-err test"),
	[]byte("-Error test"),
}

func TestAppendError(t *testing.T) {
	for _, err := range errorTable {
		expected := "-" + string(err) + "\r\n"
		test := AppendError([]byte{}, err)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}

var integerTable = []int{1000, -1, 0}

func TestAppendInteger(t *testing.T) {
	for _, i := range integerTable {
		expected := ":" + strconv.Itoa(i) + "\r\n"
		test := AppendInteger([]byte{}, i)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}

var simpleStringTable = [][]byte{
	[]byte("test"),
}

func TestAppendSimpleString(t *testing.T) {
	for _, s := range simpleStringTable {
		expected := "+" + string(s) + "\r\n"
		test := AppendSimpleString([]byte{}, s)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}
