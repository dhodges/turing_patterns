package images

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"math"

	"github.com/dhodges/turing_patterns/util"
)

// MultiScaleImage a turing pattern image that iterates using multiple turing scale variations
// see: https://softologyblog.wordpress.com/2011/07/05/multi-scale-turing-patterns/
// and: http://www.jonathanmccabe.com/Cyclic_Symmetric_Multi-Scale_Turing_Patterns.pdf
type MultiScaleImage struct {
	Width      int
	Height     int
	scales     []turingScale
	grid       [][]float64
	activators [][][]float64
	inhibitors [][][]float64
	variations [][][]float64
}

// Configuration parameters that define an image, which we optionally expect to read from a json file
type Configuration struct {
	Width  int
	Height int
	Scales []turingScale
}

type turingScale struct {
	ActivatorRadius int
	InhibitorRadius int
	SmallAmount     float64
	Weight          float64
	Symmetry        int
}

// DefaultConfig MultiScaleImages will default to these params
var DefaultConfig Configuration = Configuration{
	Scales: []turingScale{
		//turingScale{100, 200, 0.05, 1, 3},
		turingScale{20, 40, 0.04, 1, 2},
		turingScale{10, 20, 0.03, 1, 2},
		turingScale{5, 10, 0.02, 1, 2},
		turingScale{1, 2, 0.01, 1, 2},
	},
}

// MakeMultiScaleImage create a default multi-scale turing pattern image of given width and cfg.Height
func MakeMultiScaleImage(width, height int) *MultiScaleImage {
	return &MultiScaleImage{
		Width:      width,
		Height:     height,
		scales:     DefaultConfig.Scales,
		grid:       util.Make2DGridFloat64Randomised(width, height),
		activators: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
		inhibitors: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
		variations: util.Make3DGridFloat64(width, height, len(DefaultConfig.Scales)),
	}
}

// MakeMultiScaleImageFromConfig create a multi-scale turing pattern image from the given config params
func MakeMultiScaleImageFromConfig(cfg Configuration) *MultiScaleImage {
	return &MultiScaleImage{
		Width:      cfg.Width,
		Height:     cfg.Height,
		scales:     cfg.Scales,
		grid:       util.Make2DGridFloat64Randomised(cfg.Width, cfg.Height),
		activators: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
		inhibitors: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
		variations: util.Make3DGridFloat64(cfg.Width, cfg.Height, len(cfg.Scales)),
	}
}

// NextIteration generate the next variation of this image
func (img MultiScaleImage) NextIteration() {
	img.calcNextVariations()
	img.normaliseGridValues()
}

// GrayscalePixmap return a grayscale pixmap derived from the current state of grid values
func (img MultiScaleImage) GrayscalePixmap() [][]uint8 {
	pixels := util.Make2DGridUInt8(img.Width, img.Height) // pixel values derived from the grid

	// map all grid values to a grayscale value
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			pixels[x][y] = uint8(math.Trunc((img.grid[x][y] + 1) / 2 * 255))
		}
	}
	return pixels
}

func (img MultiScaleImage) calcNextVariations() {
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			for k := 0; k < len(img.scales); k++ {
				img.activators[x][y][k] = util.AverageOfPixelsWithinCircle(x, y, img.scales[k].ActivatorRadius, img.grid)
				img.activators[x][y][k] *= img.scales[k].Weight

				img.inhibitors[x][y][k] = util.AverageOfPixelsWithinCircle(x, y, img.scales[k].InhibitorRadius, img.grid)
				img.inhibitors[x][y][k] *= img.scales[k].Weight

				// the variation can be calculated as an average of values within an arbitrary radius from x,y
				// but instead we use a radius of one pixel, i.e. just the value at x,y
				// apparently a radius of one pixel produces "the sharpest, most detailed images"
				img.variations[x][y][k] = math.Abs(img.activators[x][y][k] - img.inhibitors[x][y][k])
			}

			// best variation will be the smallest
			var ( // begin with values that are arbitrary yet valid
				ndx               = 0
				bestVariation     = &img.scales[0]
				smallestVariation = img.variations[x][y][0]
			)
			for k := 0; k < len(img.scales); k++ {
				if img.variations[x][y][k] < smallestVariation {
					ndx = k
					bestVariation = &img.scales[k]
				}
			}
			if img.activators[x][y][ndx] > img.inhibitors[x][y][ndx] {
				img.grid[x][y] += bestVariation.SmallAmount
			} else {
				img.grid[x][y] -= bestVariation.SmallAmount
			}
		}
	}
}

func (img MultiScaleImage) normaliseGridValues() {
	// normalise all grid values to scale them back between -1 and +1
	// begin with the min and max values across the grid

	var ( // begin with values that are arbitrary yet valid
		smallest, largest = img.grid[0][0], img.grid[0][0]
	)
	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			smallest = math.Min(smallest, img.grid[x][y])
			largest = math.Max(largest, img.grid[x][y])
		}
	}

	for x := 0; x < img.Width; x++ {
		for y := 0; y < img.Height; y++ {
			img.grid[x][y] = (img.grid[x][y]-smallest)/(largest-smallest)*2 - 1
		}
	}
}

// ReadConfigFromJSONFile Unmarshal config from a JSON file
func ReadConfigFromJSONFile(filename string) (Configuration, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	config := Configuration{}
	if err = json.Unmarshal([]byte(file), &config); err != nil {
		log.Fatal(err)
	}

	return config, err
}

// WriteConfigToJSONFile Marshal config to a JSON file
func WriteConfigToJSONFile(cfg Configuration, filename string) error {
	configJSON, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if err := ioutil.WriteFile(filename, configJSON, 0644); err != nil {
		log.Fatal(err)
	}

	return err
}
