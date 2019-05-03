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
	raw := []byte("+test\n\r")
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
	raw := []byte(":1000\n\r")
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
	raw = []byte("0\n\r")
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
	raw := []byte("$6\n\rtesthh\n\r")
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
	raw := []byte("$6\n\rtest")
	s := NewBulkString()
	_, _, _ = s.Parse(raw)
	raw = []byte("hh\n\r")
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
	raw := []byte("*-1\n\r")
	s := NewArray()
	cmd, err, _ := s.Parse(raw)
	if err != nil {
		t.Fatal("expected", nil, "got", err)
	}
	if cmd != RespNil {
		t.Fatal("expected", RespNil, "got", cmd)
	}
}

func TestArray_Parse2(t *testing.T) {
	raw := []byte("*2\n\r+test\n\r:1000\n\r\n\r")
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
