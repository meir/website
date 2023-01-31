package lexern2

import (
	"fmt"
)

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

func (p *Page) Err(err error) {
	line, char := p.Reader.Line, p.Reader.Char
	fmt.Printf("%s:%d:%d: %s\n", p.Src, line, char, err.Error())
	panic(err)
}
