package util

import (
	"bytes"
	"image"
	"image/png"
)

func EncodePNG(img image.Image) ([]byte, error) {
	buf := &bytes.Buffer{}

	if err := png.Encode(buf, img); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
