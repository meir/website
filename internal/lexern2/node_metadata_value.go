package lexern2

import "strings"

type NodeMetadataValue struct {
	Content string
	Nodes   []NodeInterface
}

func (n *NodeMetadataValue) InternalNodes() []NodeInterface {
	return []NodeInterface{
		&NodeEscape{},
		&NodeString{},
	}
}

func (n *NodeMetadataValue) Process(p *Page) error {
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
		if err != nil {
			panic(err)
		}

		if r == ';' {
			n.Nodes = append(n.Nodes, &NodeRune{
				Content: string(n.Content),
			})
			break
		}

		n.Content += string(r)
	}
	return nil
}

func (n *NodeMetadataValue) String(p *Page, content NodeInterface, args ...string) string {
	cont := ""
	for _, node := range n.Nodes {
		cont += node.String(p, content)
	}
	return strings.TrimSpace(cont)
}

func (n *NodeMetadataValue) Detect(p *Page) (bool, error) {
	return true, nil
}
