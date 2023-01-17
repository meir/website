package lexer

import (
	"os"
	"path"
	"strings"
)

type LexerState int

const (
	RAW LexerState = iota
	TAG
	CONTENT
	META
	META_VALUE
	META_GET
)

type Page struct {
	Meta    map[string]string
	Content string
}

var absolute_root string = ""

func SetRoot(root string) {
	absolute_root = root
}

func ProcessFile(root, file string) (*Page, error) {
	return process_file(root, file, "")
}

func process_file(root, file, cont string) (*Page, error) {
	if strings.HasPrefix(file, "$") {
		root = absolute_root
		file = file[1:]
	}

	data, err := os.ReadFile(path.Join(root, file))
	if err != nil {
		return nil, err
	}

	return process_char(root, data, cont)
}

var state = RAW

func process_char(root string, file []byte, cont string) (*Page, error) {
	var buffer string
	var content string
	var token string
	var token_value string
	var skip bool = false
	var metadata map[string]string = make(map[string]string)

	for i := 0; i < len(file); i++ {
		var c rune = rune(file[i])

		switch c {
		case '{':
			if skip {
				if state == CONTENT {
					content += string(c)
				} else {
					buffer += string(c)
				}
				skip = false
			} else {
				state = TAG
			}

		case '}':
			if skip {
				if state == CONTENT {
					content += string(c)
				} else {
					buffer += string(c)
				}
				skip = false
			} else if state == CONTENT || state == TAG {
				state = RAW
				page, err := process_file(root, token, content)
				if err != nil {
					return nil, err
				}
				buffer += page.Content
				token = ""
				token_value = ""
				content = ""
			}

		case '\\':
			skip = !skip

		case '+':
			if skip {
				if state == CONTENT {
					content += string(c)
				} else {
					buffer += string(c)
				}
				skip = false
			} else {
				state = META
			}

		case ';':
			if state == META || state == META_VALUE {
				state = RAW
				if token == "content" {
					buffer += cont
				} else {
					metadata[token] = token_value
				}
				token = ""
				token_value = ""
			} else if state == META_GET {
				state = RAW
				page, err := process_file(root, token, content)
				if err != nil {
					return nil, err
				}

				if meta, ok := page.Meta[token_value]; ok {
					buffer += meta
				}
				token = ""
				token_value = ""
			}

		case ' ':
			if state == CONTENT {
				content += string(c)
			} else {
				buffer += string(c)
			}

		case '=':
			if state == META {
				state = META_VALUE
			} else {
				if state == CONTENT {
					content += string(c)
				} else {
					buffer += string(c)
				}
			}

		case ':':
			if state == TAG {
				state = CONTENT
			} else if state == META {
				state = META_GET
			} else {
				if state == CONTENT {
					content += string(c)
				} else {
					buffer += string(c)
				}
			}

		default:
			switch state {
			case META:
				token += string(c)

			case META_VALUE, META_GET:
				token_value += string(c)

			case TAG:
				token += string(c)

			case CONTENT:
				content += string(c)

			default:
				buffer += string(c)

			}
		}
	}

	return &Page{
		Meta:    metadata,
		Content: buffer,
	}, nil
}
