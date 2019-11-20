package images

import (
	"encoding/json"
	"image/color"
	"io/ioutil"
	"log"

	"github.com/dhodges/turing_patterns/hsb"
	"github.com/dhodges/turing_patterns/util"
)

// TSImageRGB an RGB reaction/diffusion image (using turing scales)
type TSImageRGB struct {
	grid   *tsGrid
	colors [][]hsb.NHSBA
}

// TSImageConfigRGB parameters that define the image
type TSImageConfigRGB struct {
	Width  int
	Height int
	Scales []turingScale
}

// MakeTSImageRGB returns a TSImageRGB with default values
func MakeTSImageRGB(width, height int) *TSImageRGB {
	img := &TSImageRGB{
		grid:   makeTuringScaleGrid(width, height, defaultTuringScales),
		colors: util.Make2DGridNHSBA(width, height),
	}
	// NB: store all colors as HSB, defaulting to a random hue
	for x := 0; x < width; x++ {
		for y := 0; y < height; y++ {
			img.colors[x][y] = hsb.NHSBA{H: util.RandFloat64(0.0, 360.0), S: 0.5, B: 1.0}
		}
	}
	return img
}

// ConfigFromFile configures TSImageRGB from the given file
func (img TSImageRGB) ConfigFromFile(configfile string) {

	file, err := ioutil.ReadFile(configfile)
	if err != nil {
		log.Fatal(err)
	}

	config := TSImageConfigRGB{}
	if err = json.Unmarshal([]byte(file), &config); err != nil {
		log.Fatal(err)
	}

	img.initFromConfig(config)
}

// initFromConfig configures TSImageRGB from the given config
func (img TSImageRGB) initFromConfig(cfg TSImageConfigRGB) {
	img.grid = makeTuringScaleGrid(cfg.Width, cfg.Height, cfg.Scales)
}

// NextIteration generates the next variation of this image
func (img TSImageRGB) NextIteration() {

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
	deltaHue := ((delta + 1) / 2) * 360.0

	if deltaHue > newColor.H {
		newColor.H += delta * 10
		newColor.S += delta * 10
	} else {
		newColor.H -= delta * 10
		newColor.S -= delta * 10
	}

	newColor.H = util.Constrain(0.0, newColor.H, 360.0)
	newColor.S = util.Constrain(0.0, newColor.S, 1.0)

	return newColor
}

// c: the previous version of this color
// delta: [-1.0 <= delta <= 1.0], indicating the change in color
func updateColorGrainy(c hsb.NHSBA, delta float64) hsb.NHSBA {
	newColor := hsb.NHSBA{H: c.H, S: c.S, B: c.B, A: c.A}

	delta = util.ToFixed(delta*100, 2)

	if delta != 0.0 {
		if 0.0 < c.H && c.H < 360.0 {
			newColor.H = util.Constrain(0.0, c.H+delta, 360.0)
		} else if 0.0 < c.S && c.S < 1.0 {
			newColor.S = util.Constrain(0.0, c.S+delta, 1.0)
		} else if 0.0 < c.B && c.B < 1.0 {
			newColor.B = util.Constrain(0.0, c.B+delta, 1.0)
		}
	}

	return newColor
}

// copyOfCurrentState return a copy of the current grid
func (img TSImageRGB) copyOfCurrentState() [][]float64 {
	return img.grid.copyOfCurrentState()
}

// OutputPNG generate a PNG file from the current iteration
func (img TSImageRGB) OutputPNG(filename string) {
	util.OutputPNG(filename, img.pixmap())
}

// pixmap return a grayscale pixmap derived from the current state of grid values
func (img TSImageRGB) pixmap() [][]color.NRGBA {
	pixels := util.Make2DGridNRGBA(img.grid.Width, img.grid.Height)

	// map all grid values to a pixel grayscale value
	for x := 0; x < img.grid.Width; x++ {
		for y := 0; y < img.grid.Height; y++ {
			pixels[x][y] = *img.colors[x][y].ToNRGBA()
		}
	}
	return pixels
}
