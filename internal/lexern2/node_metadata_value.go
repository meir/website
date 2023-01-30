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

		r, _, err := p.Reader.ReadRune()
		if err != nil {
			p.Err(err)
		}

		switch r {
		case ';':
			n.Nodes = append(n.Nodes, &NodeRune{
				Content: string(n.Content),
			})
			return nil
		default:
			n.Content += string(r)
		}
	}
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
