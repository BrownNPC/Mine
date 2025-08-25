package start

import (
	"GameFrameworkTM/components/Blocks"
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
	// each block has 3 textures.
	// Block_top.png
	// Block_side.png
	// Block_bottom.png
	top := filepath.Join(basepath, base+"_top.png")
	side := filepath.Join(basepath, base+"_side.png")
	bottom := filepath.Join(basepath, base+"_bottom.png")
	return TexturePaths{
		Top:    top,
		Bottom: bottom,
		Side:   side,
	}
}

func CreateAtlas() image.Image {
	// +1 since our for loop starts counting at 1.
	// 0 is skipped because it's air, and does not have a texture.
	// this is the same as using a map[Blocks.Type]TexturePaths
	// but more optimized.
	var textures [Blocks.TotalBlocks + 1]TexturePaths
	for i := range Blocks.TotalBlocks {
		if i == Blocks.Air {
			continue
		}
		textures[i] = getPathsForBlockID(i)
	}
	// textures are 16px
	// 3 vairants of textures for each block
	const atlasWidth = TexDimensions * 3
	// number of blocks x 16px height
	const atlasHeight = int(Blocks.TotalBlocks) * TexDimensions

	var finalAtlas = image.NewRGBA(image.Rect(0, 0, atlasWidth, atlasHeight))

	// load the textures as images
	for i, texPath := range textures {
		textures := loadTextures(texPath)
		// Y axis is block type
		// X axis are faces
		// in this order:
		// top, bottom,side
		topRect := image.Rect(0, TexDimensions*i, TexDimensions, TexDimensions*(i+1))
		bottomRect := image.Rect(TexDimensions, TexDimensions*i, TexDimensions*2, TexDimensions*(i+1))
		sideRect := image.Rect(TexDimensions*2, TexDimensions*i, TexDimensions*3, TexDimensions*(i+1))

		// draw texture variant onto atlas
		var Z = image.Pt(0, 0)
		draw.Draw(finalAtlas, topRect, textures.Top, Z, draw.Over)
		draw.Draw(finalAtlas, bottomRect, textures.Bottom, Z, draw.Over)
		draw.Draw(finalAtlas, sideRect, textures.Side, Z, draw.Over)
	}
	return finalAtlas
}

// images loaded into memory
type LoadedTextures struct {
	Top, Bottom, Side image.Image
}

func loadTextures(paths TexturePaths) LoadedTextures {
	// load files
	_top, err := os.Open(paths.Top)
	check(err)
	defer _top.Close()

	_bottom, err := os.Open(paths.Bottom)
	check(err)
	defer _bottom.Close()
	_side, err := os.Open(paths.Side)

	check(err)
	defer _side.Close()

	// decode images
	top, err := png.Decode(_top)
	check(err)
	bottom, err := png.Decode(_bottom)
	check(err)
	side, err := png.Decode(_side)
	check(err)

	// return loaded images
	return LoadedTextures{
		Top:    top,
		Bottom: bottom,
		Side:   side,
	}
}
func check(err error) {
	log.Panic("Failed to load textures for block", err)
}
