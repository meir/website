package assets

import (
	"os"
	"path/filepath"

	"github.com/tdewolff/minify"
)

var min = minify.New()

func HandleFile(file string, buildfolder string) error {
	// check file type
	// minimize css files
	// convert images to webp
	// minimize js files
	ext := filepath.Ext(file)
	filename := filepath.Base(file)

	content, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	switch ext {
	case ".css":
		content, err = min.Bytes("text/css", content)
		if err != nil {
			return err
		}
	case ".js":
		content, err = min.Bytes("text/javascript", content)
		if err != nil {
			return err
		}
	case ".png", ".jpg", ".jpeg", ".gif":

	}

	// write file to build folder/assets
	err = os.WriteFile(
		filepath.Join(buildfolder, "assets", filename),
		content,
		0644,
	)

	if err != nil {
		return err
	}

	return nil
}
