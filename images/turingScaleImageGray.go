package images

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// TSImageGray a grayscale reaction/diffusion image (using turing scales)
type TSImageGray struct {
	grid *tsGrid
}

// TSImageConfigGray parameters
type TSImageConfigGray struct {
	Width  int
	Height int
	Scales []turingScale
}

// MakeTSImageGray return a TSImageGray with default values
func MakeTSImageGray(width, height int) *TSImageGray {
	return &TSImageGray{grid: makeTuringScaleGrid(width, height, defaultTuringScales)}
}

// ConfigFromFile configures TSImageGray from the given file
func (img TSImageGray) ConfigFromFile(configfile string) *TSImageGray {
	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}

	config := TSImageConfigGray{}
	if err = json.Unmarshal([]byte(file), &config); err != nil {
		log.Fatal(err)
	}

	img.initFromConfig(config)

	return &img
}

// initFromConfig configures TSImageGray from the given config
func (img TSImageGray) initFromConfig(cfg TSImageConfigGray) {
	img.grid = makeTuringScaleGrid(cfg.Width, cfg.Height, cfg.Scales)
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
