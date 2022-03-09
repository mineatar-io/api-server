package renders

import (
	"image"
)

func RenderLeftBody(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	slimOffset := GetSlimOffset(opts.Slim)

	var (
		leftHead    *image.NRGBA = RemoveTransparency(Extract(skin, 24, 8, 8, 8))
		leftLeftArm *image.NRGBA = nil
		leftLeftLeg *image.NRGBA = nil
	)

	if IsOldSkin(skin) {
		leftLeftArm = FlipHorizontal(RemoveTransparency(Extract(skin, 40, 20, 4, 12)))
		leftLeftLeg = FlipHorizontal(RemoveTransparency(Extract(skin, 0, 20, 4, 12)))
	} else {
		leftLeftArm = RemoveTransparency(Extract(skin, 40-slimOffset, 52, 4, 12))
		leftLeftLeg = RemoveTransparency(Extract(skin, 24, 52, 4, 12))

		if opts.Overlay {
			overlaySkin := FixTransparency(skin)

			leftHead = Composite(leftHead, Extract(overlaySkin, 48, 8, 8, 8), 0, 0)
			leftLeftArm = Composite(leftLeftArm, Extract(overlaySkin, 56-slimOffset, 52, 4, 12), 0, 0)
			leftLeftLeg = Composite(leftLeftLeg, Extract(overlaySkin, 8, 52, 4, 12), 0, 0)
		}
	}

	output := image.NewNRGBA(image.Rect(0, 0, 8, 32))

	// Left Head
	output = Composite(output, leftHead, 0, 0)

	// Left Arm
	output = Composite(output, leftLeftArm, 2, 8)

	// Left Leg
	output = Composite(output, leftLeftLeg, 2, 20)

	return Scale(output, opts.Scale)
}
