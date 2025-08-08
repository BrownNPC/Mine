package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
)

type Vertex = [6]byte

func NewVertex(x, y, z int, BlockType Blocks.Type, FaceDirection Blocks.Direction, ambientOcclusion int) Vertex {
	return Vertex{
		uint8(x),
		uint8(y),
		uint8(z),
		uint8(BlockType),
		uint8(FaceDirection),
		uint8(ambientOcclusion),
	}
}

func (cm *ChunkMesh) BuildVerticies(world *World) []byte {
	chunk := &cm.Chunk
	// Chunk volume * max voxel vertices (18) * vertexSize (5)
	vertexSize := len(Vertex{})
	var vertexData = make([]byte, 0, c.CHUNK_VOLUME*18*vertexSize)
	for x := range c.CHUNK_SIZE {
		for y := range c.CHUNK_SIZE {
			for z := range c.CHUNK_SIZE {
				blockType := chunk.Get(x, y, z)
				if blockType == Blocks.Air {
					continue
				}
				// top face
				if world.IsAir(chunk, x, y+1, z) {
					// ambient occlusion
					ao := FaceAO(Blocks.Top, x, y, z, chunk, world)
					// top left, top right, bottom right, bottom left
					v0 := NewVertex(x, y+1, z, blockType, Blocks.Top, ao.TopLeft)
					v1 := NewVertex(x+1, y+1, z, blockType, Blocks.Top, ao.TopRight)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Top, ao.BottomRight)
					v3 := NewVertex(x, y+1, z+1, blockType, Blocks.Top, ao.BottomLeft)

					vertexData = addVerticies(vertexData, v0, v3, v2, v0, v2, v1)
				}
				// bottom face
				if world.IsAir(chunk, x, y-1, z) {
					// ambient occlusion
					ao := FaceAO(Blocks.Bottom, x, y, z, chunk, world)

					v0 := NewVertex(x, y, z, blockType, Blocks.Bottom, ao.TopLeft)
					v1 := NewVertex(x+1, y, z, blockType, Blocks.Bottom, ao.TopRight)
					v2 := NewVertex(x+1, y, z+1, blockType, Blocks.Bottom, ao.BottomLeft)
					v3 := NewVertex(x, y, z+1, blockType, Blocks.Bottom, ao.BottomRight)
					vertexData = addVerticies(vertexData, v0, v2, v3, v0, v1, v2)
				}
				// right face
				if world.IsAir(chunk, x+1, y, z) {
					// ambient occlusion
					ao := FaceAO(Blocks.Right, x, y, z, chunk, world)
					v0 := NewVertex(x+1, y, z, blockType, Blocks.Right, ao.BottomRight)
					v1 := NewVertex(x+1, y+1, z, blockType, Blocks.Right, ao.TopRight)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Right, ao.TopLeft)
					v3 := NewVertex(x+1, y, z+1, blockType, Blocks.Right, ao.BottomLeft)
					vertexData = addVerticies(vertexData, v0, v1, v2, v0, v2, v3)
				}
				// left face
				if world.IsAir(chunk, x-1, y, z) {
					// ambient occlusion
					ao := FaceAO(Blocks.Left, x, y, z, chunk, world)
					v0 := NewVertex(x, y, z, blockType, Blocks.Left, ao.BottomRight)
					v1 := NewVertex(x, y+1, z, blockType, Blocks.Left, ao.TopRight)
					v2 := NewVertex(x, y+1, z+1, blockType, Blocks.Left, ao.TopLeft)
					v3 := NewVertex(x, y, z+1, blockType, Blocks.Left, ao.BottomLeft)
					vertexData = addVerticies(vertexData, v0, v2, v1, v0, v3, v2)
				}
				// back face
				if world.IsAir(chunk, x, y, z-1) {
					// ambient occlusion
					ao := FaceAO(Blocks.Back, x, y, z, chunk, world)
					v0 := NewVertex(x, y, z, blockType, Blocks.Back, ao.BottomLeft)
					v1 := NewVertex(x, y+1, z, blockType, Blocks.Back, ao.TopLeft)
					v2 := NewVertex(x+1, y+1, z, blockType, Blocks.Back, ao.TopRight)
					v3 := NewVertex(x+1, y, z, blockType, Blocks.Back, ao.BottomRight)
					vertexData = addVerticies(vertexData, v0, v1, v2, v0, v2, v3)
				}
				// front face
				if world.IsAir(chunk, x, y, z+1) {
					// ambient occlusion
					ao := FaceAO(Blocks.Front, x, y, z, chunk, world)
					v0 := NewVertex(x, y, z+1, blockType, Blocks.Front, ao.BottomLeft)
					v1 := NewVertex(x, y+1, z+1, blockType, Blocks.Front, ao.TopLeft)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Front, ao.TopRight)
					v3 := NewVertex(x+1, y, z+1, blockType, Blocks.Front, ao.BottomRight)
					vertexData = addVerticies(vertexData, v0, v2, v1, v0, v3, v2)
				}
			}
		}
	}
	return vertexData
}
func addVerticies(vertexData []uint8, verticies ...Vertex) []uint8 {
	for i := range len(verticies) {
		vertex := verticies[i]
		vertexData = append(vertexData, vertex[:]...)
	}
	return vertexData
}

