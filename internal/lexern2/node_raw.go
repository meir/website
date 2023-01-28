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
		&NodeEscape{},
		&NodeString{},
		&NodeMetadata{},
		&NodeInsert{},
		&NodeTag{},
	}
}

func (n *NodeRaw) Process(p *Page) error {
	for {
		nodes, new, err := ScanContent(n, p)
		if err != nil {
			panic(err)
		}

		if new {
			runeNode := NodeRune{
				Content: string(n.Content),
			}
			n.Nodes = append(n.Nodes, &runeNode)
			n.Content = ""
			n.Nodes = append(n.Nodes, nodes...)
			continue
		}

		r, _, err := p.Reader.ReadRune()
		if err == io.EOF {
			n.Nodes = append(n.Nodes, &NodeRune{
				Content: string(n.Content),
			})
			break
		} else if err != nil {
			panic(err)
		}

		if !n.Root && r == '}' {
			break
		}

		n.Content += string(r)
	}
	return nil
}

func (n *NodeRaw) String(p *Page, content NodeInterface, args ...string) string {
	cont := ""
	for _, node := range n.Nodes {
		cont += node.String(p, content)
	}
	return strings.TrimSpace(cont)

}

func (n *NodeRaw) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	if r == ':' {
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
