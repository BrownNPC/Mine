package Blocks

type Type uint8

//go:generate stringer -type=BlockType
const (
	Air Type = iota
	Dirt
	Grass
	len
)
const TotalBlocks = len

type Direction uint8

const (
	Top Direction = iota
	Bottom
	Right
	Left
	Back
	Front
)
