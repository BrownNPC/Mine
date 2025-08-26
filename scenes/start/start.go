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
	world  *World

	atlas, HudAtlas rl.Texture2D
	chunkShader     rl.Shader
	chunkMesh       ChunkMesh

	// how far we can place or break blocks
	reach int
	// 3 modes. normal, slow, fast
	moveSpeedModes [3]int
	// index into the array above
	CurrentMoveSpeedMode int
	// the block currently in our hands.
	HeldBlock        Blocks.Type
	DebugMenuEnabled bool
}

// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {

	scene.world = NewWorld(ctx.WorldGenConfig.Width, ctx.WorldGenConfig.Heght, ctx.WorldGenConfig.Seed)
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

	tx0 := rl.GetShaderLocation(scene.chunkShader, "texture0")
	gl.UseProgram(scene.chunkShader.ID)
	gl.Uniform1i(tx0, 0)
	scene.atlas = rl.LoadTextureFromImage(atlasImg)
	scene.HudAtlas = rl.LoadTextureFromImage(atlasImg)

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
	if rl.IsKeyPressed(rl.KeyF6) {
		scene.world.Unload()
		scene.world = nil
		scene.world = NewWorld(ctx.WorldGenConfig.Width, ctx.WorldGenConfig.Heght, int64(time.Now().Nanosecond()))
		scene.world.BuildChunkMeshes()
	}
	// draw held block
	scene.DrawHeldBlock()
	// rl.DrawTextureRec(scene.atlas, rect, c.V2(topRightCorner-5, 5).R(), rl.White)

	rl.GetKeyPressed()

	rl.DrawText(fmt.Sprintf("Speed: %.2f\nCtrl to change", scene.cam.MoveSpeed), 5, 140, 20, rl.RayWhite)
	if rl.IsKeyPressed(rl.KeyLeftControl) {
		// 0, 1 or 2
		scene.CurrentMoveSpeedMode++
		scene.CurrentMoveSpeedMode %= len(scene.moveSpeedModes)
		scene.cam.MoveSpeed = float32(scene.moveSpeedModes[scene.CurrentMoveSpeedMode])
	}
	// scroll block
	if wheelMove := rl.GetMouseWheelMoveV().Y; wheelMove != 0 {
		if wheelMove > 0 {
			scene.HeldBlock++
			if scene.HeldBlock == Blocks.TotalBlocks {
				scene.HeldBlock--
			}
		} else {
			// 0 is air, we dont want the player to hold air.
			scene.HeldBlock--
			if scene.HeldBlock == 0 {
				scene.HeldBlock = 1
			}
		}
	}
	if rl.IsKeyPressed(rl.KeyF3) {
		scene.DebugMenuEnabled = !scene.DebugMenuEnabled
	}

	if scene.DebugMenuEnabled {
		rl.DrawFPS(10, 10)
		// Draw Coordinates
		rl.DrawText(scene.cam.Position.String(), 5, 30, 20, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("World Size: %dx%dx%d (%d chunks)", scene.world.Width, scene.world.Height, scene.world.Depth, scene.world.Volume), 5, 70, 20, rl.RayWhite)
		rl.DrawText(fmt.Sprintf("Held Block %s", scene.HeldBlock.String()), 5, 90, 20, rl.White)

		// draw HeldBlock
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
	rl.UnloadShader(scene.chunkShader)
	return "someOtherSceneId" // the engine will switch to the scene that is registered with this id
}
