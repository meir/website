package lexern

type LexerState int

type Lexer struct {
	root   string
	global map[string]string
}

type FileLexer struct {
	buffer  *Buffer
	skip    bool
	state   LexerState
	lexer   *Lexer
	file    string
	page    *Page
	history []*FileLexer
}

const (
	RAW LexerState = iota
	TAG
	META
	META_VALUE
	META_GET
	STRING
)

type Page struct {
	Meta    map[string]string
	Content string
}

func NewLexer(root string) *Lexer {
	return &Lexer{
		root:   root,
		global: make(map[string]string),
	}
}

func (l *Lexer) ProcessFile(root, filename string) (*Page, error) {
	return l.process_file(root, filename, "", nil)
}
