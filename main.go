package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"log"
	"math"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"
)

// see: https://softologyblog.wordpress.com/2011/07/05/multi-scale-turing-patterns/

const (
	imageWidth  int = 600
	imageHeight int = 600
)

var grid [imageWidth][imageHeight]float64 // value calculations are held here
var pixels [imageWidth][imageHeight]uint8 // pixel values derived from the grid

type turingScale struct {
	ActivatorRadius int
	InhibitorRadius int
	SmallAmount     float64
	Weight          float64
	Symmetry        int
}

var scales = []turingScale{
	turingScale{100, 200, 0.05, 1, 3},
	turingScale{20, 40, 0.04, 1, 2},
	turingScale{10, 20, 0.03, 1, 2},
	turingScale{5, 10, 0.02, 1, 2},
	turingScale{1, 2, 0.01, 1, 2},
}

// imageWidth * imageHeight * len(scales)
var activators [imageWidth][imageHeight][5]float64
var inhibitors [imageWidth][imageHeight][5]float64
var variations [imageWidth][imageHeight][5]float64

func randFloat64(min, max float64) float64 {
	// generate a random float64 between the given min and max
	return min + rand.Float64()*(max-min)
}

func init() {
	for i := 0; i < imageHeight; i++ {
		for j := 0; j < imageWidth; j++ {
			grid[i][j] = randFloat64(-1.0, 1.0)
		}
	}
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

func averageOfPixelsWithinCircle(x, y, radius int) float64 {
	sum := 0.0
	numPixelsWithinCircle := 0

	for j := y - radius; j < y+radius; j++ {
		for i := x - radius; i < x+radius; i++ {

			// only include pixel values with the image bounds
			if (i >= 0) && (i < imageWidth) &&
				(j >= 0) && (j < imageHeight) {

				if pointIsWithinCircle(i, j, x, y, radius) {
					sum += grid[i][j]
					numPixelsWithinCircle++
				}
			}
		}
	}
	return sum / float64(numPixelsWithinCircle)
}

func calcVariations() {
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			for k := 0; k < len(scales); k++ {
				activators[x][y][k] = averageOfPixelsWithinCircle(x, y, scales[k].ActivatorRadius)
				activators[x][y][k] *= scales[k].Weight

				inhibitors[x][y][k] = averageOfPixelsWithinCircle(x, y, scales[k].InhibitorRadius)
				inhibitors[x][y][k] *= scales[k].Weight

				// the variation can be calculated as an average of values within an arbitrary radius from x,y
				// but instead we use a radius of one pixel, i.e. just the value at x,y
				// apparently a radius of one pixel produces "the sharpest, most detailed images"
				variations[x][y][k] += math.Abs(activators[x][y][k] - inhibitors[x][y][k])
			}

			// which scale has the smallest i.e. the best variation
			var (
				ndx               = 0
				bestVariation     *turingScale
				smallestVariation = 100.0 // begin with impossibly large number
			)
			for k := 0; k < len(scales); k++ {
				if variations[x][y][k] < smallestVariation {
					ndx = k
					bestVariation = &scales[k]
				}
			}
			if activators[x][y][ndx] > inhibitors[x][y][ndx] {
				grid[x][y] += bestVariation.SmallAmount
			} else {
				grid[x][y] -= bestVariation.SmallAmount
			}
		}
	}
}

func normaliseGridValues() {
	// normalise all grid values to scale them back between -1 and +1
	// begin with the min and max values across the grid

	var (
		smallest = 100.0  // begin with impossibly large number
		largest  = -100.0 // begin with impossibly small number
	)
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			smallest = math.Min(smallest, grid[x][y])
			largest = math.Max(largest, grid[x][y])
		}
	}

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			grid[x][y] = (grid[x][y]-smallest)/(largest-smallest)*2 - 1
		}
	}
}

func mapValuesToGrayscale() {
	// map all grid values to a grayscale value
	// alternately, map the value into a colour palette
	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			// TODO
			// DO NOT update
			pixels[x][y] = uint8(math.Trunc((grid[x][y] + 1) / 2 * 255))
		}
	}
}

func outputPNG(index int) {
	img := image.NewNRGBA(image.Rect(0, 0, imageWidth, imageHeight))

	for y := 0; y < imageHeight; y++ {
		for x := 0; x < imageWidth; x++ {
			img.Set(x, y, color.NRGBA{
				R: uint8(pixels[x][y]),
				G: uint8(pixels[x][y]),
				B: uint8(pixels[x][y]),
				A: 255,
			})
		}
	}

	filename := fmt.Sprintf("image_%02v.png", index)

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

func iterageImage() {
	for i := 1; i < 1000; i++ {
		fmt.Println("\n-----\niteration: ", i)
		fmt.Println("variations...")
		calcVariations()

		fmt.Println("normalising grid values...")
		normaliseGridValues()

		fmt.Println("mapping to grayscale...")
		mapValuesToGrayscale()

		fmt.Println("writing PNG...")
		outputPNG(i)
	}
}

// profiling: https://blog.golang.org/profiling-go-programs

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	var seed = time.Now().UnixNano()
	fmt.Println("using seed: ", seed)
	rand.Seed(seed)

	iterageImage()
}
