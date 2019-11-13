package util

import (
	"image"
	"testing"
)

func TestMake2DGridFloat64(t *testing.T) {
	width, height := 10, 20
	grid := Make2DGridFloat64(width, height)
	if len(grid) != height {
		t.Errorf("grid has height %d, but it should be %d", len(grid), height)
	}
	if len(grid[0]) != width {
		t.Errorf("grid has width %d, but it should be %d", len(grid[0]), width)
	}
}

func TestMake3DGridFloat64(t *testing.T) {
	width, height, depth := 11, 22, 33
	grid := Make3DGridFloat64(width, height, depth)
	if len(grid) != height {
		t.Errorf("grid has height %d, but it should be %d", len(grid), height)
	}
	if len(grid[0]) != width {
		t.Errorf("grid has width %d, but it should be %d", len(grid[0]), width)
	}
	if len(grid[0][0]) != depth {
		t.Errorf("grid has depth %d, but it should be %d", len(grid[0][0]), depth)
	}
}

func TestMake2DGridUInt8(t *testing.T) {
	width, height := 27, 51
	grid := Make2DGridUInt8(width, height)
	if len(grid) != height {
		t.Errorf("grid has height %d, but it should be %d", len(grid), height)
	}
	if len(grid[0]) != width {
		t.Errorf("grid has width %d, but it should be %d", len(grid[0]), width)
	}
}

func TestPointIsWithinCircle(t *testing.T) {
	x, y := 0, 0
	radius := 10

	pointsInside := []image.Point{
		image.Point{X: 1, Y: 1},
		image.Point{X: 3, Y: 3},
		image.Point{X: 5, Y: 5},
	}
	for _, pt := range pointsInside {
		if !PointIsWithinCircle(pt.X, pt.Y, x, y, radius) {
			t.Errorf("Circle(x:%d, y:%d, radius:%d) should contain Point(%d, %d)", x, y, radius, pt.X, pt.Y)
		}
	}

	pointsOutside := []image.Point{
		image.Point{X: 51, Y: 51},
		image.Point{X: 53, Y: 53},
		image.Point{X: 55, Y: 55},
	}
	for _, pt := range pointsOutside {
		if PointIsWithinCircle(pt.X, pt.Y, x, y, radius) {
			t.Errorf("Circle(x:%d, y:%d, radius:%d) should not contain Point(%d, %d)", x, y, radius, pt.X, pt.Y)
		}
	}
}

func TestAverageOfPixelsWithinCircle(t *testing.T) {
	grid := Make2DGridFloat64(100, 100)

	for y := 0; y < 10; y++ {
		for x := 5; x < 10; x++ {
			grid[x][y] = 1.0
		}
	}
	expected := 0.55
	average := AverageOfPixelsWithinCircle(5, 5, 5, grid)
	if average != expected {
		t.Errorf("average of circle(x:%d, y:%d, radius:%d) is %.2f, should be %.2f", 5, 5, 5, average, expected)
	}
}
