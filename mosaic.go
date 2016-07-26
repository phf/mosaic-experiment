package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"os"
)

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

	saveImage("saveHm.jpg", img)
	convertImage(img)

}