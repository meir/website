package lexern2

import (
	"io"
	"strings"
)

type NodeTag struct {
	Token   string
	Nodes   []NodeInterface
	Content NodeInterface
}

func (n *NodeTag) InternalNodes() []NodeInterface {
	return []NodeInterface{
		&NodeRaw{},
	}
}

func (n *NodeTag) Process(p *Page) error {
	for {
		nodes, new, err := ScanContent(n, p)
		if err != nil {
			panic(err)
		}

		if new {
			n.Token = strings.TrimSpace(n.Token)
			n.Token = strings.Split(n.Token, ":")[0]
			n.Content = nodes[0]
			break
		}

		r, _, err := p.Reader.ReadRune()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}

		switch r {
		case '}':
			n.Token = strings.TrimSpace(n.Token)
			return nil
		default:
			n.Token += string(r)
		}
	}
	return nil
}

func (n *NodeTag) String(p *Page, content NodeInterface, args ...string) string {
	page := p.Lexer.GetPage(p.Root, n.Token)
	if page == nil {
		panic("Page not found: " + n.Token)
	}

	return page.String(p, n.Content, args...)
}

func (n *NodeTag) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		panic(err)
	}

	if r == '{' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
