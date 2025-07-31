package scenes

import (
	"GameFrameworkTM/engine"
	"GameFrameworkTM/scenes/cube"
	"GameFrameworkTM/scenes/start"
)

// register all your scenes in here
var Registered = engine.Scenes{
	"start": &start.Scene{},
	"cube": &cube.Scene{},
}
