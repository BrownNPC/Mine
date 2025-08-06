package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
)

type Vertex = [5]byte

func NewVertex(x, y, z int, BlockType Blocks.Type, FaceDirection Blocks.Direction) Vertex {
	return Vertex{
		uint8(x),
		uint8(y),
		uint8(z),
		uint8(BlockType),
		uint8(FaceDirection),
	}
}

func (cm *ChunkMesh) BuildVerticies(Neighbours [6]*c.Chunk) []byte {
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
				if chunk.IsAirNeighbours(x, y+1, z, Neighbours) {
					v0 := NewVertex(x, y+1, z, blockType, Blocks.Top)
					v1 := NewVertex(x+1, y+1, z, blockType, Blocks.Top)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Top)
					v3 := NewVertex(x, y+1, z+1, blockType, Blocks.Top)
					vertexData = addVerticies(vertexData, v0, v3, v2, v0, v2, v1)
				}
				// bottom face
				if chunk.IsAirNeighbours(x, y-1, z, Neighbours) {
					v0 := NewVertex(x, y, z, blockType, Blocks.Bottom)
					v1 := NewVertex(x+1, y, z, blockType, Blocks.Bottom)
					v2 := NewVertex(x+1, y, z+1, blockType, Blocks.Bottom)
					v3 := NewVertex(x, y, z+1, blockType, Blocks.Bottom)
					vertexData = addVerticies(vertexData, v0, v2, v3, v0, v1, v2)
				}
				// right face
				if chunk.IsAirNeighbours(x+1, y, z, Neighbours) {
					v0 := NewVertex(x+1, y, z, blockType, Blocks.Right)
					v1 := NewVertex(x+1, y+1, z, blockType, Blocks.Right)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Right)
					v3 := NewVertex(x+1, y, z+1, blockType, Blocks.Right)
					vertexData = addVerticies(vertexData, v0, v1, v2, v0, v2, v3)
				}
				// left face
				if chunk.IsAirNeighbours(x-1, y, z, Neighbours) {
					v0 := NewVertex(x, y, z, blockType, Blocks.Left)
					v1 := NewVertex(x, y+1, z, blockType, Blocks.Left)
					v2 := NewVertex(x, y+1, z+1, blockType, Blocks.Left)
					v3 := NewVertex(x, y, z+1, blockType, Blocks.Left)
					vertexData = addVerticies(vertexData, v0, v2, v1, v0, v3, v2)
				}
				// back face
				if chunk.IsAirNeighbours(x, y, z-1, Neighbours) {
					v0 := NewVertex(x, y, z, blockType, Blocks.Back)
					v1 := NewVertex(x, y+1, z, blockType, Blocks.Back)
					v2 := NewVertex(x+1, y+1, z, blockType, Blocks.Back)
					v3 := NewVertex(x+1, y, z, blockType, Blocks.Back)
					vertexData = addVerticies(vertexData, v0, v1, v2, v0, v2, v3)
				}
				// front face
				if chunk.IsAirNeighbours(x, y, z+1, Neighbours) {
					v0 := NewVertex(x, y, z+1, blockType, Blocks.Front)
					v1 := NewVertex(x, y+1, z+1, blockType, Blocks.Front)
					v2 := NewVertex(x+1, y+1, z+1, blockType, Blocks.Front)
					v3 := NewVertex(x+1, y, z+1, blockType, Blocks.Front)
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
