package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
	"GameFrameworkTM/engine"
	"fmt"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/go-gl/gl/v3.3-core/gl"
)

type Scene struct {
	cam    c.Camera
	skybox Skybox
	world  World

	atlas       rl.Texture2D
	chunkShader rl.Shader
	chunkMesh   ChunkMesh
}

// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {

	scene.world = NewWorld(8, 2, 0)

	scene.cam = c.NewCamera(scene.world.Center(), 90, 10, 0.0036)
	scene.skybox = LoadSkybox("assets/skybox.png")
	scene.chunkShader = rl.LoadShader("shader/chunk.vert", "shader/chunk.frag")

	start := time.Now()
	scene.world.BuildChunkMeshes()
	fmt.Println("Meshed in", time.Since(start))

	scene.atlas = rl.LoadTexture("assets/blocks/textures/Dirt.png")

	tx0 := gl.GetUniformLocation(scene.chunkShader.ID, gl.Str("texture0\x00"))
	gl.UseProgram(scene.chunkShader.ID)
	gl.Uniform1i(tx0, 0)

	rl.SetTargetFPS(240)
}

// update is called every frame
func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	scene.cam.Update()
	rl.BeginMode3D(scene.cam.R())
	scene.skybox.Draw(scene.cam.Position)

	// rl.DisableBackfaceCulling()
	// TODO: render world
	scene.world.Render(scene.cam, scene.chunkShader, scene.atlas)
	// rl.EnableBackfaceCulling()
	rl.EndMode3D()
	DrawCrosshair(30)
	rl.DrawText(fmt.Sprintf("Speed: %.2f\nScroll to change", scene.cam.MoveSpeed), 5, 100, 20, rl.RayWhite)
	if wheelMove := rl.GetMouseWheelMoveV().Y; wheelMove != 0 {
		if wheelMove > 0 {
			scene.cam.MoveSpeed++
		} else {
			scene.cam.MoveSpeed--
		}
	}
	if ctx.DebugMenuEnabled {
		rl.DrawFPS(10, 10)
		// Draw Coordinates
		rl.DrawText(scene.cam.Position.String(), 5, 30, ctx.DebugFontSize, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("World Size: %dx%dx%d (%d chunks)", scene.world.Width, scene.world.Height, scene.world.Depth, scene.world.Volume), 5, 70, 20, rl.RayWhite)
		ctx.MemoryStatsCords.X = 5
		ctx.MemoryStatsCords.Y = 100 + 60
	}
	if rl.IsMouseButtonDown(rl.MouseButtonLeft) {
		rl.DrawText(fmt.Sprint(scene.cam.LookVector()), 500, 0, 20, rl.Red)
		hit, pos := scene.world.RaycastVoxel(scene.cam.Position, scene.cam.LookVector(), 6)
		if hit {
			rl.DrawRectangle(0, 0, 20, 20, rl.Lime)
			rl.DrawText(fmt.Sprint(pos), 0, 400, 20, rl.Blue)
			x, y, z := pos.ToInt()
			scene.world.SetBlockID(x, y, z, Blocks.Dirt)
			chunk, _, _, _, ok := scene.world.ChunkAtWorld(x, y, z)
			if ok {
				chunk.Unload()
				vertices := chunk.BuildVerticies(&scene.world)
				chunk.Setup(vertices)
			}
		}
	}
	return false // if true is returned, Unload is called
}

// called after Update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	scene.skybox.Unload()
	return "someOtherSceneId" // the engine will switch to the scene that is registered with this id
}
