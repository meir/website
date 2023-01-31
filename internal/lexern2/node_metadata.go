package lexern2

import (
	"io"
	"strings"
)

type NodeMetadata struct {
	Token string
}

func (n *NodeMetadata) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeMetadata) Process(p *Page) error {
	for {
		r, err := p.Reader.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			p.Err(err)
		}

		switch r {
		case '=':
			value := &NodeMetadataValue{}
			value.Process(p)

			if strings.HasPrefix(n.Token, ".") {
				p.Lexer.Global[n.Token[1:]] = value
				return nil
			}

			switch n.Token {
			case "alias":
				p.Alias = value.String(p, nil)
			default:
				p.Metadata[n.Token] = value
			}
			return nil
		case ';':
			value := &NodeRune{Content: ""}
			if strings.HasPrefix(n.Token, ".") {
				p.Lexer.Global[n.Token[1:]] = value
				return nil
			}

			p.Metadata[n.Token] = value
			return nil
		case ' ':
			continue
		default:
			n.Token += string(r)
		}
	}
}

func (n *NodeMetadata) String(p *Page, content NodeInterface, args ...string) string {
	return ""
}

func (n *NodeMetadata) Detect(p *Page) (bool, error) {
	r, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		p.Err(err)
	}

	if r == '#' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
