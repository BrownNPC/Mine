package c

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type BaseMesh struct {
	Shader      rl.Shader
	VAO         uint32
	VBO         uint32
	VertexCount int32
}
type VertexAttrib struct {
	Location  uint32 // attribute location
	Count     int    // number of components (e.g. 3 for vec3)
	Type      uint32 // gl.FLOAT, gl.UNSIGNED_BYTE, etc.
	Normalize bool   // normalize the data?
}

// Setup initializes the VBO, VAO, and binds attributes

func SetupMesh(m *BaseMesh, data unsafe.Pointer, totalBytes int, format []VertexAttrib) {
	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)

	var boundVAO int32
	gl.GetIntegerv(gl.VERTEX_ARRAY_BINDING, &boundVAO)
	fmt.Println("VAO bound:", boundVAO, "expected:", m.VAO)

	var boundVBO int32
	gl.GetIntegerv(gl.ARRAY_BUFFER_BINDING, &boundVBO)
	fmt.Println("VBO bound:", boundVBO, "expected:", m.VBO)
	// bytes per vertex
	stride := 0
	for _, attr := range format {
		stride += attr.Count * SizeOfGLType(attr.Type)
	}
	// total vertices = totalBytes / bytesPerVertex
	m.VertexCount = int32(totalBytes / stride)
	gl.BufferData(gl.ARRAY_BUFFER, totalBytes, data, gl.STATIC_DRAW)

	// Setup attributes
	offset := uintptr(0)
	for _, attr := range format {
		gl.EnableVertexAttribArray(attr.Location)
		gl.VertexAttribPointer(attr.Location, int32(attr.Count), attr.Type, attr.Normalize, int32(stride), unsafe.Pointer(offset))
		offset += uintptr(attr.Count * SizeOfGLType(attr.Type))
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

}

// Render the mesh
func (m *BaseMesh) Render(mode uint32) {
	rl.DrawRenderBatchActive()

	gl.UseProgram(m.Shader.ID)
	mvp := rl.MatrixMultiply(rl.GetMatrixModelview(), rl.GetMatrixProjection())
	gl.UniformMatrix4fv(m.Shader.GetLocation(rl.ShaderLocMatrixMvp), 1, false, unsafe.SliceData(rl.MatrixToFloat(mvp)))
	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(mode, 0, m.VertexCount)
	if err := gl.GetError(); err != gl.NO_ERROR {
		fmt.Printf("GL error after DrawArrays: 0x%X\n", err)
	}

	gl.BindVertexArray(0)
}
func SizeOfGLType(glType uint32) int {
	switch glType {
	case gl.FLOAT:
		return 4
	case gl.UNSIGNED_BYTE:
		return 1
	// add more if needed
	default:
		panic(fmt.Sprintf("unsupported GL type: 0x%x", glType))
	}
}
func TotalBytes[T any](slice []T) int {
	var zero T
	return len(slice) * int(unsafe.Sizeof(zero))
}
