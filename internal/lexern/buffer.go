package lexern

import "fmt"

type State struct {
	token       string
	token_value string
	content     string
}

func (s *State) String() string {
	return fmt.Sprintf("Token: %s, Value: %s, Content: %s", s.token, s.token_value, s.content)
}

type Buffer struct {
	Content []*State
}

func (b *Buffer) Add(s string) {
	b.Content[b.Last()].content += s
}

func (b *Buffer) AddC(c rune) {
	b.Content[b.Last()].content += string(c)
}

func (b *Buffer) Up() {
	b.Content = append(b.Content, &State{})
}

func (b *Buffer) Down() *State {
	state := b.Content[b.Last()]
	b.Content = b.Content[:len(b.Content)-1]
	return state
}

func (b *Buffer) Last() int {
	if len(b.Content) == 0 {
		b.Content = append(b.Content, &State{})
	}
	return len(b.Content) - 1
}

func (b *Buffer) String() string {
	return b.Content[b.Last()].content
}

func (b *Buffer) Current() *State {
	return b.Content[b.Last()]
}
