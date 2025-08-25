package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"

	"github.com/ojrac/opensimplex-go"
)

func InitChunk(chunk *c.Chunk, noiseGen opensimplex.Noise32) {
	const (
		scale = 0.01
	)

	// compute world‐space corner of this chunk
	cx := chunk.X * c.CHUNK_SIZE
	cy := chunk.Y * c.CHUNK_SIZE
	cz := chunk.Z * c.CHUNK_SIZE

	for x := range c.CHUNK_SIZE {
		for z := range c.CHUNK_SIZE {
			wx := x + cx
			wz := z + cz

			n := noiseGen.Eval2(float32(wx)*scale, float32(wz)*scale)
			worldHeight := int(n*c.CHUNK_SIZE + c.CHUNK_SIZE)

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
