package lexern2

import (
	"io"
	"path"
	"strings"
)

type NodeInsert struct {
	token      string
	token_name string
}

func (n *NodeInsert) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeInsert) Process(p *Page) error {
	for {
		r, _, err := p.Reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if r == ';' {
			line := strings.Split(n.token, ":")
			if len(line) >= 2 {
				n.token = line[0]
				n.token_name = strings.Join(line[1:], ":")
			}
			break
		}

		n.token += string(r)
	}
	return nil
}

func (n *NodeInsert) String(p *Page, content NodeInterface, args ...string) string {
	if n.token == "content" {
		if content == nil {
			panic("Content not found")
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
	panic("File not found: " + n.token)
}

func (n *NodeInsert) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		panic(err)
	}

	if r == '+' {
		return true, nil
	}
	p.Reader.UnreadRune()

	return false, nil
}
