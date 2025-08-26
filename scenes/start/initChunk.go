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

			// for this strip of voxels
			worldHeight := world.getHeightMap(wx, wz)

			localHeight := min(worldHeight-cy, c.CHUNK_SIZE)

			for y := range localHeight {
				// worldY
				// position of this voxel within the world
				wy := y + cy
				chunk.Set(x, y, z, world.GetBlockForYlevel(wx, wy, wz, worldHeight))
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

func (world *World) getHeightMap(x, z int) int {
	noiseGen := world.NoiseGenerator
	// amplitude
	a1 := float32(world.CenterY) * 0.7
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
func (world *World) GetBlockForYlevel(blockX, blockY, blockZ int, WorldHeight int) Blocks.Type {
	// chunk levels for each block. 128 is max height
	const (
		snow  = 90
		dirt  = 85
		stone = 78
		grass = 8
	)

	// bottom chunks are stone
	if blockY < world.Height-1 {
		// create caves
		const caveScale3D = 0.09
		const caveFloorScale2D = 0.1
		// where caves can exist
		cave3d := world.NoiseGenerator.Eval3(float32(blockX)*caveScale3D, float32(blockY)*caveScale3D, float32(blockZ)*caveScale3D)
		// cave floor limit
		floorNoise := world.NoiseGenerator.Eval2(float32(blockX), float32(blockZ))
		minCaveY := int(floorNoise*3 + 3) // -1 to 0  -> 0-6
		if cave3d > 0 && blockY > minCaveY && blockY < WorldHeight-10 {
			return Blocks.Air
		}
		return Blocks.Stone
	} else {
		rng := world.RNG.Intn(7)
		ry := blockY - rng
		if snow <= ry && ry < WorldHeight {
			return Blocks.Snow
		} else if stone <= ry && ry < snow {
			return Blocks.Stone
		} else if dirt <= ry && ry < stone {
			return Blocks.Dirt
		} else if grass <= ry && ry < dirt {
			return Blocks.Grass
		} else {
			return Blocks.Sand
		}
	}
}
