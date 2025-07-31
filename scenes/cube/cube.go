package cube

import (
	"GameFrameworkTM/engine"

	rl "github.com/gen2brain/raylib-go/raylib"
)

type Scene struct {
	cam rl.Camera3D
}

func (scene *Scene) Load(ctx engine.Context) {
	var cam rl.Camera3D
	cam.Position = rl.NewVector3(4, 4, 4)
	cam.Up = rl.Vector3{Y: 1}
	cam.Fovy = 45
	cam.Projection = rl.CameraPerspective
	scene.cam = cam
}
func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	rl.UpdateCamera(&scene.cam,rl.CameraOrbital) // Orbit cam around target
	// DRAW CUBE
	rl.BeginMode3D(scene.cam)
	rl.DrawCube(rl.Vector3{}, 2, 2, 2, rl.Red)
	rl.EndMode3D()

	rl.DrawText("Holy moly!!", int32(rl.GetScreenWidth()/8), 100, 80, rl.White)
	rl.DrawText("Press Enter to go back", int32(rl.GetScreenWidth()/8), 200, 30, rl.DarkGreen)
	if rl.IsKeyPressed(rl.KeyEnter) {
		return true
	}
	return false
}
func (scene *Scene) Unload(ctx engine.Context) string {
	return "start" // go back to main menu (start)
}
