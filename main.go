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

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")
var configfile = flag.String("configfile", "", "read image config from a json file")

func readFlags() {
	flag.Parse()
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
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
		fmt.Println("using config:\n", configJSON)
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

func generateImages() {
	img := setupImage()

	for i := 1; i > 0; i++ {
		fmt.Printf("iteration %3d...\n", i)
		filename := fmt.Sprintf("image_%03d.png", i)

		img.NextIteration()
		util.OutputPNG(filename, img.GrayscalePixmap())
	}
}

func main() {
	readFlags()
	if *cpuprofile != "" {
		defer pprof.StopCPUProfile()
	}

	var seed = time.Now().UnixNano()
	fmt.Println("using seed: ", seed)
	rand.Seed(seed)

	generateImages()
}
