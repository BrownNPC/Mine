package start

import (
	c "GameFrameworkTM/components"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

// a ChunkMesh wraps a chunk and a mesh
type ChunkMesh struct {
	c.BaseMesh
	Chunk c.Chunk
}

func NewChunkMesh(x, y, z int) *ChunkMesh {
	return &ChunkMesh{
		Chunk: c.NewChunk(x, y, z),
	}
}

// upload vertices to the gpu
func (m *ChunkMesh) Setup(vertices []byte) {
	if len(vertices) == 0 {
		m.Chunk.Empty = true
		return
	} else {
		m.Chunk.Empty = false
	}

	attrib := []c.VertexAttrib{
		// 3 bytes for coordinates for each block within the chunk
		{Location: 0, Count: 3, Type: gl.UNSIGNED_BYTE, Normalize: false},
		// blocktype, face direction, ambient occlusion
		{Location: 1, Count: 1, Type: gl.UNSIGNED_BYTE, Normalize: false},
		{Location: 2, Count: 1, Type: gl.UNSIGNED_BYTE, Normalize: false},
		{Location: 3, Count: 1, Type: gl.UNSIGNED_BYTE, Normalize: false},
	}
	c.SetupMesh(&m.BaseMesh, gl.Ptr(vertices), c.TotalBytes(vertices), attrib)
}
func (m *ChunkMesh) Render(cam c.Camera, shader rl.Shader, texture rl.Texture2D) {
	if m.Chunk.Empty || !cam.IsInView(&m.Chunk) {
		return
	}
	if !cam.IsInView(&m.Chunk) {
		return
	}
	model := m.Chunk.GetModelMatrix()
	m.BaseMesh.Render(cam, shader, texture, model)
}
