package start

import (
	c "GameFrameworkTM/components"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func DrawCrosshair(size int32) {
	cx := int32(rl.GetScreenWidth() / 2)  // center x
	cy := int32(rl.GetScreenHeight() / 2) // center y
	// Get a scale factor based on the screen size
	// Choose smallest dimension to keep consistent scaling
	minDim := min(rl.GetScreenWidth(), rl.GetScreenHeight())

	// Scale factor: 1080p window = 1.0
	scale := float32(minDim) / 1080.0

	// Final crosshair size
	crosshairSize := int32(float32(size) * scale)
	lineThickness := 3 * scale

	horizontalLineStart := c.V2(cx-(crosshairSize/2), cy)
	horizontalLineEnd := c.V2(cx+(crosshairSize/2), cy)
	rl.DrawLineEx(
		horizontalLineStart.R(),
		horizontalLineEnd.R(),
		lineThickness,
		rl.LightGray,
	)
	verticalLineStart := c.V2(
		cx, cy+(crosshairSize/2),
	)
	verticalLineEnd := c.V2(
		cx, cy-(crosshairSize/2),
	)
	rl.DrawLineEx(
		verticalLineStart.R(),
		verticalLineEnd.R(),
		lineThickness,
		rl.LightGray,
	)

}
