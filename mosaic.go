package main

import (
	"fmt"
	"image"
	_ "image/jpeg"
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

func main() {
	fmt.Println("Mosaic experiment is experimental!")
	img, _ := loadImage("hm.jpg")
	fmt.Println(img)
}
