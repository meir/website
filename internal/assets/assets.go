package assets

import (
	"bytes"
	"image"
	"image/gif"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/nickalie/go-webpbin"
	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	"golang.org/x/image/webp"
)

var min = minify.New()

func init() {
	min.AddFunc("text/css", css.Minify)
	min.AddFunc("text/javascript", js.Minify)
}

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
		return handleImage(file, buildfolder, content)
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

func handleImage(file, buildfolder string, data []byte) error {
	reader := bytes.NewReader(data)

	var err error
	var img image.Image

	mimetype := http.DetectContentType(data)

	switch mimetype {
	case "image/png":
		img, err = png.Decode(reader)
	case "image/jpeg":
		img, err = jpeg.Decode(reader)
	case "image/gif":
		img, err = gif.Decode(reader)
	case "image/webp":
		img, err = webp.Decode(reader)
	default:
		panic("unknown image type: " + mimetype)
	}

	if err != nil {
		panic(err)
	}

	// convert to webp
	// write to build folder/assets
	filename := strings.TrimSuffix(filepath.Base(file), filepath.Ext(file))
	f, err := os.Create(filepath.Join(buildfolder, "assets", filename+".webp"))
	if err != nil {
		panic(err)
	}
	defer f.Close()

	err = webpbin.Encode(f, img)
	if err != nil {
		panic(err)
	}

	return err
}
