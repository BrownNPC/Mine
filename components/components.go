package components

import (
	"image/color"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Position rl.Vector2
type Velocity rl.Vector2

// tags
type Ball struct{}
type Plate struct{}

type Circle float32

// width height
type Size2D struct {
	Width  int32
	Height int32
}

type Color = color.RGBA
