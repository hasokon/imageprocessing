package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"os"
)

func main() {

	file, _ := os.Open("./source.jpg")
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Image File Decode Error")
		return
	}

	srcBounds := srcImg.Bounds()
	maxX := srcBounds.Dx()
	maxY := srcBounds.Dy()

	outImg := image.NewNRGBA(image.Rect(0, 0, maxX, maxY))

	for x := 0; x < maxX; x++ {
		for y := 0; y < maxY; y++ {
			outImg.Set(x, y, srcImg.At(x, y))
		}
	}

	newfile, _ := os.Create("output.jpg")
	defer newfile.Close()

	err = jpeg.Encode(newfile, outImg, &jpeg.Options{100})
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Write Error")
	}
}