type AmbientOcclusion struct {
	TopLeft, TopRight, BottomLeft, BottomRight int
}

// FaceAO returns 4 AO values (in the same corner order) for a face
func FaceAO(dir Blocks.Direction, x, y, z int, chunk *c.Chunk, world *World) AmbientOcclusion {
	var result AmbientOcclusion
	// check 8 surrounding blocks
	//(-1, -1)   (0, -1)   (+1, -1)
	//(-1,  0)      X      (+1,  0)
	//(-1, +1)   (0, +1)   (+1, +1)
	var top, bottom, left, right, topLeft, topRight, bottomLeft, bottomRight bool
	switch dir {
	case Blocks.Top, Blocks.Bottom:
		if dir == Blocks.Top {
			y += 1
		} else {
			y -= 1
		}
		top = !world.IsAir(chunk, x+0, y, z-1)
		bottom = !world.IsAir(chunk, x+0, y, z+1)
		left = !world.IsAir(chunk, x-1, y, z+0)
		right = !world.IsAir(chunk, x+1, y, z+0)

		topLeft = !world.IsAir(chunk, x-1, y, z-1)
		topRight = !world.IsAir(chunk, x+1, y, z-1)

		bottomLeft = !world.IsAir(chunk, x-1, y, z+1)
		bottomRight = !world.IsAir(chunk, x+1, y, z+1)
	case Blocks.Left, Blocks.Right:
		if dir == Blocks.Left {
			x -= 1
		} else {
			x += 1
		}
		top = !world.IsAir(chunk, x+0, y+1, z)
		bottom = !world.IsAir(chunk, x+0, y-1, z)
		left = !world.IsAir(chunk, x, y, z+1)
		right = !world.IsAir(chunk, x, y, z-1)

		topLeft = !world.IsAir(chunk, x, y+1, z+1)
		topRight = !world.IsAir(chunk, x, y+1, z-1)

		bottomLeft = !world.IsAir(chunk, x, y-1, z+1)
		bottomRight = !world.IsAir(chunk, x, y-1, z-1)
	case Blocks.Front, Blocks.Back:
		if dir == Blocks.Front {
			z += 1
		} else {
			z -= 1
		}
		top = !world.IsAir(chunk, x+0, y+1, z)
		bottom = !world.IsAir(chunk, x+0, y-1, z)
		left = !world.IsAir(chunk, x-1, y, z)
		right = !world.IsAir(chunk, x+1, y, z)

		topLeft = !world.IsAir(chunk, x-1, y+1, z)
		topRight = !world.IsAir(chunk, x+1, y+1, z)

		bottomLeft = !world.IsAir(chunk, x-1, y-1, z)
		bottomRight = !world.IsAir(chunk, x+1, y-1, z)
	}
	result.TopLeft = calcAO(left, top, topLeft)
	result.TopRight = calcAO(top, right, topRight)
	result.BottomRight = calcAO(bottom, right, bottomRight)
	result.BottomLeft = calcAO(bottom, left, bottomLeft)

	return result
}

func calcAO(edge1, edge2, corner bool) int {
	if edge1 && edge2 {
		return 3
	}
	return (btoi(edge1) + btoi(edge2) + btoi(corner))
}
func btoi(b bool) int {
	if b {
		return 1
	}
	return 0

}
