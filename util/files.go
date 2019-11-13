package util

import (
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
)

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
