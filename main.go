package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dhodges/turing_patterns/images"
)

// profiling: https://blog.golang.org/profiling-go-programs

var seed = time.Now().UnixNano()

// TODO add flag to set the initial random seed
// TODO add Dockerfile
// TODO check whether the output image has stabilized; i.e. shows no significant change from the previous iteration
// TODO flag to specify output directory for image files
// TODO flag to specify dumping current state (and iteration number) of grid when exiting
// TODO flag to specify initial state of grid
// TODO flag to specify max n iterations
// TODO generate animated PNGs

var profilecpu = flag.String("profilecpu", "", "write cpu profile to file")
var configfile = flag.String("configfile", "", "read image config from a json file")
var saveNth = flag.Int("saveNth", 1, "save an image file for each nth iteration (default: save every iteration")
var model = flag.String("model", "", "specify the generated color model ('gray' or 'rgb')")

func readFlags() {
	flag.Parse()
	if *profilecpu != "" {
		f, err := os.Create(*profilecpu)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
	}
}

// IterativeImage generic interface
type IterativeImage interface {
	ConfigFromFile(string)
	NextIteration()
	OutputPNG(string)
}

func setupImageDefault() IterativeImage {
	width, height := 600, 600

	fmt.Println("using config:")
	fmt.Println("Width: ", width)
	fmt.Println("Height: ", height)
	fmt.Println()

	switch *model {
	case "rgb":
		return images.MakeTSImageRGB(width, height)
	case "gray":
		return images.MakeTSImageGray(width, height)
	default:
		return images.MakeTSImageGray(width, height)
	}
}

func setupImage() IterativeImage {
	img := setupImageDefault()

	if *configfile != "" {
		img.ConfigFromFile(*configfile)
	}

	return img
}

func optionallySave(img IterativeImage, iteration int) {
	filename := fmt.Sprintf("image_%03d.png", iteration)
	if (*saveNth == 1) || (iteration%*saveNth == 0) {
		img.OutputPNG(filename)
	}
}

func generateImages() {
	img := setupImage()

	for i := 1; i > 0; i++ {
		fmt.Printf("iteration %3d...\r", i)

		img.NextIteration()
		optionallySave(img, i)
	}
}

func printInfo() {
	fmt.Println("using seed:  ", seed)
	if *saveNth > 1 {
		fmt.Printf("saving every: %d iterations\n", *saveNth)
	}
	switch *model {
	case "rgb":
		fmt.Println("generating color image")
	default:
		fmt.Println("generating grayscale image")
	}
}

func main() {
	readFlags()
	if *profilecpu != "" {
		defer pprof.StopCPUProfile()
	}

	rand.Seed(seed)

	printInfo()
	generateImages()
}
