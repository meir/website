package lexern

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (l *Lexer) process_file(root, filename, content string) (*Page, error) {
	if strings.HasPrefix(filename, "$") {
		root = l.root
		filename = filename[1:]
	}

	fl := &FileLexer{
		buffer: &Buffer{},
		skip:   false,
		state:  RAW,
		lexer:  l,
		file:   path.Join(root, filename),
	}

	data, err := os.ReadFile(path.Join(root, filename))
	if err != nil {
		panic(err)
	}

	return fl.process_char(root, data, content)
}

func (l *FileLexer) debug(i ...interface{}) {
	if l.lexer.Debug {
		fmt.Println(i...)
	}
}

func (l *FileLexer) debugf(format string, i ...interface{}) {
	if l.lexer.Debug {
		fmt.Printf(format, i...)
	}
}

func (l *FileLexer) process_char(root string, file []byte, cont string) (*Page, error) {
	page := &Page{
		Meta:    make(map[string]string),
		Content: cont,
	}

	for _, c := range string(file) {
		if l.skip {
			l.debug("skip", c)
			l.buffer.AddC(c)
			l.skip = false
			continue
		}

		switch l.state {
		case RAW:
			switch c {
			case '{':
				l.debug("raw", "{")
				l.state = TAG
				l.buffer.Up()
				l.debug("buffer +", l.buffer.Last())
			case '+':
				l.debug("raw", "+")
				l.state = META
				l.buffer.Up()
				l.debug("buffer +", l.buffer.Last())
			case '}':
				l.debug("raw", "}")
				state := l.buffer.Down()
				l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
				page, err := l.lexer.process_file(root, state.token, state.content)
				if err != nil {
					panic(err)
				}
				l.buffer.Add(page.Content)
			default:
				l.buffer.AddC(c)
			}
		case TAG:
			switch c {
			case '}':
				l.debug("tag", "}")
				l.state = RAW
				state := l.buffer.Down()
				l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
				page, err := l.lexer.process_file(root, state.token, "")
				if err != nil {
					panic(err)
				}
				l.buffer.Add(page.Content)
			case ':':
				l.debug("tag", ":")
				l.state = RAW
			case ' ':
				break
			default:
				l.buffer.Current().token += string(c)
			}
		case META:
			switch c {
			case '=':
				l.debug("meta", "=")
				l.state = META_VALUE
			case ':':
				l.debug("meta", ":")
				l.state = META_GET
			case ';':
				l.debug("meta", ";")
				state := l.buffer.Down()
				l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
				l.state = RAW
				if state.token == "content" {
					l.buffer.Add(cont)
				} else {
					if meta, ok := page.Meta[state.token]; ok {
						l.buffer.Add(meta)
					} else {
						l.buffer.Add("(undefined: " + state.token + ")")
					}
				}
			case ' ':
				break
			default:
				l.buffer.Current().token += string(c)
			}
		case META_VALUE:
			switch c {
			case ';':
				l.debug("meta_value", ";")
				state := l.buffer.Down()
				l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
				l.state = RAW
				page.Meta[state.token] = state.token_value
			case ' ':
				break
			default:
				l.buffer.Current().token_value += string(c)
			}
		case META_GET:
			switch c {
			case ';':
				l.debug("meta_get", ";")
				state := l.buffer.Down()
				l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
				l.state = RAW
				page, err := l.lexer.process_file(root, state.token, "")
				if err != nil {
					panic(err)
				}
				if meta, ok := page.Meta[state.token_value]; ok {
					l.buffer.Add(meta)
				} else {
					l.buffer.Add("(undefined: " + state.token_value + ")")
				}
			case ' ':
				break
			default:
				l.buffer.Current().token_value += string(c)
			}
		default:
			panic("Unknown state")
		}
	}

	if l.buffer.Last() > 0 {
		panic(fmt.Errorf("Unclosed tag in %s", l.file))
	}

	page.Content = l.buffer.String()

	return page, nil
}
