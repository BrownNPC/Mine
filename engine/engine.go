package engine

import (
	"errors"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// config is passed to the Run function in main.go
type Config struct {
	//for implementing letterboxing (black bars) see:https://www.raylib.com/examples/core/loader.html?name=core_window_letterbox
	// VirtualWidth, VirtualHeight int
	WindowTitle string
}

// info to pass to scenes
// eg. a camera, game map, or save file
type Context struct {
	SomeData any
}

// a scene must implement these methods
type scene interface {
	Load(Context)                        // called when this Scene is switched to
	Update(Context) (unload bool)        // called every frame
	Unload(Context) (nextSceneID string) // called after Update returns true
}

// map from string id to a Scene
type Scenes map[string]scene

func Run(scenes Scenes, cfg Config) error {
	ActiveSceneId := "start" // look for a scene named start as entry-point
	ActiveScene, ok := scenes[ActiveSceneId]
	ctx := Context{} // info to pass to scenes.
	if !ok {
		return errors.New(`Cannot start. There must be a scene with id "start" that is the entry-point`)
	} else if ActiveScene == nil {
		return errors.New("start scene cannot be nil")
	}
	// --------------BEGIN--------------
	rl.SetConfigFlags(rl.FlagWindowResizable)
	rl.InitWindow(0, 0, cfg.WindowTitle)
	rl.InitAudioDevice()
	defer rl.CloseWindow() // de-initialization
	defer rl.CloseAudioDevice()
	// -----------------------CENTER WINDOW----------------------------
	WindowWidth, WindowHeight := (rl.GetScreenWidth()*90)/100, (rl.GetScreenHeight()*90)/100
	rl.SetWindowSize(WindowWidth, WindowHeight) //90% of screen
	centerWindow()
	// ----LOAD START SCENE----
	ActiveScene.Load(ctx)
	// ----MAIN LOOP----
	for !rl.WindowShouldClose() {
		// ----FULL SCREEN ON F11----
		if rl.IsKeyPressed(rl.KeyF11) {
			rl.ToggleBorderlessWindowed()
		}
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		// -------UPDATE SCENE---------
		var unloadActiveScene bool = ActiveScene.Update(ctx)
		rl.EndDrawing()
		if unloadActiveScene {
			// -------UNLOAD SCENE-------
			var nextSceneId string = ActiveScene.Unload(ctx) // unload returns nextSceneId
			var nextScene, ok = scenes[nextSceneId]
			// ------SWITCH SCENE------
			if ok && nextScene != nil {
				ActiveSceneId = nextSceneId
				ActiveScene = nextScene
				ActiveScene.Load(ctx)
				continue
			}
			//-----ERROR HANDLING------
			if !ok {
				return fmt.Errorf(`There is no scene with id "%s", tried switching from scene "%s"`, nextSceneId, ActiveSceneId)
			} else if nextScene == nil {
				return fmt.Errorf(`scene with id "%s" is nil, tried switching from scene "%s"`, nextSceneId, ActiveSceneId)
			}
		}
	}
	return nil
}
func centerWindow() {
	WindowWidth, WindowHeight := rl.GetScreenWidth(), rl.GetScreenHeight()
	monitorWidth := rl.GetMonitorWidth(rl.GetCurrentMonitor())
	monitorHeight := rl.GetMonitorHeight(rl.GetCurrentMonitor())
	xPos := (monitorWidth - WindowWidth) / 2
	yPos := (monitorHeight - WindowHeight) / 2
	rl.SetWindowPosition(xPos, yPos)
}
