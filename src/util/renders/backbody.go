package renders

import (
	"image"
)

func RenderBackBody(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	slimOffset := GetSlimOffset(opts.Slim)

	var (
		backHead     *image.NRGBA = RemoveTransparency(Extract(skin, 24, 8, 8, 8))
		backTorso    *image.NRGBA = RemoveTransparency(Extract(skin, 32, 20, 8, 12))
		backLeftArm  *image.NRGBA = nil
		backRightArm *image.NRGBA = RemoveTransparency(Extract(skin, 52-slimOffset, 20, 4-slimOffset, 12))
		backLeftLeg  *image.NRGBA = nil
		backRightLeg *image.NRGBA = RemoveTransparency(Extract(skin, 12, 20, 4, 12))
	)

	if IsOldSkin(skin) {
		backLeftArm = FlipHorizontal(backRightArm)
		backLeftLeg = FlipHorizontal(backRightLeg)
	} else {
		backLeftArm = RemoveTransparency(Extract(skin, 44-slimOffset, 52, 4-slimOffset, 12))
		backLeftLeg = RemoveTransparency(Extract(skin, 28, 52, 4, 12))

		if opts.Overlay {
			overlaySkin := FixTransparency(skin)

			backHead = Composite(backHead, Extract(overlaySkin, 56, 8, 8, 8), 0, 0)
			backTorso = Composite(backTorso, Extract(overlaySkin, 32, 36, 8, 12), 0, 0)
			backLeftArm = Composite(backLeftArm, Extract(overlaySkin, 60-slimOffset, 52, 4-slimOffset, 64), 0, 0)
			backRightArm = Composite(backRightArm, Extract(overlaySkin, 52-slimOffset, 36, 4-slimOffset, 48), 0, 0)
			backLeftLeg = Composite(backLeftLeg, Extract(overlaySkin, 12, 52, 8, 64), 0, 0)
			backRightLeg = Composite(backRightLeg, Extract(overlaySkin, 12, 36, 8, 48), 0, 0)
		}
	}

	output := image.NewNRGBA(image.Rect(0, 0, 16, 32))

	// Face
	output = Composite(output, backHead, 4, 0)

	// Torso
	output = Composite(output, backTorso, 4, 8)

	// Left Arm
	output = Composite(output, backLeftArm, slimOffset, 8)

	// Right Arm
	output = Composite(output, backRightArm, 12, 8)

	// Left Leg
	output = Composite(output, backLeftLeg, 4, 20)

	// Right Leg
	output = Composite(output, backRightLeg, 8, 20)

	return Scale(output, opts.Scale)
}
