package imgUtil

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	_ "image/gif"
	_ "image/jpeg"
	_ "image/png"
	"io"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/devbydaniel/announcable/internal/logger"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/nfnt/resize"
)

var log = logger.Get()

type SupportedFormat string

func (f SupportedFormat) String() string {
	return string(f)
}

const (
	JPEG SupportedFormat = "jpeg"
	PNG  SupportedFormat = "png"
)

type EncodedFormat string

func (f EncodedFormat) String() string {
	return string(f)
}

const (
	WebP EncodedFormat = "webp"
)

func (f SupportedFormat) ToEncodedFormat() EncodedFormat {
	switch f {
	case JPEG:
		return WebP
	case PNG:
		return WebP
	default:
		return ""
	}
}

type ImgProcessConfig struct {
	MaxWidth uint
	Quality  int
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
	log.Debug().Str("format", format).Msg("Image decoded")
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

func Encode(img image.Image, format SupportedFormat) (*io.Reader, error) {
	log.Trace().Msg("Encode")
	// Encode the image
	imgBuf := new(bytes.Buffer)
	targetFormat := format.ToEncodedFormat()
	switch targetFormat {
	case WebP:
		if format == PNG {
			// Use lossless for PNG (which often has transparency)
			options, err := encoder.NewLosslessEncoderOptions(encoder.PresetDefault, 6)
			if err != nil {
				return nil, err
			}
			if err := webp.Encode(imgBuf, img, options); err != nil {
				return nil, err
			}
			log.Debug().Msg("PNG image encoded as lossless WebP")
		} else {
			// For JPEG (no transparency), use lossy with high quality
			options, err := encoder.NewLossyEncoderOptions(encoder.PresetPhoto, 90)
			if err != nil {
				return nil, err
			}
			if err := webp.Encode(imgBuf, img, options); err != nil {
				return nil, err
			}
			log.Debug().Msg("JPEG image encoded as lossy WebP")
		}
	default:
		return nil, errors.New("unsupported image format")
	}

	ioReader := io.Reader(imgBuf)
	return &ioReader, nil
}

func DecodeProcessEncode(img io.Reader, config *ImgProcessConfig) (*io.Reader, EncodedFormat, error) {
	log.Trace().Msg("DecodeProcessEncode")
	// Decode the image
	imgDecoded, format, err := Decode(img)
	if err != nil {
		return nil, "", err
	}

	// check if format is supported
	if !isFormatSupported(format) {
		return nil, "", fmt.Errorf("unsupported image format: %s", format)
	}

	// Process the image
	imgProcessed, err := ProcessImage(imgDecoded, config.MaxWidth, config.Quality)
	if err != nil {
		return nil, "", err
	}

	// Encode the image
	imgEncoded, err := Encode(imgProcessed, SupportedFormat(format))
	if err != nil {
		return nil, "", err
	}

	return imgEncoded, SupportedFormat(format).ToEncodedFormat(), nil
}

func isFormatSupported(format string) bool {
	log.Trace().Str("format", format).Msg("isFormatSupported")
	switch SupportedFormat(format) {
	case JPEG, PNG:
		return true
	default:
		return false
	}
}
