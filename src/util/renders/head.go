package renders

import "image"

func RenderHead(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	scale := float64(opts.Scale)
	output := image.NewNRGBA(image.Rect(0, 0, 16*opts.Scale, 19*opts.Scale-int(scale/2.0)-1))

	var (
		frontHead *image.NRGBA = RemoveTransparency(Extract(skin, 8, 8, 8, 8))
		topHead   *image.NRGBA = RemoveTransparency(Extract(skin, 8, 0, 8, 8))
		rightHead *image.NRGBA = RemoveTransparency(Extract(skin, 0, 8, 8, 8))
	)

	if opts.Overlay && !IsOldSkin(skin) {
		overlaySkin := FixTransparency(skin)

		frontHead = Composite(frontHead, Extract(overlaySkin, 40, 8, 8, 8), 0, 0)
		topHead = Composite(topHead, Extract(overlaySkin, 40, 0, 8, 8), 0, 0)
		rightHead = Composite(rightHead, Extract(overlaySkin, 32, 8, 8, 8), 0, 0)
	}

	// Front Head
	output = CompositeTransform(output, Scale(frontHead, opts.Scale), TransformForward, 8*scale, 12*scale-1)

	// Top Head
	output = CompositeTransform(output, Scale(topHead, opts.Scale), TransformUp, -4*scale, 4*scale)

	// Right Head
	output = CompositeTransform(output, Scale(rightHead, opts.Scale), TransformRight, 0, 4*scale)

	return output
}
