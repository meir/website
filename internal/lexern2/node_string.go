package lexern2

import (
	"io"
	"strings"
)

type NodeString struct {
	Content    string
	Nodes      []NodeInterface
	StringChar rune
}

func (n *NodeString) InternalNodes() []NodeInterface {
	return []NodeInterface{
		&NodeEscape{},
	}
}

func (n *NodeString) Process(p *Page) error {
	for {
		node, err := ScanContent(n, p)
		if err != nil {
			p.Err(err)
		}

		if node != nil {
			runeNode := NodeRune{
				Content: strings.TrimSpace(n.Content),
			}
			n.Nodes = append(n.Nodes, &runeNode)
			n.Content = ""
			n.Nodes = append(n.Nodes, node)
			continue
		}

		r, err := p.Reader.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			p.Err(err)
		}

		if r == n.StringChar {
			return nil
		}

		n.Content += string(r)
	}
}

func (n *NodeString) String(p *Page, content NodeInterface, args ...string) string {
	cont := ""
	for _, node := range n.Nodes {
		cont += node.String(p, content, args...)
	}
	return cont + n.Content
}

func (n *NodeString) Detect(p *Page) (bool, error) {
	r, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		p.Err(err)
	}

	if r == '`' {
		n.StringChar = r
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
