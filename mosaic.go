package main

import (
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"math"
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

	start := image.PixOffset(image.Bounds().Min.X, image.Bounds().Min.Y)
	index := start
	for y := 0; y < image.Bounds().Dy(); y++ {
		for x := 0; x < image.Bounds().Dx(); x++ {
			pixel := image.Pix[index : index+4]
			r += int64(pixel[0])
			g += int64(pixel[1])
			b += int64(pixel[2])
			a += int64(pixel[3])
			index += 4
		}
		index += image.Stride - image.Bounds().Dx()*4
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

func min(a, b int) int {
	return int(math.Min(float64(a), float64(b)))
}

func sqr(x float64) float64 {
	return x * x
}

func calculateDistance(a, b color.RGBA) float64 {
	R := sqr(float64(a.R) - float64(b.R))
	G := sqr(float64(a.G) - float64(b.G))
	B := sqr(float64(a.B) - float64(b.B))
	A := sqr(float64(a.A) - float64(b.A))
	return math.Sqrt(R + G + B + A)
}

func createMosaic(org *image.RGBA) *image.RGBA {

	img := image.NewRGBA(org.Bounds())

	max_x := org.Bounds().Max.X - 1
	max_y := org.Bounds().Max.Y - 1

	for x := 0; x <= max_x; x += size {
		for y := 0; y <= max_y; y += size {
			rect := image.Rect(x, y, min(x+size, max_x), min(y+size, max_y))
			tile := (org.SubImage(rect)).(*image.RGBA)
			avg := averageColor(*tile)

			minimal := 1e9
			distance := 0.0
			closest := ""

			for name, info := range tiles {
				distance = calculateDistance(avg, info.c)
				if distance < minimal {
					minimal = distance
					closest = name
				}
			}

			draw.Draw(img, rect, tiles[closest].i, image.ZP, draw.Src)
			//return
			//draw.DrawMask(org, rect, tiles[closest].i, rect.Bounds.Min, nil,  image.ZP, draw.Src)
		}
	}
	return img
}

func main() {
	fmt.Println("Mosaic experiment is experimental!")
	img, _ := loadImage("hm.jpg")
	rgba := convertImage(img)

	err := loadTiles()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	saveImage("saveHm.jpg", createMosaic(rgba))
}
