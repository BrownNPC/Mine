package c

import (
	"GameFrameworkTM/components/Blocks"
)

const CHUNK_SIZE = 32
const CHUNK_HALF_SIZE = CHUNK_SIZE / 2
const CHUNK_FACE_AREA = CHUNK_SIZE * CHUNK_SIZE
const CHUNK_VOLUME = CHUNK_SIZE * CHUNK_FACE_AREA

type Chunk struct {
	Blocks [CHUNK_VOLUME]Blocks.Type
}

// convert xyz cords into index of 1d array
func (c *Chunk) Linearize(x, y, z int) int {
	Index := x + CHUNK_SIZE*z + CHUNK_FACE_AREA*y
	return Index
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
