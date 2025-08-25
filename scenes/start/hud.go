package start

import (
	c "GameFrameworkTM/components"

	rl "github.com/gen2brain/raylib-go/raylib"
)

func (scene *Scene) DrawHeldBlock() {
	renderW := float32(rl.GetRenderWidth())
	renderH := float32(rl.GetRenderHeight())
	currentAspect := renderW / renderH
	targetAspect := float32(1920) / float32(1080)
	var scale float32
	if currentAspect > targetAspect {
		scale = renderH / float32(1080)
	} else {
		scale = renderW / float32(1920)
	}
	rect := AtlasCoordinates(scene.HeldBlock)
	topRightCorner := rl.GetRenderWidth()
	width := float32(120) * scale
	height := float32(120) * scale

	rl.DrawTexturePro(scene.HudAtlas, rect,
		rl.NewRectangle(float32(topRightCorner-int(width)-5), 5, width, height),
		c.V2Z.R(), 0, rl.White)
}
