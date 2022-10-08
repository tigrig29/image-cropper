package imagecropper

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func GetCoordinateExcludingTransparentArea(img image.Image) (image.Rectangle, error) {
	bounds := img.Bounds()
	coordinate := image.Rectangle{
		Min: image.Point{
			X: bounds.Max.X,
			Y: bounds.Max.Y,
		},
		Max: image.Point{
			X: bounds.Min.X,
			Y: bounds.Min.Y,
		},
	}

	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			_, _, _, a := img.At(x, y).RGBA()
			if a > 0 {
				if x < coordinate.Min.X {
					coordinate.Min.X = x
				}
				if y < coordinate.Min.Y {
					coordinate.Min.Y = y
				}
				if x > coordinate.Max.X {
					coordinate.Max.X = x
				}
				if y > coordinate.Max.Y {
					coordinate.Max.Y = y
				}
			}
		}
	}

	// Because the coordinate were acquired based on the starting point of the pixel,
	// must add +1 to the Max value
	coordinate.Max.X++
	coordinate.Max.Y++

	return coordinate, nil
}

func CropImage(img image.Image, coordinate image.Rectangle, outfile string) {
	fso, err := os.Create(outfile)
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()

	cimg := img.(SubImager).SubImage(image.Rect(coordinate.Min.X, coordinate.Min.Y, coordinate.Max.X, coordinate.Max.Y))

	png.Encode(fso, cimg)
}
