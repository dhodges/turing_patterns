package images

import (
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// tsGrid a grid of values which change with each iteration using turing scale variations
// see: https://softologyblog.wordpress.com/2011/07/05/multi-scale-turing-patterns/
// and: http://www.jonathanmccabe.com/Cyclic_Symmetric_Multi-Scale_Turing_Patterns.pdf
type tsGrid struct {
	Width      int
	Height     int
	scales     []turingScale
	grid       [][]float64
	activators [][][]float64
	inhibitors [][][]float64
	variations [][][]float64
}

// turingScale one of more of these are used to change a grid of values with each iteration
type turingScale struct {
	ActivatorRadius int
	InhibitorRadius int
	SmallAmount     float64
	Weight          float64
	Symmetry        int
}

// DefaultTuringScales default values when we have no config
var defaultTuringScales = []turingScale{
	turingScale{20, 40, 0.04, 1, 2},
	turingScale{10, 20, 0.03, 1, 2},
	turingScale{5, 10, 0.02, 1, 2},
	turingScale{1, 2, 0.01, 1, 2},
}

// makeTuringScaleGrid create a default multi-scale turing grid from the given params
func makeTuringScaleGrid(width, height int, scales []turingScale) *tsGrid {
	return &tsGrid{
		Width:      width,
		Height:     height,
		scales:     scales,
		grid:       util.Make2DGridFloat64Randomised(width, height),
		activators: util.Make3DGridFloat64(width, height, len(scales)),
		inhibitors: util.Make3DGridFloat64(width, height, len(scales)),
		variations: util.Make3DGridFloat64(width, height, len(scales)),
	}
}

// NextIteration generate the next variation of this grid of values
func (grid tsGrid) NextIteration() {
	grid.calcNextVariations()
	grid.normaliseGridValues()
}

func (grid tsGrid) sampleXY(x, y, scaleNdx int) float64 {
	grid.activators[x][y][scaleNdx] = util.AverageOfPixelsWithinCircle(x, y, grid.scales[scaleNdx].ActivatorRadius, grid.grid)
	grid.activators[x][y][scaleNdx] *= grid.scales[scaleNdx].Weight

	grid.inhibitors[x][y][scaleNdx] = util.AverageOfPixelsWithinCircle(x, y, grid.scales[scaleNdx].InhibitorRadius, grid.grid)
	grid.inhibitors[x][y][scaleNdx] *= grid.scales[scaleNdx].Weight

	// the variation can be calculated as an average of values within an arbitrary radius from x,y
	// but instead we use a radius of one pixel, i.e. just the value at x,y
	// apparently a radius of one pixel produces "the sharpest, most detailed images"
	return math.Abs(grid.activators[x][y][scaleNdx] - grid.inhibitors[x][y][scaleNdx])
}

func (grid tsGrid) calcNextVariations() {
	gridCenterX, gridCenterY := grid.Width/2, grid.Height/2
	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			for k := 0; k < len(grid.scales); k++ {
				// if symmetry > 1 then we effectively average this variation
				// with that many samples taken from the same point (x, y) rotated around
				// the image centre - this should result in an image symmetric around
				// its center point
				symmetry := float64(grid.scales[k].Symmetry)
				variation := 0.0
				for n := symmetry; n > 0.0; n-- {
					x1, y1 := util.RotateAboutAngle(x, y, 360.0/n, gridCenterX, gridCenterY)
					x1 = util.ConstrainInt(0, x1, grid.Width)
					y1 = util.ConstrainInt(0, y1, grid.Height)
					variation += grid.sampleXY(x1, y1, k)
				}
				grid.variations[x][y][k] = variation / symmetry
			}

			// best variation will be the smallest
			var ( // begin with values that are arbitrary yet valid
				ndx               = 0
				bestVariation     = &grid.scales[0]
				smallestVariation = grid.variations[x][y][0]
			)
			for k := 0; k < len(grid.scales); k++ {
				if grid.variations[x][y][k] < smallestVariation {
					ndx = k
					bestVariation = &grid.scales[k]
				}
			}
			if grid.activators[x][y][ndx] > grid.inhibitors[x][y][ndx] {
				grid.grid[x][y] += bestVariation.SmallAmount
			} else {
				grid.grid[x][y] -= bestVariation.SmallAmount
			}
		}
	}
}

func (grid tsGrid) normaliseGridValues() {
	// normalise all grid values to scale them back between -1 and +1
	// begin with the min and max values across the grid

	var ( // begin with values that are arbitrary yet valid
		smallest, largest = grid.grid[0][0], grid.grid[0][0]
	)
	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			smallest = math.Min(smallest, grid.grid[x][y])
			largest = math.Max(largest, grid.grid[x][y])
		}
	}

	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			grid.grid[x][y] = (grid.grid[x][y]-smallest)/(largest-smallest)*2 - 1
		}
	}
}

// copyOfCurrentState return a copy of the current grid
func (grid tsGrid) copyOfCurrentState() [][]float64 {
	copy := util.Make2DGridFloat64(grid.Width, grid.Height)
	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			copy[x][y] = grid.grid[x][y]
		}
	}
	return copy
}
