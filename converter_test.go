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
