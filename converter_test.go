package resp

import "testing"

func TestSimpleString_Parse_uncompleted(t *testing.T) {
	raw := []byte("+test")
	s := NewSimpleString()
	cmd, err, surplus := s.Parse(raw)
	if err != nil {
		t.Fatal("expected nil error, got ", err)
	}
	if cmd != nil {
		t.Fatal("expected nil, got ", cmd)
	}
	if len(surplus) != 0 {
		t.Fatal("expected zero surplus length, got", surplus)
	}
}

func TestSimpleString_Parse_successful(t *testing.T) {
	raw := []byte("+test\n")
	s := NewSimpleString()
	expectedCmd := &Command{
		Raw:  raw,
		Args: [][]byte{[]byte("test")},
	}

	cmd, err, surplus := s.Parse(raw)
	if err != nil {
		t.Fatal("expected nil, got ", err)
	}
	if len(surplus) != 0 {
		t.Fatal("expected zero surplus length, got", surplus)
	}
	if string(cmd.Raw) != string(expectedCmd.Raw) || string(cmd.Args[0]) != string(expectedCmd.Args[0]) {
		t.Fatal("expected", expectedCmd, "got", cmd)
	}
}

func TestInteger_Parse(t *testing.T) {
	raw := []byte(":1000\n")
	s := NewInteger()
	cmd, err, surplus := s.Parse(raw)
	if len(surplus) != 0 {
		t.Fatal("expected 0, got ", string(surplus))
	}
	if err != nil {
		t.Fatal("expected error, got ", err)
	}
	if string(cmd.Args[0]) != "1000" {
		t.Fatal("expected", "1000", "got", string(cmd.Args[0]))
	}
}

func TestInteger_Parse2(t *testing.T) {
	raw := []byte(":100")
	s := NewInteger()
	_, _, _ = s.Parse(raw)
	raw = []byte("0\n")
	cmd, err, surplus := s.Parse(raw)
	if len(surplus) != 0 {
		t.Fatal("expected 0, got ", string(surplus))
	}
	if err != nil {
		t.Fatal("expected error, got ", err)
	}
	if string(cmd.Args[0]) != "1000" {
		t.Fatal("expected", "1000", "got", string(cmd.Args[0]))
	}
}

func TestBulkString_Parse(t *testing.T) {
	raw := []byte("$6\ntesthh\n")
	s := NewBulkString()
	cmd, err, surplus := s.Parse(raw)
	if len(surplus) != 0 {
		t.Fatal("expected 0, got ", string(surplus))
	}
	if err != nil {
		t.Fatal("expected error, got ", err)
	}
	if string(cmd.Args[0]) != "testhh" {
		t.Fatal("expected", "testhh", "got", string(cmd.Args[0]))
	}
}

func TestBulkString_Parse2(t *testing.T) {
	raw := []byte("$6\ntest")
	s := NewBulkString()
	_, _, _ = s.Parse(raw)
	raw = []byte("hh\n")
	cmd, err, surplus := s.Parse(raw)
	if len(surplus) != 0 {
		t.Fatal("expected 0, got ", string(surplus))
	}
	if err != nil {
		t.Fatal("expected error, got ", err)
	}
	if string(cmd.Args[0]) != "testhh" {
		t.Fatal("expected", "testhh", "got", string(cmd.Args[0]))
	}
}

func TestArray_Parse(t *testing.T) {
	raw := []byte("*-1\n")
	s := NewArray()
	cmd, err, _ := s.Parse(raw)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if cmd != Nil {
		t.Fatal("expected", Nil, "got", cmd)
	}
}

func TestArray_Parse2(t *testing.T) {
	raw := []byte("*2\n+test\n:1000\n\n")
	s := NewArray()
	cmd, err, surplus := s.Parse(raw)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if len(surplus) != 0 {
		t.Fatal("expected", 0, "got", string(surplus))
	}
	if string(cmd.Args[0]) != "test" {
		t.Fatal("expected", "test", "got", string(cmd.Args[0]))
	}
	if string(cmd.Args[1]) != "1000" {
		t.Fatal("expected", "test", "got", string(cmd.Args[0]))
	}
}
