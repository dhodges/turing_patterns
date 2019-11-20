package util

import (
	"math"
)

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

// RotateAboutCenter rotates the point (x, y) by the given angle around the centre (0, 0)
func RotateAboutCenter(x, y int, angle float64) (x1, y1 int) {
	angle *= math.Pi / 180
	sin := math.Sin(angle)
	cos := math.Cos(angle)
	x1 = int(float64(x)*cos - float64(y)*sin)
	y1 = int(float64(y)*cos + float64(x)*sin)

	return x1, y1
}

// RotateAboutAngle rotates the point (x, y) by the given angle around the centre point (xc, yc)
func RotateAboutAngle(x, y int, angle float64, xc, yc int) (x1, y1 int) {
	x1, y1 = RotateAboutCenter(x-xc, y-yc, angle)
	x1 = x1 + xc
	y1 = y1 + yc

	return x1, y1
}
