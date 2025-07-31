package start

import (
	"GameFrameworkTM/engine"
	"os"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// start scene is the main menu
type Scene struct {
	menuItems        []string
	menuFontSize     int32
	selectedMenuItem int
}

func (scene *Scene) Load(ctx engine.Context) {
	scene.menuItems = []string{
		"See the CUBE", "Exit",
	}
	scene.menuFontSize = 80
}

func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	scene.selectedMenuItem =
		updateSelectedMenuItem(scene.selectedMenuItem, len(scene.menuItems)-1)

	rl.ClearBackground(rl.LightGray)

	drawMenuItems(scene.menuItems, scene.menuFontSize, scene.selectedMenuItem)
	if rl.IsKeyPressed(rl.KeyEnter) {
		return true // quit scene and call unload
	}
	return false
}

// called when update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	if scene.menuItems[scene.selectedMenuItem] == "Exit" {
		os.Exit(0)
	}
	return "cube"
}

func drawMenuItems(items []string, fontSize int32, selectedItem int) {

	xPos := int32(rl.GetScreenWidth() / 8) // 1/8th of window width
	for i, item := range items {
		var color = rl.DarkGray
		yPos := int32(i*100) + fontSize*2 // math probably aint mathin

		if i == selectedItem {
			item = "> " + item
			color = rl.Red
		}
		rl.DrawText(item,
			xPos, yPos,
			fontSize,
			color,
		)
	}
}
func updateSelectedMenuItem(selectedItem, numItems int) int {

	if rl.IsKeyPressed(rl.KeyUp) {
		selectedItem--
	}
	if rl.IsKeyPressed(rl.KeyDown) {
		selectedItem++
	}
	//clamp
	if selectedItem > numItems {
		selectedItem = 0
	} else if selectedItem < 0 {
		selectedItem = numItems
	}
	return selectedItem
}
