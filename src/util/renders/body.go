package renders

import "image"

func RenderBody(skin *image.NRGBA, opts RenderOptions) *image.NRGBA {
	scale := float64(opts.Scale)
	slimOffset := GetSlimOffset(opts.Slim)

	output := image.NewNRGBA(image.Rect(0, 0, 20*opts.Scale, 45*opts.Scale+int(scale*(1.0/16.0))))

	var (
		frontHead     *image.NRGBA = RemoveTransparency(Extract(skin, 8, 8, 8, 8))
		topHead       *image.NRGBA = RemoveTransparency(Extract(skin, 8, 0, 8, 8))
		rightHead     *image.NRGBA = RemoveTransparency(Extract(skin, 0, 8, 8, 8))
		frontTorso    *image.NRGBA = RemoveTransparency(Extract(skin, 20, 20, 8, 12))
		frontLeftArm  *image.NRGBA = nil
		topLeftArm    *image.NRGBA = nil
		frontRightArm *image.NRGBA = RemoveTransparency(Extract(skin, 44, 20, 4-slimOffset, 12))
		topRightArm   *image.NRGBA = RemoveTransparency(Extract(skin, 44, 16, 4-slimOffset, 4))
		rightRightArm *image.NRGBA = RemoveTransparency(Extract(skin, 40, 20, 4, 12))
		frontLeftLeg  *image.NRGBA = nil
		frontRightLeg *image.NRGBA = RemoveTransparency(Extract(skin, 4, 20, 4, 12))
		rightRightLeg *image.NRGBA = RemoveTransparency(Extract(skin, 0, 20, 4, 12))
	)

	if IsOldSkin(skin) {
		frontLeftArm = FlipHorizontal(frontRightArm)
		topLeftArm = FlipHorizontal(topRightArm)
		frontLeftLeg = FlipHorizontal(frontRightLeg)
	} else {
		frontLeftArm = RemoveTransparency(Extract(skin, 36, 52, 4-slimOffset, 12))
		topLeftArm = RemoveTransparency(Extract(skin, 36, 48, 4-slimOffset, 4))
		frontLeftLeg = RemoveTransparency(Extract(skin, 20, 52, 4, 12))

		if opts.Overlay {
			overlaySkin := FixTransparency(skin)

			frontHead = Composite(frontHead, Extract(overlaySkin, 40, 8, 8, 8), 0, 0)
			topHead = Composite(topHead, Extract(overlaySkin, 40, 0, 8, 8), 0, 0)
			rightHead = Composite(rightHead, Extract(overlaySkin, 32, 8, 8, 8), 0, 0)
			frontTorso = Composite(frontTorso, Extract(overlaySkin, 20, 36, 8, 12), 0, 0)
			frontLeftArm = Composite(frontLeftArm, Extract(overlaySkin, 52, 52, 4-slimOffset, 64), 0, 0)
			topLeftArm = Composite(topLeftArm, Extract(overlaySkin, 52, 48, 4-slimOffset, 4), 0, 0)
			frontRightArm = Composite(frontRightArm, Extract(overlaySkin, 44, 36, 4-slimOffset, 48), 0, 0)
			topRightArm = Composite(topRightArm, Extract(overlaySkin, 44, 48, 4-slimOffset, 4), 0, 0)
			rightRightArm = Composite(rightRightArm, Extract(overlaySkin, 40, 36, 4, 12), 0, 0)
			frontLeftLeg = Composite(frontLeftLeg, Extract(overlaySkin, 4, 52, 4, 12), 0, 0)
			frontRightLeg = Composite(frontRightLeg, Extract(overlaySkin, 4, 36, 4, 12), 0, 0)
			rightRightLeg = Composite(rightRightLeg, Extract(overlaySkin, 0, 36, 4, 12), 0, 0)
		}
	}

	// Right Side of Right Leg
	output = CompositeTransform(output, Scale(rightRightLeg, opts.Scale), TransformRight, 4*scale, 23*scale)

	// Front of Right Leg
	output = CompositeTransform(output, Scale(frontRightLeg, opts.Scale), TransformForward, 8*scale, 31*scale)

	// Front of Left Leg
	output = CompositeTransform(output, Scale(frontLeftLeg, opts.Scale), TransformForward, 12*scale, 31*scale)

	// Front of Torso
	output = CompositeTransform(output, Scale(frontTorso, opts.Scale), TransformForward, 8*scale, 19*scale)

	// Front of Right Arm
	output = CompositeTransform(output, Scale(frontRightArm, opts.Scale), TransformForward, float64(4+slimOffset)*scale, 19*scale-1)

	// Front of Left Arm
	output = CompositeTransform(output, Scale(frontLeftArm, opts.Scale), TransformForward, 16*scale, 21*scale-1)

	// Top of Left Arm
	output = CompositeTransform(output, Scale(topLeftArm, opts.Scale), TransformUp, -5*scale, 17*scale)

	// Right Side of Right Arm
	output = CompositeTransform(output, Scale(rightRightArm, opts.Scale), TransformRight, float64(slimOffset)*scale, float64(15-slimOffset)*scale)

	// Top of Right Arm
	output = CompositeTransform(output, Scale(topRightArm, opts.Scale), TransformUp, float64(-15+slimOffset)*scale, 15*scale)

	// Front of Head
	output = CompositeTransform(output, Scale(frontHead, opts.Scale), TransformForward, 10*scale, 13*scale-1)

	// Top of Head
	output = CompositeTransform(output, Scale(topHead, opts.Scale), TransformUp, -3*scale, 5*scale)

	// Right Side of Head
	output = CompositeTransform(output, Scale(rightHead, opts.Scale), TransformRight, 2*scale, 3*scale)

	return output
}
