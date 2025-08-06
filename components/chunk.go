package c

import (
	"GameFrameworkTM/components/Blocks"

	rl "github.com/gen2brain/raylib-go/raylib"
)

const CHUNK_SIZE = 32
const CHUNK_SIZE_HALF = CHUNK_SIZE / 2
const CHUNK_FACE_AREA = CHUNK_SIZE * CHUNK_SIZE
const CHUNK_VOLUME = CHUNK_SIZE * CHUNK_FACE_AREA

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
	// must be inside boundaries
	if x >= 0 && x < CHUNK_SIZE {
		if y >= 0 && y < CHUNK_SIZE {
			if z >= 0 && z < CHUNK_SIZE {
				if c.Get(x, y, z) != Blocks.Air {
					return false
				}
			}
		}
	}
	return true
}

func (c *Chunk) GetModelMatrix() rl.Matrix {
	position := V3(c.X, c.Y, c.Z).Scale(CHUNK_SIZE)
	x, y, z := position.XYZ()
	transform := rl.MatrixTranslate(x, y, z)
	return transform
}

// IsAirNeighbours returns true if the block at (x,y,z) is air,
// consulting neighbor chunks when (x,y,z) lies outside of the chunk.
// The neighbors slice must be length 6, ordered:
//
//	[0] = x-1, [1] = x+1
//	[2] = y-1, [3] = y+1
//	[4] = z-1, [5] = z+1
func (c *Chunk) IsAirNeighbours(x, y, z int, neighbors [6]*Chunk) bool {
	// In‐bounds? delegate to IsAir
	if x >= 0 && x < CHUNK_SIZE &&
		y >= 0 && y < CHUNK_SIZE &&
		z >= 0 && z < CHUNK_SIZE {
		return c.IsAir(x, y, z)
	}

	// Out‐of‐bounds on X?
	switch {
	case x < 0:
		if neighbors[0] != nil {
			return neighbors[0].IsAir(CHUNK_SIZE-1, y, z)
		}
		return true
	case x >= CHUNK_SIZE:
		if neighbors[1] != nil {
			return neighbors[1].IsAir(0, y, z)
		}
		return true
	}

	// Out‐of‐bounds on Y?
	switch {
	case y < 0:
		if neighbors[2] != nil {
			return neighbors[2].IsAir(x, CHUNK_SIZE-1, z)
		}
		return true
	case y >= CHUNK_SIZE:
		if neighbors[3] != nil {
			return neighbors[3].IsAir(x, 0, z)
		}
		return true
	}

	// Out‐of‐bounds on Z?
	switch {
	case z < 0:
		if neighbors[4] != nil {
			return neighbors[4].IsAir(x, y, CHUNK_SIZE-1)
		}
		return true
	case z >= CHUNK_SIZE:
		if neighbors[5] != nil {
			return neighbors[5].IsAir(x, y, 0)
		}
		return true
	}

	// If we somehow fell through (e.g. two axes OOB), treat as air
	return true
}
