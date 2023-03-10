package lexern2

import (
	"io"
	"strings"
)

type NodeRaw struct {
	Content string
	Nodes   []NodeInterface
	Root    bool
}

func (n *NodeRaw) InternalNodes() []NodeInterface {
	return []NodeInterface{
		&NodeScript{},
		&NodeEscape{},
		&NodeString{},
		&NodeMetadata{},
		&NodeInsert{},
		&NodeEach{},
		&NodeTag{},
	}
}

func (n *NodeRaw) Process(p *Page) error {
	for {
		node, err := ScanContent(n, p)
		if err != nil {
			p.Err(err)
		}

		if node != nil {
			runeNode := NodeRune{
				Content: string(n.Content),
			}
			n.Nodes = append(n.Nodes, &runeNode)
			n.Content = ""
			n.Nodes = append(n.Nodes, node)
			continue
		}

		r, err := p.Reader.ReadRune()
		if err == io.EOF {
			n.Nodes = append(n.Nodes, &NodeRune{
				Content: string(n.Content),
			})
			return nil
		} else if err != nil {
			p.Err(err)
		}

		if !n.Root && r == '}' {
			n.Nodes = append(n.Nodes, &NodeRune{
				Content: string(n.Content),
			})
			return nil
		}

		n.Content += string(r)
	}
}

func (n *NodeRaw) String(p *Page, content NodeInterface, args ...string) string {
	cont := ""
	for _, node := range n.Nodes {
		cont += node.String(p, content)
	}
	return strings.TrimSpace(cont)

}

func (n *NodeRaw) Detect(p *Page) (bool, error) {
	r, err := p.Reader.ReadRune()
	if err != nil && err.Error() != "EOF" {
		p.Err(err)
	}

	if r == ':' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
