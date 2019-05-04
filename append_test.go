package resp

import (
	"strconv"
	"testing"
)

var arrayTable = [][][]byte{
	{

		[]byte("+test\n"),
		[]byte(":100\n"),
	},
	{
		[]byte("-err\n"),
		[]byte("$-1\n"),
	},
}

func TestAppendArray(t *testing.T) {
	for _, a := range arrayTable {
		test := []byte{}
		expected := "*" + strconv.Itoa(len(a)) + "\n"
		for _, perMsg := range a {
			expected += string(perMsg)
		}
		expected += "\n"
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
		expected := "$" + strconv.Itoa(len(bulk)) + "\n" + string(bulk) + "\n"
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
		expected := "-" + string(err) + "\n"
		test := AppendError([]byte{}, err)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}

var integerTable = []int{1000, -1, 0}

func TestAppendInteger(t *testing.T) {
	for _, i := range integerTable {
		expected := ":" + strconv.Itoa(i) + "\n"
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
		expected := "+" + string(s) + "\n"
		test := AppendSimpleString([]byte{}, s)
		if string(test) != expected {
			t.Fatal("expected", expected, "got", string(test))
		}
	}
}
