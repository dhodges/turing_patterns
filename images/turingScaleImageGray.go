package images

import (
	"image/color"
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// TuringScaleImageGray a grayscale reaction/diffusion image (using turing scales)
type TuringScaleImageGray struct {
	grid *turingScaleGrid
}

// MakeTuringScaleImageGray return a TuringScaleImageGray with default values
func MakeTuringScaleImageGray(width, height int) *TuringScaleImageGray {
	return &TuringScaleImageGray{grid: makeTuringScaleGrid(width, height)}
}

// MakeTuringScaleImageGrayFromConfig return a TuringScaleImageGray with the given config
func MakeTuringScaleImageGrayFromConfig(cfg TuringScaleConfig) *TuringScaleImageGray {
	return &TuringScaleImageGray{grid: makeTuringScaleGridFromConfig(cfg)}
}

// NextIteration generate the next variation of this image
func (img TuringScaleImageGray) NextIteration() {
	img.grid.NextIteration()
}

// OutputPNG generate a PNG file from the current iteration
func (img TuringScaleImageGray) OutputPNG(filename string) {
	util.OutputPNG(filename, img.pixmap())
}

// pixmap return a grayscale pixmap derived from the current state of grid values
func (img TuringScaleImageGray) pixmap() [][]color.NRGBA {
	pixels := util.Make2DGridNRGBA(img.grid.Width, img.grid.Height)

	// map all grid values to a pixel grayscale value
	for x := 0; x < img.grid.Width; x++ {
		for y := 0; y < img.grid.Height; y++ {
			gray := uint8(math.Trunc((img.grid.grid[x][y] + 1) / 2 * 255))
			pixels[x][y] = color.NRGBA{
				R: uint8(gray),
				G: uint8(gray),
				B: uint8(gray),
				A: 255,
			}
		}
	}
	return pixels
}
