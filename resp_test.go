package resp

import (
	"math/rand"
	"testing"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

var letterRunes = []byte("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

func randStringRunes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return b
}

func TestParser_AppendRawData(t *testing.T) {
	for i := 1; i <= 5; i++ {
		randByte := randStringRunes(i * i)
		p := NewParser()
		p.AppendRawData(randByte)
		if string(randByte) != string(p.raw) {
			t.Fatal(
				"expected", randByte,
				"got", p.raw,
			)
		}
	}
}

func TestParser_Parse_unknownSymbol(t *testing.T) {
	// Spaces never becomes identifiers.
	e := []byte(" ")
	p := NewParser()
	p.AppendRawData(e)
	cmd, err := p.Parse()
	if cmd != nil {
		t.Fatal("expected nil but got: ", cmd)
	}
	if err == nil {
		t.Fatalf("expected a error but got nil")
	}
	if err.Error() != "unknown type symbol" {
		t.Fatal("expected", "unknown type symbol", "got", err.Error())
	}
}

func TestParser_Parse_uncompletedInput(t *testing.T) {
	e := []byte("+test")
	p := NewParser()
	p.AppendRawData(e)
	cmd, err := p.Parse()
	if cmd != nil {
		t.Fatal("expected nil cmd but got: ", cmd)
	}
	if err != nil {
		t.Fatal("expected nil error but got: ", err)
	}
	if len(p.raw) != 0 {
		t.Fatal("expected zero length p.raw but got: ", p.raw)
	}
}

func TestParser_Parse_successful(t *testing.T) {
	randByte := randStringRunes(5)
	e := []byte("+" + string(randByte))
	p := NewParser()
	p.AppendRawData(e)
	cmd, err := p.Parse()
	if cmd != nil || err != nil {
		t.Fatal("expected nil, got ", cmd, err)
	}
	p.AppendRawData([]byte("\r\n"))
	cmd, err = p.Parse()
	if err != nil {
		t.Fatal("expected nil error but got: ", err)
	}
	if cmd == nil {
		t.Fatalf("expected a commandd but got nil")
	}
	if string(cmd.Args[0]) != string(randByte) {
		t.Fatal("expceted", randByte, "got", cmd.Args[0])
	}
}
