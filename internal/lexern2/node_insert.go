package lexern2

import (
	"errors"
	"io"
	"path"
)

type NodeInsert struct {
	token      string
	token_name string
}

func (n *NodeInsert) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeInsert) Process(p *Page) error {
	token_name := false
	for {
		r, _, err := p.Reader.ReadRune()
		if err != nil {
			p.Err(err)
		}

		switch r {
		case ' ':
			continue
		case ':':
			if !token_name {
				token_name = true
				continue
			}
			n.token_name += string(r)
		case ';':
			return nil
		default:
			if token_name {
				n.token_name += string(r)
			} else {
				n.token += string(r)
			}
		}
	}
}

func (n *NodeInsert) String(p *Page, content NodeInterface, args ...string) string {
	if n.token == "content" {
		if content == nil {
			p.Err(errors.New("Content not found"))
		}
		return content.String(p, content, args...)
	}

	if value, ok := p.Lexer.Global[n.token]; ok {
		return value.String(p, content)
	} else if value, ok := p.Metadata[n.token]; ok {
		return value.String(p, content)
	} else if sp := p.Lexer.GetPage(path.Dir(p.Src), n.token); sp != nil {
		if data, ok := sp.Metadata[n.token_name]; ok {
			return data.String(p, content)
		}
	}

	p.Err(errors.New("File not found: " + n.token))
	return ""
}

func (n *NodeInsert) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		p.Err(err)
	}

	if r == '+' {
		return true, nil
	}
	p.Reader.UnreadRune()

	return false, nil
}
