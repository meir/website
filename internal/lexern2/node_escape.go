package lexern2

import "io"

type NodeEscape struct {
	Content rune
}

func (n *NodeEscape) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeEscape) Process(p *Page) error {
	r, _, err := p.Reader.ReadRune()
	if err != nil {
		panic(err)
	}

	n.Content = r
	return nil
}

func (n *NodeEscape) String(p *Page, content NodeInterface, args ...string) string {
	return string(n.Content)
}

func (n *NodeEscape) Detect(p *Page) (bool, error) {
	r, _, err := p.Reader.ReadRune()
	if err == io.EOF {
		return false, nil
	} else if err != nil {
		panic(err)
	}

	if r == '\\' {
		return true, nil
	}
	p.Reader.UnreadRune()
	return false, nil
}
