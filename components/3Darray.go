package c

import (
	"fmt"
	"log"
)

// ThreeDimensionalArray is a generic 3D array type backed by a flat slice.
// It stores dimensions X × Y × Z and a one-dimensional data slice of length X*Y*Z.
type ThreeDimensionalArray[T any] struct {
	X, Y, Z int
	data    []T
}

// New3dArray allocates a new ThreeDimensionalArray of dimensions x, y, z.
// The backing slice is zero-initialized with length x*y*z.
func New3dArray[T any](x, y, z int) ThreeDimensionalArray[T] {
	if x < 0 || y < 0 || z < 0 {
		panic(fmt.Sprintf("invalid dimensions: %dx%dx%d", x, y, z))
	}
	size := x * y * z
	return ThreeDimensionalArray[T]{
		X:    x,
		Y:    y,
		Z:    z,
		data: make([]T, size),
	}
}

// idx computes the flat index for coordinates (x, y, z) or returns an error.
func (a ThreeDimensionalArray[T]) idx(x, y, z int) (int, error) {
	if x < 0 || x >= a.X {
		return 0, fmt.Errorf("x index out of bounds: %d", x)
	}
	if y < 0 || y >= a.Y {
		return 0, fmt.Errorf("y index out of bounds: %d", y)
	}
	if z < 0 || z >= a.Z {
		return 0, fmt.Errorf("z index out of bounds: %d", z)
	}
	return x*(a.Y*a.Z) + y*a.Z + z, nil
}

// Get returns the element at (x, y, z) or panics if out of bounds.
func (a ThreeDimensionalArray[T]) Get(x, y, z int) T {
	i, err := a.idx(x, y, z)
	if err != nil {
		log.Panic(err)
	}
	return a.data[i]
}

// GetRef returns a pointer to the element at (x, y, z) or panics if out of bounds.
func (a ThreeDimensionalArray[T]) GetRef(x, y, z int) *T {
	i, err := a.idx(x, y, z)
	if err != nil {
		log.Panic(err)
	}
	return &a.data[i]
}

// Set assigns value at (x, y, z) or panics if out of bounds.
func (a *ThreeDimensionalArray[T]) Set(x, y, z int, value T) {
	i, err := a.idx(x, y, z)
	if err != nil {
		log.Panic(err)
	}
	a.data[i] = value
}

// Dimensions returns the size of each axis: X, Y, Z.
func (a ThreeDimensionalArray[T]) Dimensions() (x, y, z int) {
	return a.X, a.Y, a.Z
}

// get the flat 1d array backing the 3d array
func (a ThreeDimensionalArray[T]) BackingArray() []T {
	return a.data
}
