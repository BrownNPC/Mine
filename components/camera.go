package c

import (
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Camera struct {
	MoveSpeed   float32
	Sensitivity float32

	// Camera position
	Position Vec3
	// Camera target it looks-at
	Target Vec3
	FOV    float32

	Pitch, Yaw float64
}

// forwards is -Z (north)
var Forwards = V3(0, 0, -1)

func NewCamera(pos Vec3, fov, moveSpeed float32, sensitivity float32) Camera {
	return Camera{
		MoveSpeed:   moveSpeed,
		Sensitivity: sensitivity,
		Position:    pos,
		Target:      pos.Add(Forwards),
		FOV:         fov,
		Pitch:       0,
		Yaw:         0,
	}
}

// convert to raylib camera
func (c Camera) R() rl.Camera {
	return rl.Camera3D{
		Position:   c.Position.R(),
		Target:     c.Target.R(),
		Up:         V3(0, 1, 0).R(),
		Fovy:       c.FOV,
		Projection: rl.CameraPerspective,
	}
}
func (c *Camera) Update() {
	dt := rl.GetFrameTime()     // Time since last frame (in seconds)
	mouse := rl.GetMouseDelta() // Mouse movement since last frame

	// ───── Mouse Look ─────
	c.Yaw -= float64(mouse.X * c.Sensitivity)
	c.Pitch -= float64(mouse.Y * c.Sensitivity)

	// Clamp pitch
	const pitchLimit = rl.Pi/2 - 0.01
	if c.Pitch > pitchLimit {
		c.Pitch = pitchLimit
	}
	if c.Pitch < -pitchLimit {
		c.Pitch = -pitchLimit
	}

	// ───── Direction & Axes ─────
	dir := V3(
		float32(math.Cos(c.Pitch)*math.Sin(c.Yaw)),
		float32(math.Sin(c.Pitch)),
		float32(math.Cos(c.Pitch)*math.Cos(c.Yaw)),
	)
	right := dir.Cross(V3(0, 1, 0)).Norm()
	up := V3(0, 1, 0)

	// ───── Movement Input ─────
	// 1) accumulate
	move := V3(0, 0, 0)
	if rl.IsKeyDown(rl.KeyW) {
		move = move.Add(dir)
	}
	if rl.IsKeyDown(rl.KeyS) {
		move = move.Sub(dir)
	}
	if rl.IsKeyDown(rl.KeyA) {
		move = move.Sub(right)
	}
	if rl.IsKeyDown(rl.KeyD) {
		move = move.Add(right)
	}
	// vertical fly/jump
	if rl.IsKeyDown(rl.KeySpace) {
		move = move.Add(up)
	}
	if rl.IsKeyDown(rl.KeyLeftShift) {
		move = move.Sub(up)
	}

	// 2) normalize & scale
	if length := move.Len(); length > 0 {
		move = move.Scale(c.MoveSpeed * dt / length)
		c.Position = c.Position.Add(move)
	}

	// ───── Final Target ─────
	c.Target = c.Position.Add(dir)
}
