package util

import (
	"image/color"
	"math/rand"

	"github.com/dhodges/turing_patterns/hsb"
)

// Make2DGridFloat64 make a 2D array of float64
func Make2DGridFloat64(width, height int) [][]float64 {
	grid := make([][]float64, height)
	for i := range grid {
		grid[i] = make([]float64, width)
	}
	return grid
}

// Make2DGridUInt8 make a 2D array of uint8
func Make2DGridUInt8(width, height int) [][]uint8 {
	grid := make([][]uint8, height)
	for i := range grid {
		grid[i] = make([]uint8, width)
	}
	return grid
}

// Make2DGridNRGBA make a 2D array of NRGBA colors
func Make2DGridNRGBA(width, height int) [][]color.NRGBA {
	grid := make([][]color.NRGBA, height)
	for i := range grid {
		grid[i] = make([]color.NRGBA, width)
	}
	return grid
}

// Make2DGridNHSBA make a 2D array of NHSBA colors
func Make2DGridNHSBA(width, height int) [][]hsb.NHSBA {
	grid := make([][]hsb.NHSBA, height)
	for i := range grid {
		grid[i] = make([]hsb.NHSBA, width)
	}
	return grid
}

// Make2DGridFloat64Randomised make a 2D array of random float64 values
func Make2DGridFloat64Randomised(width, height int) [][]float64 {
	grid := Make2DGridFloat64(width, height)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			grid[i][j] = RandFloat64(-1.0, 1.0)
		}
	}
	return grid
}

// Make3DGridUInt8 make a 3D array of UInt8
func Make3DGridUInt8(width, height, depth int) [][][]uint8 {
	grid := make([][][]uint8, height)
	for i := range grid {
		grid[i] = make([][]uint8, width)
		for j := range grid[i] {
			grid[i][j] = make([]uint8, depth)
		}
	}
	return grid
}

// Make3DGridFloat64 make a 3D array of float64
func Make3DGridFloat64(width, height, depth int) [][][]float64 {
	grid := make([][][]float64, height)
	for i := range grid {
		grid[i] = make([][]float64, width)
		for j := range grid[i] {
			grid[i][j] = make([]float64, depth)
		}
	}
	return grid
}

// RandFloat64 generate a random float64 between the given min and max
func RandFloat64(min, max float64) float64 {
	return min + rand.Float64()*(max-min)
}

func sqr(x int) int {
	return x * x
}

// PointIsWithinCircle does the given point exist within (or on) the given circle?
func PointIsWithinCircle(xp, yp, x, y, radius int) bool {
	// xp, yp: the point
	// x, y, radius: the circle

	// radius >= Math.sqrt((xp - x)² + (yp - y)²)
	// i.e.
	// radius² >= (xp - x)² + (yp - y)²

	return sqr(radius) >= sqr(xp-x)+sqr(yp-y)
}

// AverageOfPixelsWithinCircle return average of all pixel values in the given circle
func AverageOfPixelsWithinCircle(x, y, radius int, grid [][]float64) float64 {
	// x, y, radius: the circle of values from which to derive an average
	// grid: the grid of values from which the circles are found
	sum := 0.0
	numPixelsWithinCircle := 1.0

	for i := x - radius; i < x+radius; i++ {
		for j := y - radius; j < y+radius; j++ {

			// only include pixel values within the image bounds
			if (i >= 0) && (i < len(grid[0])) &&
				(j >= 0) && (j < len(grid)) {

				if PointIsWithinCircle(i, j, x, y, radius) {
					sum += grid[i][j]
					numPixelsWithinCircle++
				}
			}
		}
	}
	return sum / numPixelsWithinCircle
}
