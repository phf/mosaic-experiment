package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
	"path/filepath"
)

const size = 16 // size in pixels (square)

var tiles map[string]info

type info struct {
	c color.RGBA
	i *image.RGBA
}

func averageColor(image image.RGBA) (average color.RGBA) {
	var r, g, b, a int64
	for x := 0; x < image.Bounds().Dx(); x++ {
		for y := 0; y < image.Bounds().Dy(); y++ {
			color := image.RGBAAt(x, y)
			r += int64(color.R)
			g += int64(color.G)
			b += int64(color.B)
			a += int64(color.A)
		}
	}
	pixels := int64(image.Bounds().Dx() * image.Bounds().Dy())
	average.R = uint8(r / pixels)
	average.G = uint8(g / pixels)
	average.B = uint8(b / pixels)
	average.A = uint8(a / pixels)
	return
}

func loadImage(path string) (image.Image, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	img, _, err := image.Decode(file)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func saveImage(path string, img image.Image) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	err = jpeg.Encode(file, img, nil)
	if err != nil {
		return err
	}
	return nil
}

func convertImage(img image.Image) *image.RGBA {
	newimg := image.NewRGBA(img.Bounds())

	draw.Draw(newimg, newimg.Bounds(), img, image.ZP, draw.Src)

	return newimg
}

func loadTile(path string, inf os.FileInfo, err error) error {
	if !inf.Mode().IsRegular() {
		return nil
	}

	tile, err := loadImage(path)
	if err != nil {
		return err
	}

	converted := convertImage(tile)
	averageColor := averageColor(*converted)

	tiles[path] = info{averageColor, converted}

	return nil
}

func loadTiles() error {
	tiles = make(map[string]info)

	err := filepath.Walk("./tiles", loadTile)
	if err != nil {
		return errors.New("The directory can't be walked")
	}
	return nil
}

func main() {
	fmt.Println("Mosaic experiment is experimental!")
	img, _ := loadImage("hm.jpg")
	rgba := convertImage(img)
	saveImage("saveHm.jpg", rgba)

	err := loadTiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%v\n", tiles)
}
