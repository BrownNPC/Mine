package start

import (
	c "GameFrameworkTM/components"
	"unsafe"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Skybox struct {
	model  rl.Model
	Unload func()
}

func LoadSkybox(location string) Skybox {

	// load skybox shader and set required locations
	skyboxShader := rl.LoadShader("shader/skybox.vs", "shader/skybox.fs")

	setShaderIntValue(skyboxShader, "environmentMap", rl.MapCubemap)

	cube := rl.GenMeshCube(1, 1, 1)
	skybox := rl.LoadModelFromMesh(cube)
	skybox.Materials.Shader = skyboxShader

	// Load skybox image
	skyboxImg := rl.LoadImage(location)
	defer rl.UnloadImage(skyboxImg) // Free CPU memory after cubemap is created

	skyboxTexture := rl.LoadTextureCubemap(skyboxImg, rl.CubemapLayoutLineVertical)
	rl.SetMaterialTexture(skybox.Materials, rl.MapCubemap, skyboxTexture)

	return Skybox{model: skybox, Unload: func() {
		rl.UnloadTexture(skyboxTexture)
		rl.UnloadShader(skyboxShader)
		rl.UnloadModel(skybox)
	}}
}

func (skybox Skybox) Draw(cameraCoordinates c.Vec3) {
	rl.DisableBackfaceCulling()
	rl.DisableDepthMask()

	rl.DrawModel(skybox.model, cameraCoordinates.R(), 1, rl.White)

	rl.EnableBackfaceCulling()
	rl.EnableDepthMask()
}
func setShaderIntValue(shader rl.Shader, name string, value int32) {
	rl.SetShaderValue(
		shader,
		rl.GetShaderLocation(shader, name),
		unsafe.Slice((*float32)(unsafe.Pointer(&value)), 4),
		rl.ShaderUniformInt,
	)
}
