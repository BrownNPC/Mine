package start

import (
	"GameFrameworkTM/components/Blocks"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"log"
	"os"
	"path/filepath"
)

type TexturePaths struct {
	Top, Bottom, Side string
}

const TexDimensions = 16

func getPathsForBlockID(blockID Blocks.Type) TexturePaths {
	const basepath = "assets/blocks/textures"
	if blockID == Blocks.Air {
		panic("No texture for block ID 0 (Air)")
	}
	base := blockID.String()
	return TexturePaths{
		Top:    filepath.Join(basepath, base+"_top.png"),
		Bottom: filepath.Join(basepath, base+"_bottom.png"),
		Side:   filepath.Join(basepath, base+"_side.png"),
	}
}

func CreateAtlas() image.Image {
	// Prepare texture paths for all blocks except Air
	var textures [Blocks.TotalBlocks + 1]TexturePaths
	for i := range Blocks.TotalBlocks {
		if i == Blocks.Air {
			continue
		}
		texPath := getPathsForBlockID(Blocks.Type(i))
		if texPath == (TexturePaths{}) {
			log.Panicf("path for block ID %s is empty", Blocks.Type(i).String())
		}
		textures[i] = texPath
	}

	const atlasWidth = TexDimensions * 3
	const atlasHeight = (int(Blocks.TotalBlocks)) * TexDimensions

	finalAtlas := image.NewRGBA(image.Rect(0, 0, atlasWidth, atlasHeight))

	for i, texPath := range textures {
		if i == int(Blocks.Air) || texPath == (TexturePaths{}) {
			continue
		}

		loaded := loadTextures(texPath)

		// Calculate rectangles for top, bottom, side
		topRect := image.Rect(0, TexDimensions*i, TexDimensions, TexDimensions*(i+1))
		bottomRect := image.Rect(TexDimensions, TexDimensions*i, TexDimensions*2, TexDimensions*(i+1))
		sideRect := image.Rect(TexDimensions*2, TexDimensions*i, TexDimensions*3, TexDimensions*(i+1))

		// Draw them on the atlas
		draw.Draw(finalAtlas, topRect, loaded.Top, image.Point{}, draw.Over)
		draw.Draw(finalAtlas, bottomRect, loaded.Bottom, image.Point{}, draw.Over)
		draw.Draw(finalAtlas, sideRect, loaded.Side, image.Point{}, draw.Over)
	}

	return finalAtlas
}

// images loaded into memory
type LoadedTextures struct {
	Top, Bottom, Side image.Image
}

func loadTextures(paths TexturePaths) LoadedTextures {
	_top, err := os.Open(paths.Top)
	check(err, paths.Top)
	defer _top.Close()

	_bottom, err := os.Open(paths.Bottom)
	check(err, paths.Bottom)
	defer _bottom.Close()

	_side, err := os.Open(paths.Side)
	check(err, paths.Side)
	defer _side.Close()

	top, err := png.Decode(_top)
	if err != nil {
		log.Panicf("failed to decode image %s: %v", paths.Top, err)
	}

	bottom, err := png.Decode(_bottom)
	if err != nil {
		log.Panicf("failed to decode image %s: %v", paths.Bottom, err)
	}

	side, err := png.Decode(_side)
	if err != nil {
		log.Panicf("failed to decode image %s: %v", paths.Side, err)
	}

	return LoadedTextures{
		Top:    top,
		Bottom: bottom,
		Side:   side,
	}
}

func check(err error, path string) {
	if err != nil {
		log.Panic(fmt.Errorf("error for path: %s: %w", path, err))
	}
}
