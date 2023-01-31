package lexern2

import (
	"errors"
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
		node, err := ScanContent(n, p)
		if err != nil {
			p.Err(err)
		}

		if node != nil {
			n.Token = strings.TrimSpace(n.Token)
			n.Token = strings.Split(n.Token, ":")[0]
			n.Content = node
			return nil
		}

		r, err := p.Reader.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			p.Err(err)
		}

		switch r {
		case '}':
			n.Token = strings.TrimSpace(n.Token)
			return nil
		default:
			n.Token += string(r)
		}
	}
}

func (n *NodeTag) String(p *Page, content NodeInterface, args ...string) string {
	page := p.Lexer.GetPage(p.Root, n.Token)
	if page == nil {
		p.Err(errors.New("Page not found: " + n.Token))
	}

	for key, data := range p.Metadata {
		page.Metadata[key] = data
	}

	return page.String(p, n.Content, args...)
}

func (n *NodeTag) Detect(p *Page) (bool, error) {
	r, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		p.Err(err)
	}

	if r == '{' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
