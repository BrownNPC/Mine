package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
	"fmt"
	"math"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/ojrac/opensimplex-go"
)

type World struct {
	Width, Height       int
	Depth, Volume, Area int
	CenterXZ, CenterY   int
	Chunks              c.ThreeDimensionalArray[*ChunkMesh]
	NoiseGenerator      opensimplex.Noise32
}

// NewWorld generates terrain using the provided seed.
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
	world.Chunks = c.New3dArray[*ChunkMesh](world.Width, world.Height, world.Depth)

	world.NoiseGenerator = opensimplex.New32(seed)

	// arrange chunks
	start := time.Now()
	for x := range world.Width {
		for y := range world.Height {
			for z := range world.Depth {
				mesh := NewChunkMesh(x, y, z)
				world.InitChunk(&mesh.Chunk)
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
	for x := range world.Width {
		for y := range world.Height {
			for z := range world.Depth {

				chunk := world.Chunks.Get(x, y, z)
				// Gather pointers to neighboring chunks

				// Setup this chunk with its neighbors
				vertices := chunk.BuildVerticies(world)
				chunk.Setup(vertices)

			}
		}
	}
}

func (world *World) Render(cam c.Camera, shader rl.Shader, textures rl.Texture2D) {
	chunks := world.Chunks.BackingArray()
	for _, chunk := range chunks {
		chunk.Render(cam, shader, textures)
	}
}

// you may need to adapt c.Vec3 and Blocks.Air names to your codebase
// RaycastVoxelNormal returns (hit, voxelPos, faceNormal).
// voxelPos is the integer voxel coordinates of the solid block hit.
// faceNormal is a vector with one of: (-1,0,0),(1,0,0),(0,-1,0),(0,1,0),(0,0,-1),(0,0,1).
// If the ray starts inside a solid voxel, we return that voxel and faceNormal = c.V3Z.
func (world *World) RaycastVoxel(origin, dir c.Vec3, maxDist float32) (bool, c.Vec3, c.Vec3) {
	// helper: check solid block
	isSolid := func(ix, iy, iz int) bool {
		return world.GetBlockID(ix, iy, iz) != Blocks.Air
	}

	// Use float64 for the math
	ox := float64(origin.X)
	oy := float64(origin.Y)
	oz := float64(origin.Z)
	dx := float64(dir.X)
	dy := float64(dir.Y)
	dz := float64(dir.Z)

	// normalize direction so t is in world units (distance)
	dirLen := math.Sqrt(dx*dx + dy*dy + dz*dz)
	if dirLen == 0 {
		return false, c.V3Z, c.V3Z
	}
	dx /= dirLen
	dy /= dirLen
	dz /= dirLen

	maxT := float64(maxDist) // since dir is now unit-length

	// voxel coordinates containing origin: use floor
	ix := int(math.Floor(ox))
	iy := int(math.Floor(oy))
	iz := int(math.Floor(oz))

	// steps
	stepX := 1
	if dx < 0 {
		stepX = -1
	}
	stepY := 1
	if dy < 0 {
		stepY = -1
	}
	stepZ := 1
	if dz < 0 {
		stepZ = -1
	}

	// tMax: distance along ray to the first voxel boundary on each axis
	var tMaxX, tMaxY, tMaxZ float64
	fracX := ox - math.Floor(ox) // fractional part in [0,1)
	fracY := oy - math.Floor(oy)
	fracZ := oz - math.Floor(oz)

	if dx == 0 {
		tMaxX = math.Inf(1)
	} else if stepX > 0 {
		tMaxX = (1.0 - fracX) / dx
	} else {
		tMaxX = fracX / -dx
	}

	if dy == 0 {
		tMaxY = math.Inf(1)
	} else if stepY > 0 {
		tMaxY = (1.0 - fracY) / dy
	} else {
		tMaxY = fracY / -dy
	}

	if dz == 0 {
		tMaxZ = math.Inf(1)
	} else if stepZ > 0 {
		tMaxZ = (1.0 - fracZ) / dz
	} else {
		tMaxZ = fracZ / -dz
	}

	// tDelta: how far we must travel along the ray to cross one voxel on that axis
	var tDeltaX, tDeltaY, tDeltaZ float64
	if dx == 0 {
		tDeltaX = math.Inf(1)
	} else {
		tDeltaX = math.Abs(1.0 / dx)
	}
	if dy == 0 {
		tDeltaY = math.Inf(1)
	} else {
		tDeltaY = math.Abs(1.0 / dy)
	}
	if dz == 0 {
		tDeltaZ = math.Inf(1)
	} else {
		tDeltaZ = math.Abs(1.0 / dz)
	}

	// If starting inside a solid voxel, return it immediately (no reliable normal).
	if isSolid(ix, iy, iz) {
		return true, c.V3(float32(ix), float32(iy), float32(iz)), c.V3Z
	}

	// Track which axis & sign we last stepped on so we can produce a normal.
	lastAxis := -1 // 0=x, 1=y, 2=z
	lastStep := 0  // ±1

	t := 0.0
	for t <= maxT {
		// choose smallest tMax to step (tie-breaking: X, then Y, then Z — same as your original)
		if tMaxX <= tMaxY && tMaxX <= tMaxZ {
			ix += stepX
			lastAxis = 0
			lastStep = stepX
			t = tMaxX
			tMaxX += tDeltaX
		} else if tMaxY <= tMaxZ {
			iy += stepY
			lastAxis = 1
			lastStep = stepY
			t = tMaxY
			tMaxY += tDeltaY
		} else {
			iz += stepZ
			lastAxis = 2
			lastStep = stepZ
			t = tMaxZ
			tMaxZ += tDeltaZ
		}

		if t > maxT {
			break
		}
		if isSolid(ix, iy, iz) {
			// normal points from the solid block toward the empty space we came from:
			switch lastAxis {
			case 0:
				// stepped in X, so normal is -stepX on X
				return true, c.V3(float32(ix), float32(iy), float32(iz)), c.V3(float32(-lastStep), 0, 0)
			case 1:
				return true, c.V3(float32(ix), float32(iy), float32(iz)), c.V3(0, float32(-lastStep), 0)
			case 2:
				return true, c.V3(float32(ix), float32(iy), float32(iz)), c.V3(0, 0, float32(-lastStep))
			default:
				// shouldn't happen, but return without a normal
				return true, c.V3(float32(ix), float32(iy), float32(iz)), c.V3Z
			}
		}
	}

	return false, c.V3Z, c.V3Z
}

// get a blockID from any world coordinate
func (world *World) GetBlockID(x, y, z int) Blocks.Type {
	// divFloor returns (chunkIndex, localIndex) where localIndex is in [0..c.CHUNK_SIZE-1]
	cx, lx := divFloor(x, c.CHUNK_SIZE)
	cy, ly := divFloor(y, c.CHUNK_SIZE)
	cz, lz := divFloor(z, c.CHUNK_SIZE)

	ch := world.Chunks.Get(cx, cy, cz)
	if ch == nil {
		return Blocks.Air
	}

	return ch.Chunk.Get(lx, ly, lz)
}

// convert a local chunk voxel position into a world voxel position
func (world *World) LocalChunkPosToWorldPos(chunk *c.Chunk, lx, ly, lz int) (wx, wy, wz int) {
	// chunk.X/Y/Z are chunk indices (in chunk-space).
	// Multiply by chunk size to get the world coordinate of the chunk origin,
	// then add the local offset inside the chunk.
	wx = chunk.X*c.CHUNK_SIZE + lx
	wy = chunk.Y*c.CHUNK_SIZE + ly
	wz = chunk.Z*c.CHUNK_SIZE + lz
	return
}

// check if a block is air. If out of bounds, check neighbour chunks
func (world *World) IsAir(chunk *c.Chunk, lx, ly, lz int) bool {
	if lx >= 0 && lx < c.CHUNK_SIZE && ly >= 0 && ly < c.CHUNK_SIZE && lz >= 0 && lz < c.CHUNK_SIZE {
		if chunk == nil {
			return true
		}
		return chunk.IsAir(lx, ly, lz)
	}
	wx, wy, wz := world.LocalChunkPosToWorldPos(chunk, lx, ly, lz)
	return world.GetBlockID(wx, wy, wz) == Blocks.Air

}

// SetBlockID sets the block at world coordinates (x,y,z).
// Returns true if the block was set, false if the chunk was missing/empty.
func (world *World) SetBlockID(x, y, z int, id Blocks.Type) bool {

	cx, lx := divFloor(x, c.CHUNK_SIZE)
	cy, ly := divFloor(y, c.CHUNK_SIZE)
	cz, lz := divFloor(z, c.CHUNK_SIZE)

	ch := world.Chunks.Get(cx, cy, cz)
	if ch == nil {
		return false
	}

	ch.Chunk.Set(lx, ly, lz, id)

	return true
}

// ChunkAtWorld returns the chunk containing the world voxel (x,y,z),
// the local indices inside that chunk (lx,ly,lz), and ok==true if the chunk exists and is not empty.
func (world *World) ChunkAtWorld(x, y, z int) (chunk *ChunkMesh, ok bool) {

	cx, _ := divFloor(x, c.CHUNK_SIZE)
	cy, _ := divFloor(y, c.CHUNK_SIZE)
	cz, _ := divFloor(z, c.CHUNK_SIZE)

	var chunkMesh *ChunkMesh = world.Chunks.Get(cx, cy, cz)
	if chunkMesh == nil {
		return nil, false
	}

	return chunkMesh, true
}

// recalculate the mesh for this chunk and it's neighbours
func (world *World) RefreshChunkMesh(chunkMesh *ChunkMesh) {
	chunk := &chunkMesh.Chunk
	top := world.Chunks.Get(chunk.X, chunk.Y+1, chunk.Z)
	bottom := world.Chunks.Get(chunk.X, chunk.Y-1, chunk.Z)
	left := world.Chunks.Get(chunk.X-1, chunk.Y, chunk.Z)
	right := world.Chunks.Get(chunk.X+1, chunk.Y, chunk.Z)
	front := world.Chunks.Get(chunk.X, chunk.Y, chunk.Z-1)
	back := world.Chunks.Get(chunk.X, chunk.Y, chunk.Z+1)
	var neighbours = [...]*ChunkMesh{top, bottom, left, right, front, back, chunkMesh}
	for _, mesh := range neighbours {
		if mesh != nil {
			mesh.Unload()
			vertices := mesh.BuildVerticies(world)
			mesh.Setup(vertices)
		}
	}

}

// divFloor returns (chunkIndex, localIndex) where localIndex is in [0..c.CHUNK_SIZE-1]
func divFloor(n, size int) (chunkIdx, localIdx int) {
	chunkIdx = n / size
	localIdx = n % size
	if localIdx < 0 {
		chunkIdx--
		localIdx += size
	}
	return
}
