package images

import (
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// turingScaleGrid a grid of values which change with each iteration using turing scale variations
// see: https://softologyblog.wordpress.com/2011/07/05/multi-scale-turing-patterns/
// and: http://www.jonathanmccabe.com/Cyclic_Symmetric_Multi-Scale_Turing_Patterns.pdf
type turingScaleGrid struct {
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

// makeTuringScaleGrid create a default multi-scale turing grid of given width and weight
func makeTuringScaleGrid(width, height int) *turingScaleGrid {
	return &turingScaleGrid{
		Width:      width,
		Height:     height,
		scales:     DefaultConfig.Scales,
		grid:       util.Make2DGridFloat64Randomised(width, height),
		activators: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
		inhibitors: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
		variations: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
	}
}

// makeTuringScaleGrid create a default multi-scale turing grid from the given params
func makeTuringScaleGridFromConfig(cfg TuringScaleConfig) *turingScaleGrid {
	return &turingScaleGrid{
		Width:      cfg.Width,
		Height:     cfg.Height,
		scales:     cfg.Scales,
		grid:       util.Make2DGridFloat64Randomised(cfg.Width, cfg.Height),
		activators: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
		inhibitors: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
		variations: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
	}
}

// NextIteration generate the next variation of this grid of values
func (grid turingScaleGrid) NextIteration() {
	grid.calcNextVariations()
	grid.normaliseGridValues()
}

func (grid turingScaleGrid) calcNextVariations() {
	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			for k := 0; k < len(grid.scales); k++ {
				grid.activators[x][y][k] = util.AverageOfPixelsWithinCircle(x, y, grid.scales[k].ActivatorRadius, grid.grid)
				grid.activators[x][y][k] *= grid.scales[k].Weight

				grid.inhibitors[x][y][k] = util.AverageOfPixelsWithinCircle(x, y, grid.scales[k].InhibitorRadius, grid.grid)
				grid.inhibitors[x][y][k] *= grid.scales[k].Weight

				// the variation can be calculated as an average of values within an arbitrary radius from x,y
				// but instead we use a radius of one pixel, i.e. just the value at x,y
				// apparently a radius of one pixel produces "the sharpest, most detailed images"
				grid.variations[x][y][k] = math.Abs(grid.activators[x][y][k] - grid.inhibitors[x][y][k])
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

func (grid turingScaleGrid) normaliseGridValues() {
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
func (grid turingScaleGrid) copyOfCurrentState() [][]float64 {
	copy := util.Make2DGridFloat64(grid.Width, grid.Height)
	for x := 0; x < grid.Width; x++ {
		for y := 0; y < grid.Height; y++ {
			copy[x][y] = grid.grid[x][y]
		}
	}
	return copy
}
