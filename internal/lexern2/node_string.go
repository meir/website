package lexern2

import "io"

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
			break
		} else if err != nil {
			panic(err)
		}

		if r == n.StringChar {
			break
		}

		n.Content += string(r)
	}
	return nil
}

func (n *NodeString) String(p *Page, content NodeInterface, args ...string) string {
	cont := ""
	for _, node := range n.Nodes {
		cont += node.String(p, content, args...)
	}
	return cont + n.Content
}

func (n *NodeString) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err != nil && err.Error() != "EOF" {
		panic(err)
	}

	if r == '`' {
		n.StringChar = r
		return true, nil
	}

	p.Reader.UnreadRune()
	return false, nil
}
