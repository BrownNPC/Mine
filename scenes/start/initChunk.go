package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
)

// InitChunk generates terrain
func (world *World) InitChunk(chunk *c.Chunk) {
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

			worldHeight := world.getHeight(wx, wz)

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
func (world *World) getHeight(x, z int) int {
	noiseGen := world.NoiseGenerator
	// amplitude
	a1 := float32(world.CenterY)
	a2, a4, a8 := a1*0.5, a1*0.25, a1*0.125
	// frequency
	const f1 = 0.005
	const f2, f4, f8 = f1 * 2, f1 * 4, f1 * 8
	height := noiseGen.Eval2(float32(x)*f1, float32(z)*f1)*a1 + a1
	height += noiseGen.Eval2(float32(x)*f2, float32(z)*f2)*a2 + a2
	height += noiseGen.Eval2(float32(x)*f4, float32(z)*f4)*a4 + a4
	height += noiseGen.Eval2(float32(x)*f8, float32(z)*f8)*a8 + a8

	height = max(height, 1)
	return int(height)
}
