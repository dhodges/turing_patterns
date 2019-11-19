package images

import (
	"image/color"
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// TSImageGray a grayscale reaction/diffusion image (using turing scales)
type TSImageGray struct {
	grid *tsGrid
}

// MakeTSImageGray return a TSImageGray with default values
func MakeTSImageGray(width, height int) *TSImageGray {
	return &TSImageGray{grid: makeTuringScaleGrid(width, height)}
}

// MakeTSImageGrayFromConfig return a TSImageGray with the given config
func MakeTSImageGrayFromConfig(cfg TuringScaleConfig) *TSImageGray {
	return &TSImageGray{grid: makeTuringScaleGridFromConfig(cfg)}
}

// NextIteration generate the next variation of this image
func (img TSImageGray) NextIteration() {
	img.grid.NextIteration()
}

// OutputPNG generate a PNG file from the current iteration
func (img TSImageGray) OutputPNG(filename string) {
	util.OutputPNG(filename, img.pixmap())
}

// pixmap return a grayscale pixmap derived from the current state of grid values
func (img TSImageGray) pixmap() [][]color.NRGBA {
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
