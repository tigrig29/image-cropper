package main

import (
	"fmt"
	"image"
	"image/png"
	"os"
)

type SubImager interface {
	SubImage(r image.Rectangle) image.Image
}

func main() {
	f, err := os.Open("../assets/images/c1-default-0.png")
	if err != nil {
		fmt.Println("open:", err)
		return
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		fmt.Println("decode:", err)
		return
	}

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

	// ピクセル始点基準で判別しているため、 Max（終点）は +1 する
	coordinate.Max.X++
	coordinate.Max.Y++

	fmt.Printf("%d, %d, %d, %d \n", coordinate.Min.X, coordinate.Min.Y, coordinate.Max.X, coordinate.Max.Y)
	fmt.Printf("%d, %d\n", coordinate.Dx(), coordinate.Dy())

	fso, err := os.Create("../assets/images/c1-default-0.out.png")
	if err != nil {
		fmt.Println("create:", err)
		return
	}
	defer fso.Close()

	cimg := img.(SubImager).SubImage(image.Rect(coordinate.Min.X, coordinate.Min.Y, coordinate.Max.X, coordinate.Max.Y))

	png.Encode(fso, cimg)
}
