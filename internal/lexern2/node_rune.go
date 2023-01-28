package lexern2

type NodeRune struct {
	Content string
}

func ScanContent(n NodeInterface, p *Page) ([]NodeInterface, bool, error) {
	nodes := []NodeInterface{}

	for _, node := range n.InternalNodes() {
		if ok, err := node.Detect(p); err != nil {
			panic(err)
		} else if ok {
			nodes = append(nodes, node)
			node.Process(p)
			break
		}
	}

	return nodes, len(nodes) > 0, nil
}

func (n *NodeRune) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeRune) Process(p *Page) error {
	r, _, err := p.Reader.ReadRune()
	if err != nil {
		panic(err)
	}

	n.Content += string(r)
	return nil
}

func (n *NodeRune) String(p *Page, content NodeInterface, args ...string) string {
	return string(n.Content)
}

func (n *NodeRune) Detect(p *Page) (bool, error) {
	return true, nil
}
