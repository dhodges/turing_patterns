package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"runtime/pprof"
	"time"

	multiscale "github.com/dhodges/turing_patterns/multi_scale"
)

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

	img := multiscale.MakeImage(600, 600)

	for i := 1; i > 0; i++ {
		fmt.Printf("iteration %2d...\n", i)
		filename := fmt.Sprintf("image_%02d.png", i)
		img.NextIteration()
		multiscale.OutputPNG(filename, img.GrayscalePixmap())
	}
}
