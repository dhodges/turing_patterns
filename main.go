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

var profilecpu = flag.String("profilecpu", "", "write cpu profile to file")
var configfile = flag.String("configfile", "", "read image config from a json file")
var saveNth = flag.Int("saveNth", 1, "save an image file for each nth iteration (default: save every iteration")

// TODO add flag to set the initial random seed
// TODO add Dockerfile
// TODO extract golang interface for generic image generator
// TODO check whether the output image has stabilized; i.e. shows no significant change from the previous iteration
// TODO flag to specify colour model (gray, rgb, hsl)
// TODO flag to specify output directory for image files
// TODO flag to specify dumping current state (and iteration number) of grid when exiting
// TODO flag to specify initial state of grid
// TODO flag to specify max n iterations
// TODO generate animated PNGs

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

func setupImageDefault() *images.TSImageRGB {
	width, height := 600, 600

	fmt.Println("using config:")
	fmt.Println("Width: ", width)
	fmt.Println("Height: ", height)
	fmt.Println()

	return images.MakeTSImageRGB(width, height)
}

func setupImage() *images.TSImageRGB {
	img := setupImageDefault()

	if *configfile != "" {
		img.ConfigFromFile(*configfile)
	}

	return img
}

func optionallySaveImage(img *images.TSImageRGB, iteration int) {
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
		optionallySaveImage(img, i)
	}
}

func printInfo() {
	fmt.Println("using seed:  ", seed)
	if *saveNth > 1 {
		fmt.Printf("saving every: %d iterations\n", *saveNth)
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
