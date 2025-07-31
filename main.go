package main

import (
	"GameFrameworkTM/engine"
	"GameFrameworkTM/scenes"
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// You can register scenes in scenes/register.go

// You can edit the window title in this file.
func main() {
	rl.SetTraceLogLevel(rl.LogError)
	err := engine.Run(scenes.Registered, engine.Config{
		WindowTitle: "change this in main.go",
	})
	if err != nil {
		fmt.Println(err)
	}
}
