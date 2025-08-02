package start

import (
	c "GameFrameworkTM/components"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type ChunkMesh struct {
	c.BaseMesh
	chunk c.Chunk
}

func NewChunkMesh(shader rl.Shader, chunk c.Chunk) ChunkMesh {
	mesh := ChunkMesh{
		BaseMesh: c.BaseMesh{},
		chunk:    chunk,
	}
	return mesh
}
func (m *ChunkMesh) Setup() {
	verticies := buildTriangles(&m.chunk)
	attrib := []c.VertexAttrib{
		// 3 bytes for coordinates for each block within the chunk
		{Location: 0, Count: 3, Type: gl.UNSIGNED_BYTE, Normalize: false},
		// 2 byes for blockType and blockFaceDirection
		{Location: 1, Count: 1, Type: gl.UNSIGNED_BYTE, Normalize: false},
		{Location: 2, Count: 1, Type: gl.UNSIGNED_BYTE, Normalize: false},
	}
	c.SetupMesh(&m.BaseMesh, gl.Ptr(verticies), c.TotalBytes(verticies), attrib)
}
