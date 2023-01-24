package main

import (
	"fmt"
	"image"
	"image/png"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	basefilepath := "assets/pose/c1-default-0.png"
	workpath := "assets"
	outpath := "dist"

	err := os.RemoveAll(outpath)
	if err != nil {
		fmt.Println("remove outpath:", err)
		return
	}

	err2 := os.Mkdir(outpath, os.ModePerm)
	if err2 != nil {
		fmt.Println("mkdir outpath:", err2)
		return
	}

	hoge(basefilepath, workpath, outpath)
}

func hoge(basefilepath string, workpath string, outpath string) {
	// 引数１の画像をオープン
	f, err := os.Open(basefilepath)
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

	// 画像の非透明領域の座標を取得　……　α
	coordinate, err := GetCoordinateExcludingTransparentArea(img)
	if err != nil {
		fmt.Println("GetCoordinateExcludingTransparentArea:", err)
		return
	}

	// 上余白は身長表現に利用するので切り取り対象外 = 0 にする
	coordinate.Min.Y = 0

	// 引数２のフォルダにある全画像を α で切り取り → 引数３のフォルダに出力
	files, err := ioutil.ReadDir(workpath)
	if err != nil {
		fmt.Println("readdir:", err)
		return
	}
	CropImageFilesRecursive(files, coordinate, workpath, outpath)
}

func CropImageFilesRecursive(files []fs.FileInfo, coordinate image.Rectangle, workpath string, outpath string) {
	for _, file := range files {
		fmt.Println(file.Name())
		// 対象がディレクトリ → そのディレクトリに対して処理し直す
		if file.IsDir() {
			workpath := filepath.Join(workpath, file.Name())
			outpath := filepath.Join(outpath, file.Name())
			files, err := ioutil.ReadDir(workpath)
			if err != nil {
				fmt.Println("readdir:", err)
				return
			}
			CropImageFilesRecursive(files, coordinate, workpath, outpath)
			continue
		}
		// ファイルオープン
		path := filepath.Join(workpath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("open:", err)
			return
		}
		// img 取得
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			fmt.Println("decode:", err)
			return
		}
		// フォルダなかったら作成
		if _, err := os.Stat(outpath); os.IsNotExist(err) {
			os.MkdirAll(outpath, os.ModePerm)
		}
		// 切り取り
		outfile := filepath.Join(outpath, file.Name())
		CropImage(img, coordinate, outfile)
	}
}

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
