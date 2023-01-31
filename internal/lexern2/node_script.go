package lexern2

import "io"

type NodeScript struct {
	Content string
}

func (n *NodeScript) InternalNodes() []NodeInterface {
	return []NodeInterface{}
}

func (n *NodeScript) Process(p *Page) error {
	for {
		r, err := p.Reader.ReadRune()
		if err != nil {
			p.Err(err)
		}

	RUNE_SWITCH:
		switch r {
		case '<':
			//check if it is a closing tag
			tag := "/script>"
			n.Content += "<"
			for i := 0; i < len(tag); i++ {
				r, err := p.Reader.ReadRune()
				if err != nil && err != io.EOF {
					return err
				}

				if err != io.EOF {
					n.Content += string(r)
				}

				if r != rune(tag[i]) {
					break RUNE_SWITCH
				}

				return nil
			}

		default:
			n.Content += string(r)
		}
	}
}

func (n *NodeScript) String(p *Page, content NodeInterface, args ...string) string {
	return string(n.Content)
}

func (n *NodeScript) Detect(p *Page) (bool, error) {
	tag := []rune("<script")
	for index, tag_rune := range tag {
		r, err := p.Reader.Peek(index)
		if err != nil && err != io.EOF {
			p.Err(err)
		}

		if r != tag_rune {
			return false, nil
		}
	}
	return true, nil
}
