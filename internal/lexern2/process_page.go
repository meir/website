package lexern2

func (p *Page) Process() {
	p.RootNode = &NodeRaw{
		Content: p.Content,
		Nodes:   []NodeInterface{},
		Root:    true,
	}

	p.RootNode.Process(p)
}

func (p *Page) String(_p *Page, content NodeInterface, args ...string) string {
	return p.RootNode.String(_p, content, args...)
}
