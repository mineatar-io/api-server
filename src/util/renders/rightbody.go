package renders

import (
	"image"
)

func RenderRightBody(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	var (
		rightHead     *image.NRGBA = RemoveTransparency(Extract(skin, 0, 8, 8, 8))
		rightRightArm *image.NRGBA = RemoveTransparency(Extract(skin, 40, 20, 4, 12))
		rightRightLeg *image.NRGBA = RemoveTransparency(Extract(skin, 0, 20, 4, 12))
	)

	if opts.Overlay && !IsOldSkin(skin) {
		overlaySkin := FixTransparency(skin)

		rightHead = Composite(rightHead, Extract(overlaySkin, 32, 8, 8, 8), 0, 0)
		rightRightArm = Composite(rightRightArm, Extract(overlaySkin, 40, 36, 4, 12), 0, 0)
		rightRightLeg = Composite(rightRightLeg, Extract(overlaySkin, 0, 36, 4, 12), 0, 0)
	}

	output := image.NewNRGBA(image.Rect(0, 0, 8, 32))

	// Right Head
	output = Composite(output, rightHead, 0, 0)

	// Right Arm
	output = Composite(output, rightRightArm, 2, 8)

	// Right Leg
	output = Composite(output, rightRightLeg, 2, 20)

	return Scale(output, opts.Scale)
}
