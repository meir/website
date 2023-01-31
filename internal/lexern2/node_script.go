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
		r, _, err := p.Reader.ReadRune()
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
				r, _, err := p.Reader.ReadRune()
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
	tag := "<script"
	read := 0
	for i := 0; i < len(tag); i++ {
		r, _, err := p.Reader.ReadRune()
		if err != nil && err != io.EOF {
			return false, err
		}

		if err != io.EOF {
			read++
		}

		if r != rune(tag[i]) {
			//unread
			for i := read; i > 0; i-- {
				err := p.Reader.UnreadByte()
				if err != nil {
					return false, err
				}
			}
			return false, nil
		}
	}

	n.Content = tag
	return true, nil
}
