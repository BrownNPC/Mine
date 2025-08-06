package c

import (
	"fmt"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Vec2 rl.Vector2
type Vec3 rl.Vector3

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func V2[T Number](x, y T) Vec2 {
	return Vec2{float32(x), float32(y)}
}
func V3[T Number](x, y, z T) Vec3 {
	return Vec3{float32(x), float32(y), float32(z)}
}

// convert to raylib vector
func (v Vec2) R() rl.Vector2 {
	return rl.Vector2(v)
}

// convert to raylib vector
func (v Vec3) R() rl.Vector3 {
	return rl.Vector3(v)
}

// vector 2 zero
var V2Z = V2(0, 0)

// vector 3 zero
var V3Z = V3(0, 0, 0)

//
// === Vec2 Math ===
//

func (v Vec2) Add(o Vec2) Vec2      { return V2(v.X+o.X, v.Y+o.Y) }
func (v Vec2) Sub(o Vec2) Vec2      { return V2(v.X-o.X, v.Y-o.Y) }
func (v Vec2) Mul(o Vec2) Vec2      { return V2(v.X*o.X, v.Y*o.Y) }
func (v Vec2) Scale(s float32) Vec2 { return V2(v.X*s, v.Y*s) }
func (v Vec2) Dot(o Vec2) float32   { return v.X*o.X + v.Y*o.Y }
func (v Vec2) Len() float32         { return float32(math.Hypot(float64(v.X), float64(v.Y))) }
func (v Vec2) Norm() Vec2 {
	l := v.Len()
	if l == 0 {
		return V2Z
	}
	return v.Scale(1 / l)
}
func (v Vec2) Dist(o Vec2) float32 {
	return v.Sub(o).Len()
}
func (v Vec2) XY() (float32, float32) {
	return v.X, v.Y
}

//
// === Vec3 Math ===
//

func (v Vec3) Add(o Vec3) Vec3      { return V3(v.X+o.X, v.Y+o.Y, v.Z+o.Z) }
func (v Vec3) Sub(o Vec3) Vec3      { return V3(v.X-o.X, v.Y-o.Y, v.Z-o.Z) }
func (v Vec3) Mul(o Vec3) Vec3      { return V3(v.X*o.X, v.Y*o.Y, v.Z*o.Z) }
func (v Vec3) Scale(s float32) Vec3 { return V3(v.X*s, v.Y*s, v.Z*s) }
func (v Vec3) Dot(o Vec3) float32   { return v.X*o.X + v.Y*o.Y + v.Z*o.Z }
func (v Vec3) Len() float32 {
	return float32(math.Sqrt(float64(v.X*v.X + v.Y*v.Y + v.Z*v.Z)))
}
func (v Vec3) Norm() Vec3 {
	l := v.Len()
	if l == 0 {
		return V3Z
	}
	return v.Scale(1 / l)
}
func (v Vec3) Dist(o Vec3) float32 {
	return v.Sub(o).Len()
}
func (v Vec3) Cross(o Vec3) Vec3 {
	return V3(
		v.Y*o.Z-v.Z*o.Y,
		v.Z*o.X-v.X*o.Z,
		v.X*o.Y-v.Y*o.X,
	)
}
func (v Vec3) Floor() Vec3 {
	return V3(float32(math.Floor(float64(v.X))), float32(math.Floor(float64(v.Y))), float32(math.Floor(float64(v.Z))))
}
func (v Vec3) Round() Vec3 {
	return V3(float32(math.Round(float64(v.X))), float32(math.Round(float64(v.Y))), float32(math.Round(float64(v.Z))))
}
func (v Vec3) ToInt() (int, int, int) {
	return int(v.X), int(v.Y), int(v.Z)
}
func (v Vec3) XYZ() (float32, float32, float32) {
	return v.X, v.Y, v.Z
}

func (v Vec3) String() string {
	return fmt.Sprintf("XYZ: %.0f %.0f %.0f", v.X, v.Y, v.Z)
}
func (v Vec2) String() string {
	return fmt.Sprintf("XY :%.2f Y:%.2f", v.X, v.Y)
}
