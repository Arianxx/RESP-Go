package resp

import (
	"fmt"
)

// Command represents a command parsed from the input bytes.
type Command struct {
	// Raw is the original data received from the client.
	Raw []byte
	// Args is the a series of arguments that make up the command.
	Args [][]byte
}

func NewCommand() *Command {
	return &Command{[]byte{}, [][]byte{}}
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
