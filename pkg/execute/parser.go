package execute

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
)

// Parser represents a parser.
type Parser struct {
	scanner *Scanner
	buf     struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{scanner: NewScanner(r)}
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.scanner.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }

func (p *Parser) parseKeyValues(kv map[string]string) (out map[string]string, err error) {
	token, key := p.scan()
	if token != WORD {
		return nil, errors.New("KV need key")
	}
	token, _ = p.scan()
	if token != COLON {
		return nil, errors.New("KV need :")
	}
	token, value := p.scan()
	if token == WORD || token == NUMBER {
		kv[key] = value
	} else {
		return nil, errors.New("KV needs value")
	}

	token, _ = p.scan()
	if token == COMMA {
		return p.parseKeyValues(kv)
	} else {
		p.unscan()
	}
	return kv, nil
}

func (p *Parser) collectBlockContent() (content string, err error) {
	var buffer bytes.Buffer

	for depth := 1; depth > 0;
	{
		token, str := p.scan()
		if token == BLOCK_START {
			depth = depth + 1
			buffer.Write([]byte(str))
		} else if token == BLOCK_END {
			depth = depth - 1
			if depth > 0 {
				buffer.Write([]byte(str))
			}
		} else {
			buffer.Write([]byte(str))
		}
	}
	return buffer.String(), nil
}

func (p *Parser) parseBlock() (plan Plan, err error) {
	times := 1
	mode := "s"

	token, val := p.scan()

	if token == NUMBER {
		times, _ = strconv.Atoi(val)
		token, val = p.scan()
	} else {
		times = 1
	}
	if token == WORD {
		if val == "p" || val == "s" {
			mode = val
		} else {
			return nil, errors.New("Only s or p accepted")
		}

		token, val = p.scan()
	}
	var content string
	if token == BLOCK_START {
		content, err = p.collectBlockContent()
		if err != nil {
			errors.New("")
		}
		_ = content
		token, val = p.scan()
	} else {
		return nil, errors.New("Unknown token for Block")
	}

	var kv map[string]string
	if token == COMMA {
		kv, err = p.parseKeyValues(make(map[string]string))
		if err != nil {
			return nil, err
		}
		token, _ = p.scan()
	} else {
		kv = make(map[string]string)
	}

	block := &Block{
		times: times,
		mode:  mode,
		Block: content,
		KV:    kv,
	}

	if token == PLUS || token == MULTIPLY {
		next, err := p.parseBlock()
		if err != nil {
			return nil, err
		}

		return &Operator{
			operand: token,
			left:    block,
			right:   next,
		}, nil
	}

	return block, nil

}

func (p *Parser) parseBlockContent() (plan *Block, err error) {
	token, val := p.scan()

	var content string
	if token != WORD {
		return nil, errors.New("Only s or p accepted")
	}
	content = val
	token, val = p.scan()

	var kv map[string]string
	if token == COMMA {
		kv, err = p.parseKeyValues(make(map[string]string))
		if err != nil {
			return nil, err
		}
		token, _ = p.scan()
	} else {
		kv = make(map[string]string)
	}

	block := &Block{
		times: 1,
		mode:  "s",
		Block: content,
		KV:    kv,
	}

	return block, nil

}

func Parse(molecule string) (plan Plan, err error) {
	parser := NewParser(strings.NewReader(molecule))
	return parser.parseBlock()
}

func ParseBlock(block string) (*Block, error) {
	parser := NewParser(strings.NewReader(block))
	return parser.parseBlockContent()
}