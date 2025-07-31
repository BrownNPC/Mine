# Raylib Golang "GameFramework"

A simple template for making games with raylib in 300 lines. You are expected
to modify this codebase to suit your needs.


Featuring Seperate "Scenes" with their own state


### Getting started

Just create a repository using this as the template. By default it includes a spinning cube demo. and a main menu. See below on how to run it.
#### Running
If you are on windows, you only need to place `raylib.dll` next to the `main.go` file. See [PureGo](https://github.com/gen2brain/raylib-go/?tab=readme-ov-file#purego-without-cgo-ie-cgo_enabled0) 


For linux and mac you will need some dependencies.
See [raylib-go](https://github.com/gen2brain/raylib-go/) for instructions. 

once you have the dependencies, run `go mod tidy` and then `go run .`


# Adding Scenes
A scene is just a struct that holds the data for your game, updates the logic and draws things to the screen.
```go
// a scene must implement the engine.scene interface
type scene interface {
	Load(Context)                        // called when this Scene is switched to
	Update(Context) (unload bool)        // called every frame
	Unload(Context) (nextSceneID string) // called after Update returns true. Switches to nextSceneID
}
```

put your scenes inside of its own seperate package:
`scenes/<MyScenePackage>` for example:
```
scenes/
├── RegisteredScenes.go // register your scenes in this file (more info later)
| // cube scene package
├── cube 
│   └── cube.go // spinning cube
| // start scene package (loaded first)
└── start 
    ├── start.go
    └── systems.go
```

Scenes must be registered inside of [`scenes/register.go`](https://github.com/BrownNPC/Golang-Raylib-GameFramework/blob/master/scenes/register.go) 
```go
// register all your scenes in here
var Registered = engine.Scenes{
	"start": &start.Scene{},
	"cube": &cube.Scene{},
}
```


here is a blank scene that will never quit
```go
package myscene

import "GameFrameworkTM/engine"

type Scene struct {
}
// Load is called once the scene is switched to
func (scene *Scene) Load(ctx engine.Context) {
}
// update is called every frame
func (scene *Scene) Update(ctx engine.Context) (unload bool) {
	return false // if true is returned, Unload is called
}
// called after Update returns true
func (scene *Scene) Unload(ctx engine.Context) (nextSceneID string) {
	return "someOtherSceneId"  // the engine will switch to the scene that is registered with this id
}
```
### Demo scenes
#### [`scenes/cube/cube.go`](https://github.com/BrownNPC/Golang-Raylib-GameFramework/blob/master/scenes/cube/cube.go)
this scene features a spinning cube (wow!)

#### [`scenes/start/start.go`](https://github.com/BrownNPC/Golang-Raylib-GameFramework/blob/master/scenes/start/start.go)
this is the main menu scene that is loaded first, press enter to select an option.






By the way, engine.Context is supposed to be used by you to implement a feature. 
```go
// info to pass to scenes
// eg. a camera, game map, or save file
type Context struct {
	SomeData any
}
```
### Code description

#### `main.go`

The main.go file isn't important, you can only change the window title here.

#### [`engine/engine.go`](https://github.com/BrownNPC/Golang-Raylib-GameFramework/blob/master/engine/engine.go)

The [`engine.Run`](https://github.com/BrownNPC/Golang-Raylib-GameFramework/blob/master/engine/engine.go#L33)
function initializes a window and centers it. Its only 100 lines.
looking at the code comments is recommended. The engine.Context, which is supposed to be used by you to implement features is defined here.





### What next?
Add what you need. Here are some features you could implement:

- [Letterboxing (black bars)](https://www.raylib.com/examples/core/loader.html?name=core_window_letterbox)
- Asset management
- Save Files
- anything else your game may need

