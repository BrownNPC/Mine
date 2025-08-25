package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"

	"github.com/ojrac/opensimplex-go"
)

// InitChunk generates terrain
func (world *World) InitChunk(chunk *c.Chunk, noiseGen opensimplex.Noise32) {
	const (
		scale = 0.01
	)

	// compute world‐space corner of this chunk
	cx := chunk.X * c.CHUNK_SIZE
	cy := chunk.Y * c.CHUNK_SIZE
	cz := chunk.Z * c.CHUNK_SIZE

	for x := range c.CHUNK_SIZE {
		wx := x + cx
		for z := range c.CHUNK_SIZE {
			wz := z + cz

			worldHeight := world.getHeight(noiseGen, wx, wz)

			localHeight := min(worldHeight-cy, c.CHUNK_SIZE)

			for y := range localHeight {
				// wy := y + cy
				chunk.Set(x, y, z, Blocks.Dirt)
			}
		}
	}
	for _, block := range chunk.Blocks {
		if block != Blocks.Air {
			chunk.Empty = false
			return
		}
	}
	chunk.Empty = true
}
func (world *World) getHeight(noiseGen opensimplex.Noise32, x, z int) int {
	// amplitude
	a1 := float32(world.CenterY)
	// frequency
	const f1 = 0.005
	height := noiseGen.Eval2(float32(x)*f1, float32(z)*f1)*a1 + a1
	return int(height)
}
