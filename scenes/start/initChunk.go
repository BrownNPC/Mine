package start

import (
	c "GameFrameworkTM/components"
	"GameFrameworkTM/components/Blocks"
)

func InitChunk() *c.Chunk {
	var chunk c.Chunk
	for x := range c.CHUNK_SIZE {
		for y := range c.CHUNK_SIZE {
			for z := range c.CHUNK_SIZE {
				chunk.Set(x, y, z, Blocks.Type(x+y+z))
			}
		}
	}
	return &chunk
}
