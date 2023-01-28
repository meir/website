package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"strings"
	"time"

	"github.com/flamingo-development/static/internal/lexern2"
)

var folder = flag.String("i", "./site", "Folder to build from")
var output = flag.String("o", "./build", "Folder to build to")

func main() {
	flag.Parse()

	lexer := lexern2.NewLexer(*folder)
	os.MkdirAll(path.Join(*output, "/assets"), 0755)

	lexer.SetGlobal("build_date", time.Now().Format("2006-01-02 15:04:05"))
	lexer.SetGlobal("year", time.Now().Format("2006"))
	lexer.SetGlobal("version", "0.0.1")
	lexer.SetGlobal("time", time.Now().Format("15:04:05"))
	lexer.SetGlobal("go_version", runtime.Version())
	lexer.SetGlobal("go_os", runtime.GOOS)
	lexer.SetGlobal("go_arch", runtime.GOARCH)

	err := load_dir(*folder, lexer)
	if err != nil {
		panic(err)
	}

	lexer.Process()

	pages := lexer.GetByMetaKey("url")
	for _, page := range pages {
		filepath := page.Metadata["url"].String(page, nil)
		content := page.String(page, nil)
		p := path.Join(*output, filepath)
		os.MkdirAll(p, 0755)
		fmt.Printf("----- %s ------\n%s\n====================\n", path.Join(p, "index.html"), content)
		err := os.WriteFile(
			path.Join(p, "index.html"),
			[]byte(content),
			0755,
		)
		if err != nil {
			panic(err)
		}
	}
}

func load_dir(f string, lexer *lexern2.Lexer) error {
	dir, err := os.ReadDir(f)
	if err != nil {
		panic(err)
	}

	for _, file := range dir {
		if file.IsDir() {
			err := load_dir(path.Join(f, file.Name()), lexer)
			if err != nil {
				panic(err)
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

		lexer.LoadFile(
			lexern2.FileProcessOptions{
				Root: *folder,
				File: path.Join(f, file.Name()),
			},
		)
	}

	return nil
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
