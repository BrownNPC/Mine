package c

import (
	"GameFrameworkTM/components/Blocks"
	"math"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const CHUNK_SIZE = 32
const CHUNK_SIZE_HALF = CHUNK_SIZE / 2
const CHUNK_FACE_AREA = CHUNK_SIZE * CHUNK_SIZE
const CHUNK_VOLUME = CHUNK_SIZE * CHUNK_FACE_AREA

var CHUNK_SPHERE_RADIUS = CHUNK_SIZE_HALF * math.Sqrt(3)

type Chunk struct {
	X, Y, Z int
	Empty   bool
	Blocks  [CHUNK_VOLUME]Blocks.Type
}

func NewChunk(X, Y, Z int) Chunk {
	return Chunk{
		X: X,
		Y: Y,
		Z: Z,
	}
}

// convert xyz cords into index of 1d array
func (c *Chunk) Linearize(x, y, z int) int {
	index := x + (CHUNK_SIZE * z) + (CHUNK_FACE_AREA * y)
	return index
}
func (c *Chunk) Get(x, y, z int) Blocks.Type {
	return c.Blocks[c.Linearize(x, y, z)]
}

func (c *Chunk) Set(x, y, z int, T Blocks.Type) {
	c.Blocks[c.Linearize(x, y, z)] = T
}
func (c *Chunk) IsAir(x, y, z int) bool {
	if 0 <= x && x < CHUNK_SIZE {
		if 0 <= y && y < CHUNK_SIZE {
			if 0 <= z && z < CHUNK_SIZE {
				return c.Get(x, y, z) == Blocks.Air
			}
		}
	}
	return false
}

func (c *Chunk) GetModelMatrix() rl.Matrix {
	position := V3(c.X, c.Y, c.Z).Scale(CHUNK_SIZE)
	x, y, z := position.XYZ()
	transform := rl.MatrixTranslate(x, y, z)
	return transform
}

// center in world coordinates
func (c *Chunk) Center() Vec3 {
	center := V3(
		float32(c.X)+0.5,
		float32(c.Y)+0.5,
		float32(c.Z)+0.5,
	).Scale(CHUNK_SIZE)
	return center
}
