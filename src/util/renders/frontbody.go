package renders

import (
	"image"
)

func RenderFrontBody(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	slimOffset := GetSlimOffset(opts.Slim)

	var (
		frontHead  *image.NRGBA = RemoveTransparency(Extract(skin, 8, 8, 8, 8))
		frontTorso *image.NRGBA = RemoveTransparency(Extract(skin, 20, 20, 8, 12))
		leftArm    *image.NRGBA = nil
		rightArm   *image.NRGBA = RemoveTransparency(Extract(skin, 44, 20, 4-slimOffset, 12))
		leftLeg    *image.NRGBA = nil
		rightLeg   *image.NRGBA = RemoveTransparency(Extract(skin, 4, 20, 4, 12))
	)

	if IsOldSkin(skin) {
		leftArm = FlipHorizontal(rightArm)
		leftLeg = FlipHorizontal(rightLeg)
	} else {
		leftArm = RemoveTransparency(Extract(skin, 36, 52, 4-slimOffset, 12))
		leftLeg = RemoveTransparency(Extract(skin, 20, 52, 4, 12))

		if opts.Overlay {
			overlaySkin := FixTransparency(skin)

			frontHead = Composite(frontHead, Extract(overlaySkin, 40, 8, 8, 8), 0, 0)
			frontTorso = Composite(frontTorso, Extract(overlaySkin, 20, 36, 8, 12), 0, 0)
			leftArm = Composite(leftArm, Extract(overlaySkin, 52, 52, 4-slimOffset, 64), 0, 0)
			rightArm = Composite(rightArm, Extract(overlaySkin, 44, 36, 4-slimOffset, 48), 0, 0)
			leftLeg = Composite(leftLeg, Extract(overlaySkin, 4, 52, 4, 12), 0, 0)
			rightLeg = Composite(rightLeg, Extract(overlaySkin, 4, 36, 4, 12), 0, 0)
		}
	}

	output := image.NewNRGBA(image.Rect(0, 0, 16, 32))

	// Face
	output = Composite(output, frontHead, 4, 0)

	// Torso
	output = Composite(output, frontTorso, 4, 8)

	// Left Arm
	output = Composite(output, leftArm, 12, 8)

	// Right Arm
	output = Composite(output, rightArm, slimOffset, 8)

	// Left Leg
	output = Composite(output, leftLeg, 8, 20)

	// Right Leg
	output = Composite(output, rightLeg, 4, 20)

	return Scale(output, opts.Scale)
}
