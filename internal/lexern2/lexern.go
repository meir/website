package lexern2

import "strings"

type NodeType int

const (
	RAW NodeType = iota
)

type NodeInterface interface {
	String(p *Page, content NodeInterface, args ...string) string
	InternalNodes() []NodeInterface
	Process(*Page) error
	Detect(*Page) (bool, error)
}

type Page struct {
	Lexer *Lexer

	RootNode  NodeInterface
	Root      string
	Alias     string
	Src       string
	Content   string
	Metadata  map[string]NodeInterface
	Arguments map[string]string
	Reader    strings.Reader
}

type PageInterface interface {
	Process() *Page
	String() string
}

type Lexer struct {
	Root   string
	Pages  []*Page
	Global map[string]NodeInterface
	States map[string]bool
}

type FileProcessOptions struct {
	Root      string
	File      string
	Content   string
	Arguments map[string]string
	Lexer     *Lexer
}

type LexerInterface interface {
	LoadFile(options FileProcessOptions)
	Process() *Page
	String(*Page) string
	GetByMetaKey(key string) []*Page
	GetByMetaValue(key string, value string) []*Page
	GetPageByAlias(alias string) *Page
}

func NewLexer(root string) *Lexer {
	return &Lexer{
		Root:   root,
		Pages:  []*Page{},
		Global: map[string]NodeInterface{},
		States: map[string]bool{},
	}
}

func (l *Lexer) SetGlobal(key string, value string) {
	l.Global[key] = &NodeRune{
		Content: value,
	}
}
