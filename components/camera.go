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
	Aspect     float32

	FovyRad, FovxRad float32
}

// forwards is -Z (north)
var Forwards = V3(0, 0, -1)

func NewCamera(pos Vec3, fov, moveSpeed float32, sensitivity float32) Camera {
	cam := Camera{
		MoveSpeed:   moveSpeed,
		Sensitivity: sensitivity,
		Position:    pos,
		Target:      pos.Add(Forwards),
		FOV:         fov,
		Pitch:       0,
		Yaw:         0,
	}
	cam.CaclulateFOV(fov)
	return cam
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
	c.Aspect = float32(rl.GetRenderWidth()) / float32(rl.GetRenderHeight())
	c.CaclulateFOV(c.FOV)
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
	forwards := dir.Sub(V3(0, dir.Y, 0))
	if rl.IsKeyDown(rl.KeyW) {
		move = move.Add(forwards)
	}
	if rl.IsKeyDown(rl.KeyS) {
		move = move.Sub(forwards)
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
func (c *Camera) LookVector() Vec3 {
	return V3(
		float32(math.Cos(c.Pitch)*math.Sin(c.Yaw)),
		float32(math.Sin(c.Pitch)),
		float32(math.Cos(c.Pitch)*math.Cos(c.Yaw)),
	)
}

// convert fov in degrees to radians and get vertical and horizontal fov
func (c *Camera) CaclulateFOV(degFOV float32) {
	aspect := float32(rl.GetRenderWidth()) / float32(rl.GetRenderHeight())
	fovyRad := degFOV * rl.Deg2rad
	fovxRad := 2 * math.Atan(math.Tan(float64(fovyRad/2))*float64(aspect))
	c.FovyRad = fovyRad
	c.FovxRad = float32(fovxRad)
}

// perform frustum culling math
func (c *Camera) IsInView(chunk *Chunk) bool {
	const NEAR = 0.1
	const FAR = 10000
	r := CHUNK_SPHERE_RADIUS
	halfX := float64(c.FovxRad * 0.5)
	halfY := float64(c.FovyRad * 0.5)

	factorY := 1.0 / math.Cos(halfY)
	tanY := math.Tan(halfY)

	factorX := 1.0 / math.Cos(halfX)
	tanX := math.Tan(halfX)

	// camera axes
	forward := c.LookVector().Norm()
	right := forward.Cross(V3(0, 1, 0)).Norm()
	up := right.Cross(forward).Norm()

	center := chunk.Center().Sub(c.Position)

	// outside NEAR and FAR planes?
	sz := float64(center.Dot(forward))
	if sz < NEAR-r || sz > FAR+r {
		return false
	}

	// outside TOP and BOTTOM planes?
	sy := float64(center.Dot(up))
	dist := factorY*r + sz*tanY

	if sy < -dist || sy > dist {
		return false
	}

	sx := float64(center.Dot(right))
	// outside the LEFT and RIGHT plane
	dist = factorX*r + sz*tanX
	if sx < -dist || sx > dist {
		return false
	}

	return true
}
