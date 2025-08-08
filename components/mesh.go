package c

import (
	"fmt"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type BaseMesh struct {
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

func SetupMesh(m *BaseMesh, vertexData unsafe.Pointer, totalBytes int, format []VertexAttrib) {
	gl.GenVertexArrays(1, &m.VAO)
	gl.BindVertexArray(m.VAO)

	gl.GenBuffers(1, &m.VBO)
	gl.BindBuffer(gl.ARRAY_BUFFER, m.VBO)

	// bytes per vertex
	stride := 0
	for _, attr := range format {
		stride += attr.Count * SizeOfGLType(attr.Type)
	}
	// total vertices = totalBytes / bytesPerVertex
	m.VertexCount = int32(totalBytes / stride)
	gl.BufferData(gl.ARRAY_BUFFER, totalBytes, vertexData, gl.STATIC_DRAW)

	// Setup attributes
	offset := uintptr(0)
	for _, attr := range format {
		gl.EnableVertexAttribArray(attr.Location)

		bytesPerAttr := uintptr(attr.Count * SizeOfGLType(attr.Type))

		// If it's an integer type (and you declared uvec*/ivec* in GLSL)
		if attr.Normalize == false && isIntegerGLType(attr.Type) {
			gl.VertexAttribIPointer(
				attr.Location,
				int32(attr.Count),
				attr.Type,
				int32(stride),
				unsafe.Pointer(offset),
			)
		} else {
			// float attributes or normalized ints
			gl.VertexAttribPointer(
				attr.Location,
				int32(attr.Count),
				attr.Type,
				attr.Normalize,
				int32(stride),
				unsafe.Pointer(offset),
			)
		}

		offset += bytesPerAttr
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	gl.BindVertexArray(0)

}

// Render the mesh
// texture is needed for texture id
func (m *BaseMesh) Render(cam Camera, shader rl.Shader, texture rl.Texture2D, model rl.Matrix) {
	rl.DrawRenderBatchActive()

	gl.UseProgram(shader.ID)

	locModel := shader.GetLocation(rl.ShaderLocMatrixModel)
	locView := shader.GetLocation(rl.ShaderLocMatrixView)
	locProjection := shader.GetLocation(rl.ShaderLocMatrixProjection)

	view := rl.GetCameraMatrix(cam.R())
	projection := rl.GetMatrixProjection()

	rl.SetShaderValueMatrix(shader, locModel, model)
	rl.SetShaderValueMatrix(shader, locView, view)
	rl.SetShaderValueMatrix(shader, locProjection, projection)

	gl.ActiveTexture(gl.TEXTURE0)
	gl.BindTexture(gl.TEXTURE_2D, texture.ID)

	gl.BindVertexArray(m.VAO)
	gl.DrawArrays(gl.TRIANGLES, 0, m.VertexCount)
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
func isIntegerGLType(t uint32) bool {
	switch t {
	case gl.UNSIGNED_BYTE, gl.INT, gl.UNSIGNED_INT:
		return true
	default:
		return false
	}

}
func TotalBytes[T any](slice []T) int {
	var zero T
	return len(slice) * int(unsafe.Sizeof(zero))
}

// Unload deletes the VAO and VBO from GPU memory.
func (m *BaseMesh) Unload() {
	if m.VAO != 0 {
		gl.DeleteVertexArrays(1, &m.VAO)
		m.VAO = 0
	}
	if m.VBO != 0 {
		gl.DeleteBuffers(1, &m.VBO)
		m.VBO = 0
	}
	m.VertexCount = 0
}
