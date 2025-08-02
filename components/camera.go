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
	// Adjust yaw (horizontal angle) and pitch (vertical angle) based on mouse movement
	// Subtracting to make right/left and up/down feel natural
	c.Yaw -= float64(mouse.X * c.Sensitivity)
	c.Pitch -= float64(mouse.Y * c.Sensitivity)

	// Clamp the pitch to avoid gimbal lock (can't look straight up/down)
	const pitchLimit = rl.Pi/2 - 0.01
	if c.Pitch > pitchLimit {
		c.Pitch = pitchLimit
	}
	if c.Pitch < -pitchLimit {
		c.Pitch = -pitchLimit
	}

	// ───── Direction Vector ─────
	// Convert yaw and pitch into a normalized direction vector using spherical coordinates
	// This is the direction the camera is looking (forward)
	dir := V3(
		float32(math.Cos(c.Pitch)*math.Sin(c.Yaw)), // X axis
		float32(math.Sin(c.Pitch)),                 // Y axis (up/down)
		float32(math.Cos(c.Pitch)*math.Cos(c.Yaw)), // Z axis
	)

	// ───── Local Axes ─────
	// Calculate the right vector (camera's X axis) using cross product: right = forward × world up
	right := dir.Cross(V3(0, 1, 0)).Norm()

	// Calculate the up vector (camera's Y axis) using cross product: up = right × forward
	up := right.Cross(dir).Norm()

	// ───── Movement Input ─────
	// Move the camera along its local axes depending on which keys are pressed
	// All movement is scaled by dt
	if rl.IsKeyDown(rl.KeyW) {
		c.Position = c.Position.Add(dir.Scale(c.MoveSpeed * dt)) // Forward
	}
	if rl.IsKeyDown(rl.KeyS) {
		c.Position = c.Position.Sub(dir.Scale(c.MoveSpeed * dt)) // Backward
	}
	if rl.IsKeyDown(rl.KeyA) {
		c.Position = c.Position.Sub(right.Scale(c.MoveSpeed * dt)) // Left (strafe)
	}
	if rl.IsKeyDown(rl.KeyD) {
		c.Position = c.Position.Add(right.Scale(c.MoveSpeed * dt)) // Right (strafe)
	}
	if rl.IsKeyDown(rl.KeySpace) {
		c.Position = c.Position.Add(up.Scale(c.MoveSpeed * dt)) // Up (fly/jump)
	}
	if rl.IsKeyDown(rl.KeyLeftShift) {
		c.Position = c.Position.Sub(up.Scale(c.MoveSpeed * dt)) // Down (descend)
	}

	// where it looks. the direction vector + position
	c.Target = c.Position.Add(dir)
}
