package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	"github.com/dhodges/turing_patterns/images"
	"github.com/dhodges/turing_patterns/util"
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

func setupImageConfigured(configfile string) *images.MultiScaleImage {
	var img *images.MultiScaleImage

	if cfg, err := images.ReadConfigFromJSONFile(configfile); err != nil {
		log.Fatal(err)
	} else {
		configJSON, err := json.MarshalIndent(cfg, "", "  ")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("using config:\n", string(configJSON))
		fmt.Println()
		img = images.MakeMultiScaleImageFromConfig(cfg)
	}
	return img
}

func setupImageDefault() *images.MultiScaleImage {
	width, height := 600, 600
	img := images.MakeMultiScaleImage(width, height)

	fmt.Println("using config:")
	fmt.Println("Width: ", width)
	fmt.Println("Height: ", height)
	fmt.Println("Scales: ", images.DefaultConfig.Scales)
	fmt.Println()

	return img
}

func setupImage() *images.MultiScaleImage {
	var img *images.MultiScaleImage

	if *configfile != "" {
		img = setupImageConfigured(*configfile)
	} else {
		img = setupImageDefault()
	}
	return img
}

func optionallySaveImage(img *images.MultiScaleImage, iteration int) {
	filename := fmt.Sprintf("image_%03d.png", iteration)
	if (*saveNth == 1) || (iteration%*saveNth == 0) {
		util.OutputPNG(filename, img.GrayscalePixmap())
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
