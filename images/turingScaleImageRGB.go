package images

import (
	"image/color"
	"math"

	"github.com/dhodges/turing_patterns/hsb"
	"github.com/dhodges/turing_patterns/util"
)

// TuringScaleImageRGB an RGB reaction/diffusion image (using turing scales)
type TuringScaleImageRGB struct {
	grid   *turingScaleGrid
	colors [][]hsb.NHSBA
}

// MakeTuringScaleImageRGB returns a TuringScaleImageRGB with default values
func MakeTuringScaleImageRGB(width, height int) *TuringScaleImageRGB {
	img := &TuringScaleImageRGB{
		grid:   makeTuringScaleGrid(width, height),
		colors: util.Make2DGridNHSBA(width, height),
	}
	// NB: store all colors as HSB, defaulting to a random hue
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.colors[x][y] = hsb.NHSBA{H: util.RandFloat64(0.0, 360.0), S: 1.0, B: 0.5}
		}
	}
	return img
}

// MakeTuringScaleImageRGBFromConfig return a TuringScaleImageRGB with the given config
func MakeTuringScaleImageRGBFromConfig(cfg TuringScaleConfig) *TuringScaleImageRGB {
	img := &TuringScaleImageRGB{
		grid:   makeTuringScaleGridFromConfig(cfg),
		colors: util.Make2DGridNHSBA(cfg.Width, cfg.Height),
	}
	// NB: store all colors as HSB values, defaulting to a muted teal
	for x := 0; x < img.grid.Width; x++ {
		for y := 0; y < img.grid.Height; y++ {
			img.colors[x][y] = hsb.NHSBA{H: 180.0, S: 0.5, B: 0.5}
		}
	}
	return img
}

// NextIteration generates the next variation of this image
func (img TuringScaleImageRGB) NextIteration() {

	// we are interested in the change from the previous iteration to the next
	previousGrid := img.copyOfCurrentState()

	img.grid.NextIteration()

	for x := 0; x < img.grid.Width; x++ {
		for y := 0; y < img.grid.Height; y++ {
			delta := img.grid.grid[x][y] - previousGrid[x][y]
			img.colors[x][y] = updateColor(img.colors[x][y], delta)
		}
	}
}

// c: the previous version of this color
// delta: [-1.0 <= delta <= 1.0], indicating the change in color
func updateColor(c hsb.NHSBA, delta float64) hsb.NHSBA {
	newColor := hsb.NHSBA{H: c.H, S: 1.0, B: 0.5, A: 1.0}

	delta = toFixed(delta*100, 2)

	if delta != 0.0 {
		newColor.H = ((delta + 1) / 2) * 360.0
	}

	return newColor
}

// c: the previous version of this color
// delta: [-1.0 <= delta <= 1.0], indicating the change in color
func updateColorGrainy(c hsb.NHSBA, delta float64) hsb.NHSBA {
	newColor := hsb.NHSBA{H: c.H, S: c.S, B: c.B, A: c.A}

	delta = toFixed(delta*100, 2)

	if delta != 0.0 {
		if 0.0 < c.H && c.H < 360.0 {
			newColor.H = constrain(0.0, c.H+delta, 360.0)
		} else if 0.0 < c.S && c.S < 1.0 {
			newColor.S = constrain(0.0, c.S+delta, 1.0)
		} else if 0.0 < c.B && c.B < 1.0 {
			newColor.B = constrain(0.0, c.B+delta, 1.0)
		}
	}

	return newColor
}

// see: https://stackoverflow.com/a/29786394
func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

// see: https://stackoverflow.com/a/29786394
func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

// copyOfCurrentState return a copy of the current grid
func (img TuringScaleImageRGB) copyOfCurrentState() [][]float64 {
	return img.grid.copyOfCurrentState()
}

// OutputPNG generate a PNG file from the current iteration
func (img TuringScaleImageRGB) OutputPNG(filename string) {
	util.OutputPNG(filename, img.pixmap())
}

// pixmap return a grayscale pixmap derived from the current state of grid values
func (img TuringScaleImageRGB) pixmap() [][]color.NRGBA {
	pixels := util.Make2DGridNRGBA(img.grid.Width, img.grid.Height)

	// map all grid values to a pixel grayscale value
	for x := 0; x < img.grid.Width; x++ {
		for y := 0; y < img.grid.Height; y++ {
			pixels[x][y] = *img.colors[x][y].ToNRGBA()
		}
	}
	return pixels
}

// constrain the given float within the range min <= n <= max
func constrain(min, n, max float64) float64 {
	n = math.Max(min, n)
	n = math.Min(n, max)
	return n
}
