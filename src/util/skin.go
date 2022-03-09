package util

import (
	"bytes"
	_ "embed"
	"image"
	"image/draw"
	"image/png"
	"log"
)

var (
	//go:embed steve.png
	rawSteveData []byte

	//go:embed alex.png
	rawAlexData []byte

	steveSkin *image.NRGBA = image.NewNRGBA(image.Rect(0, 0, 64, 64))
	alexSkin  *image.NRGBA = image.NewNRGBA(image.Rect(0, 0, 64, 64))
)

func init() {
	rawSteveSkin, err := png.Decode(bytes.NewReader(rawSteveData))

	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(steveSkin, rawSteveSkin.Bounds(), rawSteveSkin, image.Pt(0, 0), draw.Src)

	rawAlexSkin, err := png.Decode(bytes.NewReader(rawAlexData))

	if err != nil {
		log.Fatal(err)
	}

	draw.Draw(alexSkin, rawAlexSkin.Bounds(), rawAlexSkin, image.Pt(0, 0), draw.Src)
}

func GetDefaultSkin(slim bool) *image.NRGBA {
	if slim {
		return alexSkin
	}

	return steveSkin
}
