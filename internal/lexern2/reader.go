package lexern2

import "strings"

type Reader struct {
	*strings.Reader
	Line  int
	Lines []int
	Char  int
	Index int
}

func NewReader(content string) *Reader {
	return &Reader{
		Reader: strings.NewReader(content),
		Line:   1,
		Lines:  []int{},
		Char:   0,
		Index:  0,
	}
}

func (r *Reader) ReadRune() (rune, int, error) {
	_r, size, err := r.Reader.ReadRune()
	if err != nil {
		return _r, size, err
	}

	r.Index++
	if _r == '\n' {
		r.Line++
		r.Lines = append(r.Lines, r.Char)
		r.Char = 0
	} else {
		r.Char++
	}
	return _r, size, err
}

func (r *Reader) UnreadRune() error {
	err := r.Reader.UnreadRune()
	if err != nil {
		return err
	}
	r.Index--
	if r.Char == 0 {
		r.Line--
		r.Char = r.Lines[len(r.Lines)-1]
		r.Lines = r.Lines[:len(r.Lines)-1]
	} else {
		r.Char--
	}
	return nil
}
