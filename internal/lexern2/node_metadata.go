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
		r, _, err := p.Reader.ReadRune()
		if err != nil {
			panic(err)
		}

		if r == '=' {
			value := &NodeMetadataValue{}
			value.Process(p)

			if strings.HasPrefix(n.Token, ".") {
				p.Lexer.Global[n.Token[1:]] = value
				break
			}

			switch n.Token {
			case "alias":
				p.Alias = value.String(p, nil)
			default:
				p.Metadata[n.Token] = value
			}
			break
		}

		if r == ' ' {
			continue
		}

		n.Token += string(r)
	}
	return nil
}

func (n *NodeMetadata) String(p *Page, content NodeInterface, args ...string) string {
	return ""
}

func (n *NodeMetadata) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		panic(err)
	}

	if r == '#' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
