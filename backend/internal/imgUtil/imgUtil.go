package imgUtil

import (
	"bytes"
	"errors"
	"image"
	_ "image/gif"
	"image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/devbydaniel/release-notes-go/internal/logger"
	"github.com/nfnt/resize"
)

var log = logger.Get()

type ImgProcessConfig struct {
	MaxWidth uint
	Quality  int
	Format   string
}

func VerifyImageType(img multipart.File) bool {
	log.Trace().Msg("VerifyImageType")
	// Read the first 512 bytes to detect content type
	buff := make([]byte, 512)
	_, err := img.Read(buff)
	if err != nil {
		return false
	}

	// Reset the file pointer
	img.Seek(0, io.SeekStart)

	// Verify file type
	contentType := http.DetectContentType(buff)
	if !strings.HasPrefix(contentType, "image/") {
		return false
	}
	return true
}

func Decode(file io.Reader) (image.Image, string, error) {
	log.Trace().Msg("Decode")
	// Decode the image
	img, format, err := image.Decode(file)
	if err != nil {
		return nil, "", err
	}
	return img, format, nil
}

func ProcessImage(img image.Image, maxWidth uint, quality int) (image.Image, error) {
	log.Trace().Msg("ProcessImage")
	// Resize the image if it's too wide
	bounds := img.Bounds()
	width := uint(bounds.Dx())
	if width > maxWidth {
		log.Debug().Uint("width", width).Uint("maxWidth", maxWidth).Msg("Resizing image")
		img = resize.Resize(maxWidth, 0, img, resize.Lanczos3)
	}
	return img, nil
}

func Encode(img image.Image, format string) (*io.Reader, error) {
	log.Trace().Msg("Encode")
	// Encode the image
	imgBuf := new(bytes.Buffer)
	switch format {
	case "jpeg":
		if err := jpeg.Encode(imgBuf, img, nil); err != nil {
			return nil, err
		}
	default:
		return nil, errors.New("unsupported image format")
	}

	ioReader := io.Reader(imgBuf)
	return &ioReader, nil
}

func DecodeProcessEncode(img io.Reader, config *ImgProcessConfig) (*io.Reader, error) {
	log.Trace().Msg("DecodeProcessEncode")
	// Decode the image
	imgDecoded, _, err := Decode(img)
	if err != nil {
		return nil, err
	}

	// Process the image
	imgProcessed, err := ProcessImage(imgDecoded, config.MaxWidth, config.Quality)
	if err != nil {
		return nil, err
	}

	// Encode the image
	imgEncoded, err := Encode(imgProcessed, config.Format)
	if err != nil {
		return nil, err
	}

	return imgEncoded, nil
}
