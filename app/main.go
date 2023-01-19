package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/flamingo-development/static/internal/lexern"
)

var folder = flag.String("i", "./site", "Folder to build from")
var output = flag.String("o", "./build", "Folder to build to")

func main() {
	flag.Parse()

	lexer := lexern.NewLexer(*folder)
	os.MkdirAll(path.Join(*output, "/assets"), 0755)

	pages, err := load_dir(*folder, lexer)
	if err != nil {
		panic(err)
	}

	//write files into output folder in folders named after the url
	for _, page := range pages {
		if page.Meta["url"] != "" {
			os.MkdirAll(path.Join(*output, page.Meta["url"]), 0755)
			err := os.WriteFile(path.Join(
				*output,
				page.Meta["url"],
				"index.html",
			), []byte(page.Content), 0644)
			if err != nil {
				panic(err)
			}
		}
	}

}

func load_dir(f string, lexer *lexern.Lexer) (map[string]*lexern.Page, error) {
	dir, err := os.ReadDir(f)
	if err != nil {
		panic(err)
	}

	pages := make(map[string]*lexern.Page)

	for _, file := range dir {
		if file.IsDir() {
			subpages, err := load_dir(path.Join(f, file.Name()), lexer)
			if err != nil {
				panic(err)
			}

			for url, page := range subpages {
				pages[url] = page
			}
			continue
		}

		if !strings.HasSuffix(file.Name(), ".htm") {
			_, err := copy(path.Join(f, file.Name()), path.Join(*output, "assets", file.Name()))
			if err != nil {
				panic(err)
			}
			continue
		}

		page, err := lexer.ProcessFile(f, file.Name())
		if err != nil {
			panic(err)
		}

		if page.Meta["url"] != "" {
			pages[page.Meta["url"]] = page
		}
	}

	return pages, nil
}

func copy(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}
