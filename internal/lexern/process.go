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

func (l *FileLexer) process_char(root string, file []byte, cont string) error {
	l.page = &Page{
		Meta:    make(map[string]string),
		Content: cont,
	}

	l.history = append(l.history, l)

	for _, c := range string(file) {
		if l.skip {
			l.buffer.AddC(c)
			l.skip = false
			continue
		}
		charsinline := 0

		switch l.state {
		case RAW:
			switch c {
			case '{':
				l.state = TAG
				l.buffer.Up()
			case '+':
				l.state = META
				l.buffer.Up()
			case '}':
				state := l.buffer.Down()
				page, err := l.lexer.process_file(root, state.token, state.content, l.history)
				if err != nil {
					panic(err)
				}
				l.buffer.Add(page.Content)
			case '\\':
				l.skip = !l.skip
			case '"', '\'', '`':
				l.state = STRING
				l.buffer.Up()
				l.buffer.Current().token_value = string(c)
			case '\n':
				if charsinline > 0 {
					l.buffer.AddC(c)
				}
			case ' ':
				if charsinline > 0 {
					l.buffer.AddC(c)
				}
			default:
				l.buffer.AddC(c)
				charsinline++
			}
		case STRING:
			switch c {
			case '\\':
				l.skip = !l.skip
			case '"', '\'', '`':
				if l.buffer.Current().token_value == string(c) {
					l.state = RAW
					state := l.buffer.Down()
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
				l.state = RAW
				state := l.buffer.Down()
				page, err := l.lexer.process_file(root, state.token, "", l.history)
				if err != nil {
					panic(err)
				}
				l.buffer.Add(page.Content)
			case ':':
				l.state = RAW
			case ' ':
				break
			case '\n':
				break
			default:
				l.buffer.Current().token += string(c)
			}
		case META:
			switch c {
			case '=':
				l.state = META_VALUE
			case ':':
				l.state = META_GET
			case ';':
				state := l.buffer.Down()
				l.state = RAW
				if state.token == "content" {
					l.buffer.Add(cont)
				} else {
					if meta, ok := l.page.Meta[state.token]; ok {
						l.buffer.Add(meta)
					} else if meta, ok := l.lexer.global[state.token]; ok {
						l.buffer.Add(meta)
					} else {
						l.buffer.Add("(undefined: " + state.token + ")")
					}
				}
			case ' ':
				break
			case '\n':
				break
			default:
				l.buffer.Current().token += string(c)
			}
		case META_VALUE:
			switch c {
			case ';':
				state := l.buffer.Down()
				l.state = RAW
				if strings.HasSuffix(state.token_value, " ") {
					state.token_value = state.token_value[:len(state.token_value)-1]
				}
				if strings.HasPrefix(state.token, ".") {
					l.lexer.global[state.token[1:]] = state.token_value
				} else {
					l.page.Meta[state.token] = state.token_value
				}
			case ' ':
				value := l.buffer.Current().token_value
				if value == "" || value[len(value)-1] == ' ' {
					break
				}
				l.buffer.Current().token_value += string(c)
			case '\n':
				break
			default:
				l.buffer.Current().token_value += string(c)
			}
		case META_GET:
			switch c {
			case ';':
				state := l.buffer.Down()
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
			case '\n':
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
