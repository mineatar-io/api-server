package renders

import (
	"image"
	"image/draw"

	"github.com/mineatar-io/api-server/src/util/renders/matrix"
	drw "golang.org/x/image/draw"
	"golang.org/x/image/math/f64"
)

var (
	SkewA            float64       = 26.0 / 45.0
	SkewB            float64       = SkewA * 2.0
	TransformForward matrix.Matrix = matrix.Matrix{
		XX: 1, YX: -SkewA,
		XY: 0, YY: SkewB,
		X0: 0, Y0: SkewA,
	}
	TransformUp matrix.Matrix = matrix.Matrix{
		XX: 1, YX: -SkewA,
		XY: 1, YY: SkewA,
		X0: 0, Y0: 0,
	}
	TransformRight matrix.Matrix = matrix.Matrix{
		XX: 1, YX: SkewA,
		XY: 0, YY: SkewB,
		X0: 0, Y0: 0,
	}
)

func Extract(img *image.NRGBA, x, y, width, height int) *image.NRGBA {
	output := image.NewNRGBA(image.Rect(0, 0, width, height))

	draw.Draw(output, output.Bounds(), img, image.Pt(x, y), draw.Src)

	return output
}

func Scale(img *image.NRGBA, scale int) *image.NRGBA {
	if scale == 1 {
		return img
	}

	bounds := img.Bounds().Max
	output := image.NewNRGBA(image.Rect(0, 0, bounds.X*scale, bounds.Y*scale))

	for x := 0; x < bounds.X; x++ {
		for y := 0; y < bounds.Y; y++ {
			color := img.At(x, y)

			for sx := 0; sx < scale; sx++ {
				for sy := 0; sy < scale; sy++ {
					output.Set(x*scale+sx, y*scale+sy, color)
				}
			}
		}
	}

	return output
}

func RemoveTransparency(img *image.NRGBA) *image.NRGBA {
	output := image.NewNRGBA(img.Bounds())

	for i, l := 0, len(img.Pix); i < l; i += 4 {
		output.Pix[i] = img.Pix[i]
		output.Pix[i+1] = img.Pix[i+1]
		output.Pix[i+2] = img.Pix[i+2]
		output.Pix[i+3] = 255
	}

	return output
}

func IsOldSkin(img *image.NRGBA) bool {
	return img.Bounds().Max.Y < 64
}

func Composite(bottom, top *image.NRGBA, x, y int) *image.NRGBA {
	output := image.NewNRGBA(bottom.Bounds())

	topBounds := top.Bounds().Max

	draw.Draw(output, bottom.Bounds(), bottom, image.Pt(0, 0), draw.Src)
	draw.Draw(output, image.Rect(0, 0, topBounds.X+x, topBounds.Y+y), top, image.Pt(-x, -y), draw.Over)

	return output
}

func FlipHorizontal(img *image.NRGBA) *image.NRGBA {
	data := img.Pix
	bounds := img.Bounds()

	output := image.NewNRGBA(bounds)

	for x := 0; x < bounds.Max.X; x++ {
		for y := 0; y < bounds.Max.Y; y++ {
			fx := bounds.Max.X - x - 1
			fi := fx*4 + y*4*bounds.Max.X
			ix := x*4 + y*4*bounds.Max.X

			for i := 0; i < 4; i++ {
				output.Pix[ix+i] = data[fi+i]
			}
		}
	}

	return output
}

func FixTransparency(img *image.NRGBA) *image.NRGBA {
	a := img.Pix[0:4]

	if a[3] == 0 {
		return img
	}

	output := Clone(img)

	for i, l := 0, len(output.Pix); i < l; i += 4 {
		if output.Pix[i+0] != a[0] || output.Pix[i+1] != a[1] || output.Pix[i+2] != a[2] || output.Pix[i+3] != a[3] {
			continue
		}

		output.Pix[i+3] = 0
	}

	return output
}

func Clone(img *image.NRGBA) *image.NRGBA {
	bounds := img.Bounds()
	output := image.NewNRGBA(bounds)

	draw.Draw(output, bounds, img, image.Pt(0, 0), draw.Src)

	return output
}

func GetSlimOffset(slim bool) int {
	if slim {
		return 1
	}

	return 0
}

func CompositeTransform(bottom, top *image.NRGBA, mat matrix.Matrix, x, y float64) *image.NRGBA {
	output := image.NewNRGBA(bottom.Bounds())

	draw.Draw(output, bottom.Bounds(), bottom, image.Pt(0, 0), draw.Src)

	transformer := drw.NearestNeighbor

	fx, fy := float64(x), float64(y)

	m := mat.Translate(fx, fy)

	transformer.Transform(output, f64.Aff3{m.XX, m.XY, m.X0, m.YX, m.YY, m.Y0}, top, top.Bounds(), draw.Over, nil)

	return output
}
