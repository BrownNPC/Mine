package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Scene struct {
	cam         c.Camera
	skybox      Skybox
	chunkMesh   ChunkMesh
	dirtTexture rl.Texture2D
}

// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {
	scene.cam = c.NewCamera(c.V3(0, 0, 0), 90, 100, 0.0036)
	scene.skybox = LoadSkybox("assets/skybox.png")
	chunkShader := rl.LoadShader("shader/chunk.vert", "shader/chunk.frag")

	scene.dirtTexture = rl.LoadTexture("assets/blocks/textures/Dirt.png")
	tx0 := gl.GetUniformLocation(chunkShader.ID, gl.Str("texture0\x00"))
	gl.UseProgram(chunkShader.ID)
	gl.Uniform1i(tx0, 0)

	scene.chunkMesh = ChunkMesh{}
	scene.chunkMesh.Setup(chunkShader, InitChunk())
	rl.SetTargetFPS(240)
}

// update is called every frame
func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	scene.cam.Update()
	rl.BeginMode3D(scene.cam.R())
	scene.skybox.Draw(scene.cam.Position)

	rl.DisableBackfaceCulling()
	scene.chunkMesh.Render(gl.TRIANGLES, &scene.dirtTexture)
	rl.EnableBackfaceCulling()
	rl.EndMode3D()
	DrawCrosshair(30)
	if ctx.DebugMenuEnabled {
		rl.DrawFPS(10, 10)
		// Draw Coordinates
		rl.DrawText(scene.cam.Position.String(), 10, 30, ctx.DebugFontSize, rl.RayWhite)
		ctx.MemoryStatsCords.X = 10
		ctx.MemoryStatsCords.Y = 50
	}
	return false // if true is returned, Unload is called
}

// called after Update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	scene.skybox.Unload()
	return "someOtherSceneId" // the engine will switch to the scene that is registered with this id
}
