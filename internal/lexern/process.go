package lexern

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (l *Lexer) process_file(root, filename, content string, history []*FileLexer) (*Page, error) {
	if strings.HasPrefix(filename, "$") {
		root = l.root
		filename = filename[1:]
	}

	if history == nil {
		history = []*FileLexer{}
	} else {
		for _, flh := range history {
			if flh.file == path.Join(root, filename) {
				return flh.page, nil
			}
		}
	}

	fl := &FileLexer{
		buffer:  &Buffer{},
		skip:    false,
		state:   RAW,
		lexer:   l,
		file:    path.Join(root, filename),
		history: history,
	}

	data, err := os.ReadFile(path.Join(root, filename))
	if err != nil {
		panic(err)
	}

	err = fl.process_char(root, data, content)
	if err != nil {
		return nil, err
	}

	return fl.page, nil
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

func (l *FileLexer) process_char(root string, file []byte, cont string) error {
	l.page = &Page{
		Meta:    make(map[string]string),
		Content: cont,
	}

	l.history = append(l.history, l)

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
				page, err := l.lexer.process_file(root, state.token, state.content, l.history)
				if err != nil {
					panic(err)
				}
				l.buffer.Add(page.Content)
			case '\\':
				l.debug("raw", "\\")
				l.skip = !l.skip
			case '"', '\'', '`':
				l.debug("raw", string(c))
				l.state = STRING
				l.buffer.Up()
				l.buffer.Current().token_value = string(c)
			default:
				l.buffer.AddC(c)
			}
		case STRING:
			switch c {
			case '"', '\'', '`':
				l.debug("string", string(c))
				if l.buffer.Current().token_value == string(c) {
					l.state = RAW
					state := l.buffer.Down()
					l.debug("buffer -", l.buffer.Last(), state.token, state.token_value)
					l.buffer.Add(state.content)
				} else {
					l.buffer.AddC(c)
				}
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
				page, err := l.lexer.process_file(root, state.token, "", l.history)
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
					if meta, ok := l.page.Meta[state.token]; ok {
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
				l.page.Meta[state.token] = state.token_value
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
				page, err := l.lexer.process_file(root, state.token, "", l.history)
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

	l.page.Content = l.buffer.String()

	return nil
}
