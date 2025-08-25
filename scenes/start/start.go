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

	// how far we can place or break blocks
	reach int
	// 3 modes. normal, slow, fast
	moveSpeedModes [3]int
	// index into the array above
	CurrentMoveSpeedMode int
	// the block currently in our hands.
	HeldBlock Blocks.Type
}

// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {

	scene.world = NewWorld(8, 8, 0)
	scene.reach = 9
	scene.HeldBlock = Blocks.Grass

	scene.cam = c.NewCamera(scene.world.Center(), 90, 10, 0.0036)
	scene.moveSpeedModes = [3]int{15, 5, 35}

	scene.skybox = LoadSkybox("assets/skybox.png")
	scene.chunkShader = rl.LoadShader("shader/chunk.vert", "shader/chunk.frag")

	start := time.Now()
	scene.world.BuildChunkMeshes()
	fmt.Println("Meshed in", time.Since(start))

	atlas := CreateAtlas()
	atlasImg := rl.NewImageFromImage(atlas)
	defer rl.UnloadImage(atlasImg)

	scene.atlas = rl.LoadTextureFromImage(atlasImg)

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

	scene.world.Render(scene.cam, scene.chunkShader, scene.atlas)
	rl.EndMode3D()
	DrawCrosshair(30)
	if scene.HeldBlock == Blocks.Air {
		panic("player should never hold air")
	}
	// draw held block
	rect := AtlasCoordinates(scene.HeldBlock)
	topRightCorner := rl.GetRenderWidth()
	rl.DrawTextureRec(scene.atlas, rect, c.V2(topRightCorner-5, 5).R(), rl.White)

	rl.GetKeyPressed()

	rl.DrawTexture(scene.atlas, 300, 300, rl.White)
	rl.DrawText(fmt.Sprintf("Speed: %.2f\n Ctrl to change", scene.cam.MoveSpeed), 5, 100, 20, rl.RayWhite)
	if rl.IsKeyPressed(rl.KeyLeftControl) {
		// 0, 1 or 2
		scene.CurrentMoveSpeedMode++
		scene.CurrentMoveSpeedMode %= len(scene.moveSpeedModes)
		scene.cam.MoveSpeed = float32(scene.moveSpeedModes[scene.CurrentMoveSpeedMode])
	}
	if ctx.DebugMenuEnabled {
		rl.DrawFPS(10, 10)
		// Draw Coordinates
		rl.DrawText(scene.cam.Position.String(), 5, 30, ctx.DebugFontSize, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("World Size: %dx%dx%d (%d chunks)", scene.world.Width, scene.world.Height, scene.world.Depth, scene.world.Volume), 5, 70, 20, rl.RayWhite)
		ctx.MemoryStatsCords.X = 5
		ctx.MemoryStatsCords.Y = 100 + 60
	}
	if rl.IsMouseButtonPressed(rl.MouseButtonLeft) {
		scene.breakBlock()
	} else if rl.IsMouseButtonPressed(rl.MouseButtonRight) {
		scene.placeBlock()
	}

	return false // if true is returned, Unload is called
}
func (scene *Scene) placeBlock() {
	hit, pos, normal := scene.world.RaycastVoxel(scene.cam.Position, scene.cam.LookVector(), float32(scene.reach))
	if hit {
		x, y, z := pos.Add(normal).ToInt()
		ok := scene.world.SetBlockID(x, y, z, scene.HeldBlock)
		if !ok {
			return
		}
		chunk, ok := scene.world.ChunkAtWorld(x, y, z)
		if ok {
			chunk.Chunk.Empty = false
			scene.world.RefreshChunkMesh(chunk)
		}
	}
}

func (scene *Scene) breakBlock() {
	hit, pos, _ := scene.world.RaycastVoxel(scene.cam.Position, scene.cam.LookVector(), float32(scene.reach))
	if hit {
		x, y, z := pos.ToInt()
		scene.world.SetBlockID(x, y, z, Blocks.Air)
		chunk, ok := scene.world.ChunkAtWorld(x, y, z)
		if ok {
			scene.world.RefreshChunkMesh(chunk)
		}
	}
}

// called after Update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	scene.skybox.Unload()
	return "someOtherSceneId" // the engine will switch to the scene that is registered with this id
}
