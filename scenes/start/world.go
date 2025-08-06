package start

import (
	c "GameFrameworkTM/components"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ojrac/opensimplex-go"
)

type World struct {
	Width, Height       int
	Depth, Volume, Area int
	CenterXZ, CenterY   int
	Chunks              c.ThreeDimensionalArray[ChunkMesh]
	NoiseGenerator      opensimplex.Noise32
}

func NewWorld(width, height int, seed int64) World {
	world := World{
		Width:  width,
		Height: height,
		Depth:  width,
	}
	world.Area = world.Width * world.Depth
	world.Volume = world.Area * world.Height
	world.CenterXZ = world.Width * c.CHUNK_SIZE_HALF
	world.CenterY = world.Height * c.CHUNK_SIZE_HALF
	// volume
	world.Chunks = c.New3dArray[ChunkMesh](world.Width, world.Height, world.Depth)

	world.NoiseGenerator = opensimplex.New32(seed)

	// arrange chunks
	start := time.Now()
	for x := range world.Width {
		for y := range world.Height {
			for z := range world.Depth {
				mesh := NewChunkMesh(x, y, z)
				InitChunk(&mesh.Chunk, world.NoiseGenerator)
				world.Chunks.Set(x, y, z, mesh)
			}
		}
	}
	fmt.Println("Arranged chunks in", time.Since(start))
	return world
}

func (world World) Center() c.Vec3 {
	return c.V3(world.CenterXZ, world.Height*c.CHUNK_SIZE, world.CenterXZ)
}

// BuildChunkMeshes initiates building of all chunk meshes
// and uploads them to the GPU
func (world *World) BuildChunkMeshes() {
	dirs := [6][3]int{
		{-1, 0, 0}, // -X
		{+1, 0, 0}, // +X
		{0, -1, 0}, // -Y
		{0, +1, 0}, // +Y
		{0, 0, -1}, // -Z
		{0, 0, +1}, // +Z
	}
	for x := range world.Width {
		for y := range world.Height {
			for z := range world.Depth {
				chunk := world.Chunks.Get(x, y, z) // value copy

				// Gather pointers to neighboring chunks
				var neighbors [6]*c.Chunk

				for i, d := range dirs {
					nx, ny, nz := x+d[0], y+d[1], z+d[2]
					if nx >= 0 && nx < world.Width &&
						ny >= 0 && ny < world.Height &&
						nz >= 0 && nz < world.Depth {
						// Get neighbor by reference (via index in backing slice)
						neighbors[i] = &world.Chunks.GetRef(nx, ny, nz).Chunk
					}
				}

				// Setup this chunk with its neighbors
				vertices := chunk.BuildVerticies(neighbors)
				chunk.Setup(vertices)

				// Store back updated chunk
				world.Chunks.Set(x, y, z, chunk)
			}
		}
	}
}

func (world *World) Render(cam c.Camera, shader rl.Shader, textures rl.Texture2D) {
	chunks := world.Chunks.BackingArray()
	for i := range chunks {
		chunks[i].Render(cam, shader, textures)
	}
}
