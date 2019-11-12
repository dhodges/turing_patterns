package multiscale

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
)

type turingScale struct {
	ActivatorRadius int
	InhibitorRadius int
	SmallAmount     float64
	Weight          float64
	Symmetry        int
}

// Image a turing pattern image that iterates using multiple turing scale variations
// see: https://softologyblog.wordpress.com/2011/07/05/multi-scale-turing-patterns/
type Image struct {
	Width      int
	Height     int
	scales     []turingScale
	grid       [][]float64
	activators [][][]float64
	inhibitors [][][]float64
	variations [][][]float64
}

// MakeImage create an initialised multi-scale turing pattern image of given width and height
func MakeImage(width, height int) *Image {
	scales := []turingScale{
		turingScale{100, 200, 0.05, 1, 3},
		turingScale{20, 40, 0.04, 1, 2},
		turingScale{10, 20, 0.03, 1, 2},
		turingScale{5, 10, 0.02, 1, 2},
		turingScale{1, 2, 0.01, 1, 2},
	}
	return &Image{
		Width:      width,
		Height:     height,
		scales:     scales,
		grid:       make2DGridFloat64Randomised(width, height),
		activators: make3DGridFloat64(width, height, len(scales)),
		inhibitors: make3DGridFloat64(width, height, len(scales)),
		variations: make3DGridFloat64(width, height, len(scales)),
	}
}

// NextIteration generate the next variation of this image
func (img Image) NextIteration() {
	img.calcNextVariations()
	img.normaliseGridValues()
}

// GrayscalePixmap return a grayscale pixmap derived from the current state of grid values
func (img Image) GrayscalePixmap() [][]uint8 {
	pixels := make2DGridUInt8(img.Width, img.Height) // pixel values derived from the grid

	// map all grid values to a grayscale value
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			pixels[x][y] = uint8(math.Trunc((img.grid[x][y] + 1) / 2 * 255))
		}
	}
	return pixels
}

func make2DGridFloat64(width, height int) [][]float64 {
	grid := make([][]float64, height)
	for i := range grid {
		grid[i] = make([]float64, width)
	}
	return grid
}

func make2DGridUInt8(width, height int) [][]uint8 {
	grid := make([][]uint8, height)
	for i := range grid {
		grid[i] = make([]uint8, width)
	}
	return grid
}

func make2DGridFloat64Randomised(width, height int) [][]float64 {
	grid := make2DGridFloat64(width, height)
	for i := 0; i < height; i++ {
		for j := 0; j < width; j++ {
			grid[i][j] = randFloat64(-1.0, 1.0)
		}
	}
	return grid
}

func make3DGridFloat64(width, height, depth int) [][][]float64 {
	grid := make([][][]float64, height)
	for i := range grid {
		grid[i] = make([][]float64, width)
		for j := range grid[i] {
			grid[i][j] = make([]float64, depth)
		}
	}
	return grid
}

func randFloat64(min, max float64) float64 {
	// generate a random float64 between the given min and max
	return min + rand.Float64()*(max-min)
}

func sqr(x int) int {
	return x * x
}

func pointIsWithinCircle(xp, yp, x, y, radius int) bool {
	// does the given point exist within (or on) the given circle?
	// xp, yp: the point
	// x, y, radius: the circle

	// radius <= Math.sqrt((xp - x)² + (yp - y)²)
	// i.e.
	// radius² <= (xp - x)² + (yp - y)²

	return sqr(radius) <= sqr(xp-x)+sqr(yp-y)
}

func (img Image) averageOfPixelsWithinCircle(x, y, radius int) float64 {
	// x, y, radius: the circle of values from which to derive an average
	sum := 0.0
	numPixelsWithinCircle := 1.0

	for j := y - radius; j < y+radius; j++ {
		for i := x - radius; i < x+radius; i++ {

			// only include pixel values with the image bounds
			if (i >= 0) && (i < img.Width) &&
				(j >= 0) && (j < img.Height) {

				if pointIsWithinCircle(i, j, x, y, radius) {
					sum += img.grid[i][j]
					numPixelsWithinCircle++
				}
			}
		}
	}
	return sum / numPixelsWithinCircle
}

func (img Image) calcNextVariations() {
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			for k := 0; k < len(img.scales); k++ {
				img.activators[x][y][k] = img.averageOfPixelsWithinCircle(x, y, img.scales[k].ActivatorRadius)
				img.activators[x][y][k] *= img.scales[k].Weight

				img.inhibitors[x][y][k] = img.averageOfPixelsWithinCircle(x, y, img.scales[k].InhibitorRadius)
				img.inhibitors[x][y][k] *= img.scales[k].Weight

				// the variation can be calculated as an average of values within an arbitrary radius from x,y
				// but instead we use a radius of one pixel, i.e. just the value at x,y
				// apparently a radius of one pixel produces "the sharpest, most detailed images"
				img.variations[x][y][k] += math.Abs(img.activators[x][y][k] - img.inhibitors[x][y][k])
			}

			// which scale has the smallest i.e. the best variation
			var (
				ndx               = 0
				bestVariation     *turingScale
				smallestVariation = 100.0 // begin with impossibly large number
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

func (img Image) normaliseGridValues() {
	// normalise all grid values to scale them back between -1 and +1
	// begin with the min and max values across the grid

	var (
		smallest = 100.0  // begin with impossibly large number
		largest  = -100.0 // begin with impossibly small number
	)
	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			smallest = math.Min(smallest, img.grid[x][y])
			largest = math.Max(largest, img.grid[x][y])
		}
	}

	for y := 0; y < img.Height; y++ {
		for x := 0; x < img.Width; x++ {
			img.grid[x][y] = (img.grid[x][y]-smallest)/(largest-smallest)*2 - 1
		}
	}
}

// OutputPNG export this image as a PNG
func OutputPNG(filename string, pixmap [][]uint8) {
	height := len(pixmap)
	width := len(pixmap[0])
	img := image.NewNRGBA(image.Rect(0, 0, width, height))

	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(pixmap[x][y]),
				G: uint8(pixmap[x][y]),
				B: uint8(pixmap[x][y]),
				A: 255,
			})
		}
	}

	f, err := os.Create(filename)
	if err != nil {
		log.Fatal(err)
	}

	if err := png.Encode(f, img); err != nil {
		f.Close()
		log.Fatal(err)
	}

	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}
