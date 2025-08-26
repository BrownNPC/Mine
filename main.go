package main

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/engine"
	"GameFrameworkTM/scenes"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// You can register scenes in scenes/register.go

// You can edit the window title in this file.
func main() {
	rl.SetTraceLogLevel(rl.LogDebug)
	err := engine.Run(scenes.Registered, engine.Config{
		WindowTitle:   "Mine",
		MinScreenSize: c.V2(640, 480),
		ExitKey:       rl.KeyNumLock,
	})
	if err != nil {
		fmt.Println(err)
	}
}
