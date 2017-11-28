package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"os"
)

/*
x1y1	x2y2
	  xy
x3y3	x4y4
*/

const (
	NEAREST_NAIGHBOR = iota
	BILINEAR
	BICUBIC
)

func bicubic(src image.Image, magnification float64, x, y int) color.Color {
}

func bilinear(src image.Image, magnification float64, x, y int) color.Color {
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	originX := float64(x) / magnification
	originY := float64(y) / magnification

	x0, y0 := int(originX), int(originY)
	if x0 >= srcDx-1 {
		x0 = srcDx - 2
	}
	if y0 >= srcDy-1 {
		y0 = srcDy - 2
	}

	r := make([]uint32, 4)
	g := make([]uint32, 4)
	b := make([]uint32, 4)
	a := make([]uint32, 4)

	r[0], g[0], b[0], a[0] = src.At(x0, y0).RGBA()
	r[1], g[1], b[1], a[1] = src.At(x0+1, y0).RGBA()
	r[2], g[2], b[2], a[2] = src.At(x0, y0+1).RGBA()
	r[3], g[3], b[3], a[3] = src.At(x0+1, y0+1).RGBA()

	dx := originX - float64(x0)
	dy := originY - float64(y0)

	weight := make([]float64, 4)

	weight[0] = (1 - dx) * (1 - dy)
	weight[1] = dx * (1 - dy)
	weight[2] = (1 - dx) * dy
	weight[3] = dx * dy

	newR := uint8(0)
	for i := 0; i < 4; i++ {
		newR = newR + uint8(float64(r[i]>>8)*weight[i])
	}

	newG := uint8(0)
	for i := 0; i < 4; i++ {
		newG = newG + uint8(float64(g[i]>>8)*weight[i])
	}

	newB := uint8(0)
	for i := 0; i < 4; i++ {
		newB = newB + uint8(float64(b[i]>>8)*weight[i])
	}

	newA := uint8(0)
	for i := 0; i < 4; i++ {
		newA = newA + uint8(float64(a[i]>>8)*weight[i])
	}

	return color.RGBA{newR, newG, newB, newA}
}

func nearestNeighbor(src image.Image, magnification float64, x, y int) color.Color {
	return src.At(int(float64(x)/magnification+0.5), int(float64(y)/magnification+0.5))
}

func enlargement(src image.Image, magnification float64, algo int) image.Image {
	srcBounds := src.Bounds()
	dX := srcBounds.Dx()
	dY := srcBounds.Dy()

	newDx := int(float64(dX) * magnification)
	newDy := int(float64(dY) * magnification)

	out := image.NewNRGBA(image.Rect(0, 0, newDx, newDy))

	algorithm := nearestNeighbor
	switch algo {
	case NEAREST_NAIGHBOR:
		algorithm = nearestNeighbor
	case BILINEAR:
		algorithm = bilinear
	case BICUBIC:
		algorithm = bicubic
	}

	for x := 0; x < newDx; x++ {
		for y := 0; y < newDy; y++ {
			newColor := algorithm(src, magnification, x, y)
			out.Set(x, y, newColor)
		}
	}
	return out
}

func main() {

	file, _ := os.Open("./source.jpg")
	defer file.Close()

	srcImg, _, err := image.Decode(file)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Image File Decode Error")
		return
	}

	nn, _ := os.Create("nearest_naighbor.jpg")
	defer nn.Close()

	err = jpeg.Encode(nn, enlargement(srcImg, 2.0, NEAREST_NAIGHBOR), &jpeg.Options{100})
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Write Error")
	}

	bl, _ := os.Create("bilinear.jpg")
	defer bl.Close()

	err = jpeg.Encode(bl, enlargement(srcImg, 2.0, BILINEAR), &jpeg.Options{100})
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Write Error")
	}
}
