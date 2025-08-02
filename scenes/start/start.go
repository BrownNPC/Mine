package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Scene struct {
	cam      c.Camera
	skybox   Skybox
	quadMesh QuadMesh
}

// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {
	scene.cam = c.NewCamera(c.V3(0, 0, 0), 90, 1.8, 0.0036)
	scene.skybox = LoadSkybox("assets/skybox.png")
	// chunkShader := rl.LoadShader("shader/chunk.vert", "shader/chunk.vert")
	quadShader := rl.LoadShader("shader/quad.vert", "shader/quad.frag")
	scene.quadMesh = QuadMesh{}
	scene.quadMesh.Setup(quadShader)
}

// update is called every frame
func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	scene.cam.Update()
	rl.BeginMode3D(scene.cam.R())
	scene.skybox.Draw(scene.cam.Position)

	rl.DisableBackfaceCulling()
	scene.quadMesh.Render(gl.TRIANGLES)
	rl.EnableBackfaceCulling()
	rl.EndMode3D()
	DrawCrosshair(30)
	// Draw Coordinates
	rl.DrawFPS(10, 5)
	rl.DrawText(scene.cam.Position.String(), 10, 30, 20, rl.RayWhite)
	return false // if true is returned, Unload is called
}

// called after Update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	scene.skybox.Unload()
	return "someOtherSceneId" // the engine will switch to the scene that is registered with this id
}
