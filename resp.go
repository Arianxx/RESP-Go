package resp

import (
	"fmt"
	"github.com/arianxx/camellia-io"
	"log"
)

// Command represents a command parsed from the input bytes.
type Command struct {
	// Raw is the original data received from the client.
	Raw []byte
	// Args is the a series of arguments that make up the command.
	Args [][]byte
	Type string
}

func NewCommand(t string) *Command {
	return &Command{[]byte{}, [][]byte{}, t}
}

// Parser tries to parse the input bytes to a series of Command.
type Parser struct {
	raw []byte
	// nowConverter represents the converter that is now used to process the input stream
	nowConverter *Converter
	// end indicates whether it should continue parsing.
	end bool
}

func NewParser() *Parser {
	return &Parser{raw: []byte{}}
}

// AppendRawData appends a raw data to the parser.
func (p *Parser) AppendRawData(in []byte) {
	p.raw = append(p.raw, in...)
}

/* Parse parses the received data.
 * Return a Command pointers if there are some completed parsings, otherwise it remains nil.
 * Return a errProtocol if the error data has been detected.
 * If a err was be returned, then the parser will be unusable and should be replaced by a new Parser.
 */
func (p *Parser) Parse() (cmd *Command, err error) {
	if p.end {
		return nil, fmt.Errorf("unusable ended parser")
	}

	if p.nowConverter != nil {
		cmd, err, p.raw = (*p.nowConverter).Parse(p.raw)
		if err != nil {
			p.end = true
			cmd = nil
			return
		}
		if cmd != nil {
			p.nowConverter = nil
			return
		}
	} else if len(p.raw) != 0 {
		for s, f := range converters {
			if s == p.raw[0] {
				c := f()
				p.nowConverter = &c
				break
			}
		}

		if p.nowConverter == nil {
			p.end = true
			return nil, fmt.Errorf("unknown type symbol")
		}

		return p.Parse()
	}

	return
}

type Server struct {
	*camellia.Server
	proc CommandProc

	Event *camellia.Event
}

func NewServer(net, addr string, f CommandProc) (*Server, error) {
	s := &Server{camellia.NewServer(), f, nil}
	lis, err := camellia.NewListener(net, addr, s.El)
	if err != nil {
		return nil, err
	}
	s.AddListener(lis)
	s.Event = &camellia.Event{Data: s.loopProc}
	s.AddEvent(s.Event)
	return s, nil
}

func (s *Server) StartServe() error {
	if s.proc == nil {
		log.Fatal("empty proc")
	}

	return s.Server.StartServe()
}

func (s *Server) loopProc(el *camellia.EventLoop, connPtr *interface{}) {
	conn := (*connPtr).(*camellia.Conn)
	parserInterface := conn.GetContext()
	if parserInterface == nil {
		parserInterface = NewParser()
		conn.SetContext(parserInterface)
	}
	parser := parserInterface.(*Parser)

	parser.AppendRawData(conn.Read())
	cmd, err := parser.Parse()
	if err != nil {
		conn.SetContext(NewParser())
		conn.Write(AppendError([]byte{}, []byte(err.Error())))
		return
	}

	if cmd != nil {
		res := &[]byte{}
		s.proc(cmd, conn, res)
		conn.Write(*res)
		return
	}
}

type CommandProc func(cmd *Command, conn *camellia.Conn, res *[]byte)
