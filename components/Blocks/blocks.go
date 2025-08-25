package Blocks

type Type uint8

//go:generate stringer -type=Type
const (
	Air Type = iota
	Dirt
	Grass
	PalmLeaves
	PalmPlanks
	PalmLog
	_len
)
const TotalBlocks = _len

type Direction uint8

const (
	Top Direction = iota
	Bottom
	Right
	Left
	Back
	Front
)
