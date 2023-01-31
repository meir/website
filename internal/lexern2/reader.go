package lexern2

import (
	"io"
)

type Reader struct {
	Content []rune
	Line    int
	Lines   []int
	Char    int
	Index   int64
}

func NewReader(content string) *Reader {
	return &Reader{
		Content: []rune(content),
		Line:    1,
		Lines:   []int{},
		Char:    0,
		Index:   0,
	}
}

func (r *Reader) Peek(i int) (rune, error) {
	if r.Index+int64(i) >= int64(len(r.Content)) {
		return 0, io.EOF
	}
	return r.Content[r.Index+int64(i)], nil
}

func (r *Reader) ReadRune() (rune, error) {
	if r.Index >= int64(len(r.Content)) {
		return 0, io.EOF
	}

	ch := r.Content[r.Index]

	r.Index++
	if ch == '\n' {
		r.Line++
		r.Lines = append(r.Lines, r.Char)
		r.Char = 0
	} else {
		r.Char++
	}
	return ch, nil
}

func (r *Reader) UnreadRune() error {
	if r.Index <= 0 {
		return io.EOF
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
