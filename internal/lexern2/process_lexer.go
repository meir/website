package lexern2

import (
	"fmt"
	"os"
	"path"
	"strings"
)

func (l *Lexer) LoadFile(options FileProcessOptions) {
	content, err := os.ReadFile(options.File)
	if err != nil {
		panic(err)
	}
	page := &Page{
		Lexer:    l,
		Root:     options.Root,
		Alias:    options.File,
		Src:      options.File,
		Metadata: map[string]NodeInterface{},
		Reader:   *strings.NewReader(string(content)),
	}
	l.Pages = append(l.Pages, page)
}

func (l *Lexer) Process() {
	for _, page := range l.Pages {
		fmt.Println("Processing", page.Src)
		page.Process()
	}
	return
}

func (l *Lexer) String(p *Page) string {
	return p.String(p, nil)
}

func (l *Lexer) GetByMetaKey(key string) []*Page {
	pages := []*Page{}
	for _, page := range l.Pages {
		if _, ok := page.Metadata[key]; ok {
			pages = append(pages, page)
		}
	}
	return pages
}

func (l *Lexer) GetByMetaValue(key string, value string) []*Page {
	pages := []*Page{}
	for _, page := range l.Pages {
		if page.Metadata[key].String(page, nil) == value {
			pages = append(pages, page)
		}
	}
	return pages
}

func (l *Lexer) GetPage(root, s string) *Page {
	if p := l.GetPageByPath(root, s); p != nil {
		return p
	} else if p := l.GetPageByAlias(s); p != nil {
		return p
	}
	return nil
}

func (l *Lexer) GetPageByPath(root, p string) *Page {
	if strings.HasPrefix(p, "$") {
		if strings.HasPrefix(l.Root, "./") {
			l.Root = strings.Replace(l.Root, "./", "", 1)
		}
		if strings.HasPrefix(l.Root, "/") {
			l.Root = strings.Replace(l.Root, "/", "", 1)
		}
		p = strings.Replace(p, "$", l.Root, 1)
	} else {
		p = path.Join(root, p)
	}

	for _, page := range l.Pages {
		if strings.HasPrefix(page.Src, "./") {
			page.Src = strings.Replace(page.Src, "./", "", 1)
		}
		if strings.HasPrefix(page.Src, "/") {
			page.Src = strings.Replace(page.Src, "/", "", 1)
		}

		if page.Src == p {
			return page
		}
	}
	return nil
}

func (l *Lexer) GetPageByAlias(alias string) *Page {
	for _, page := range l.Pages {
		if page.Alias == alias {
			return page
		}
	}
	return nil
}
