package start

import (
	c "GameFrameworkTM/components"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type QuadMesh struct {
	c.BaseMesh
}

func (m *QuadMesh) Setup(shader rl.Shader) {
	vertexData := m.GetVertexData()
	m.Shader = shader
	format := []c.VertexAttrib{
		{Location: 0, Count: 3, Type: gl.FLOAT, Normalize: false},
		{Location: 1, Count: 3, Type: gl.FLOAT, Normalize: false},
	}
	c.SetupMesh(&m.BaseMesh, gl.Ptr(vertexData), c.TotalBytes(vertexData), format)
}

func (mesh *QuadMesh) GetVertexData() []float32 {
	positions := [][3]float32{
		{0.5, 0.5, 0}, {-0.5, 0.5, 0}, {-0.5, -0.5, 0},
		{0.5, 0.5, 0}, {-0.5, -0.5, 0}, {0.5, -0.5, 0},
	}
	colors := [][3]float32{
		{0, 1, 0}, {1, 0, 0}, {1, 1, 0},
		{0, 1, 0}, {1, 1, 0}, {0, 0, 1},
	}

	result := make([]float32, 0, len(positions)*6)
	for i := 0; i < len(positions); i++ {
		pos := positions[i]
		col := colors[i]
		result = append(result, pos[0], pos[1], pos[2], col[0], col[1], col[2])
	}
	return result
}
