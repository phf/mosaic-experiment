package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/jpeg"
	"os"
)

type info struct {
	c color.RGBA
	i image.RGBA
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

func main() {
	fmt.Println("Mosaic experiment is experimental!")
	img, _ := loadImage("hm.jpg")
	rgba := convertImage(img)
	saveImage("saveHm.jpg", rgba)
	average := averageColor(*rgba)
	fmt.Printf("%v\n", average)
}
