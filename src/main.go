package main

import (
	"fmt"
	"image"
	"imagecropper/imagecropper"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {
	basefilepath := "../assets/base.png"
	subpath := "../assets/sub"
	outpath := "../dist"

	// 引数受け取り

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
	coordinate, err := imagecropper.GetCoordinateExcludingTransparentArea(img)
	if err != nil {
		fmt.Println("GetCoordinateExcludingTransparentArea:", err)
		return
	}

	// 引数１の画像を α で切り取り → 引数３のフォルダに出力
	imagecropper.CropImage(img, coordinate, outpath+"/base.png")

	// 引数２のフォルダにある全画像を α で切り取り → 引数３のフォルダに出力
	files, err := ioutil.ReadDir(subpath)
	if err != nil {
		fmt.Println("readdir:", err)
		return
	}
	for _, file := range files {
		path := filepath.Join(subpath, file.Name())
		f, err := os.Open(path)
		if err != nil {
			fmt.Println("open:", err)
			return
		}
		img, _, err := image.Decode(f)
		f.Close()
		if err != nil {
			fmt.Println("decode:", err)
			return
		}
		imagecropper.CropImage(img, coordinate, outpath+"/"+file.Name())
	}
}
