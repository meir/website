package lexern2

import (
	"errors"
	"io"
	"strings"
)

type NodeEach struct {
	Token     string
	TokenName string
	Content   NodeInterface
}

func (n *NodeEach) InternalNodes() []NodeInterface {
	return []NodeInterface{
		&NodeRaw{},
	}
}

func (n *NodeEach) Process(p *Page) error {
	token_name := false
	for {
		node, err := ScanContent(n, p)
		if err != nil {
			p.Err(err)
		}

		if node != nil {
			n.Token = strings.TrimSpace(n.Token)
			n.TokenName = strings.TrimSpace(n.TokenName)

			n.Content = node
			return nil
		}

		r, err := p.Reader.ReadRune()
		if err == io.EOF {
			return nil
		} else if err != nil {
			p.Err(err)
		}

		switch r {
		case '}':
			p.Err(errors.New("Unexpected '}'; 'each' expects content!"))
		case '=':
			if !token_name {
				token_name = true
				continue
			}
			fallthrough

		default:
			if token_name {
				n.TokenName += string(r)
			} else {
				n.Token += string(r)
			}
		}
	}
}

func (n *NodeEach) String(p *Page, content NodeInterface, args ...string) string {
	pages := []*Page{}

	if n.TokenName == "" {
		pages = p.Lexer.GetByMetaKey(n.Token)
	} else {
		pages = p.Lexer.GetByMetaValue(n.Token, n.TokenName)
	}

	cont := ""
	for _, page := range pages {
		cont += n.Content.String(page, nil)
	}

	return cont
}

func (n *NodeEach) Detect(p *Page) (bool, error) {
	tag := []rune("{each ")
	for index, tag_rune := range tag {
		r, err := p.Reader.Peek(index)
		if err != nil && err != io.EOF {
			p.Err(err)
		}

		if r != tag_rune {
			return false, nil
		}
	}

	for i := 0; i < len(tag); i++ {
		p.Reader.ReadRune()
	}

	return true, nil
}
