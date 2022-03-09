package renders

import "image"

func RenderFace(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	output := RemoveTransparency(Extract(skin, 8, 8, 8, 8))

	if opts.Overlay && !IsOldSkin(skin) {
		output = Composite(output, Extract(skin, 40, 8, 8, 8), 0, 0)
	}

	return Scale(output, opts.Scale)
}
