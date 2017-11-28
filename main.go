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

var a float64 = -1.0

func calcWeightBicubic(t float64) float64 {
	if t < 0 {
		t = t * -1
	}
	if t <= 1 {
		return (a+2)*t*t*t - (a+3)*t*t + 1
	}
	if t <= 2 {
		return a*t*t*t - 5*a*t*t + 8*a*t - 4*a
	}
	return 0
}

func bicubic(src image.Image, magnification float64, x, y int) color.Color {
	srcBounds := src.Bounds()
	maxDx := float64(srcBounds.Dx() - 2)
	maxDy := float64(srcBounds.Dy() - 2)

	originX := float64(x) / magnification
	originY := float64(y) / magnification

	for originX >= maxDx {
		originX = originX - 1
	}

	for originY >= maxDy {
		originY = originY - 1
	}

	for originX < 1 {
		originX = originX + 1
	}

	for originY < 1 {
		originY = originY + 1
	}

	x0, y0 := int(originX), int(originY)

	r := make([]uint32, 16)
	g := make([]uint32, 16)
	b := make([]uint32, 16)
	a := make([]uint32, 16)

	r[0], g[0], b[0], a[0] = src.At(x0-1, y0-1).RGBA()
	r[1], g[1], b[1], a[1] = src.At(x0, y0-1).RGBA()
	r[2], g[2], b[2], a[2] = src.At(x0+1, y0-1).RGBA()
	r[3], g[3], b[3], a[3] = src.At(x0+2, y0-1).RGBA()

	r[4], g[4], b[4], a[4] = src.At(x0-1, y0).RGBA()
	r[5], g[5], b[5], a[5] = src.At(x0, y0).RGBA()
	r[6], g[6], b[6], a[6] = src.At(x0+1, y0).RGBA()
	r[7], g[7], b[7], a[7] = src.At(x0+2, y0).RGBA()

	r[8], g[8], b[8], a[8] = src.At(x0-1, y0+1).RGBA()
	r[9], g[9], b[9], a[9] = src.At(x0, y0+1).RGBA()
	r[10], g[10], b[10], a[10] = src.At(x0+1, y0+1).RGBA()
	r[11], g[11], b[11], a[11] = src.At(x0+2, y0+1).RGBA()

	r[12], g[12], b[12], a[12] = src.At(x0-1, y0+2).RGBA()
	r[13], g[13], b[13], a[13] = src.At(x0, y0+2).RGBA()
	r[14], g[14], b[14], a[14] = src.At(x0+1, y0+2).RGBA()
	r[15], g[15], b[15], a[15] = src.At(x0+2, y0+2).RGBA()

	dx := make([]float64, 4)
	dy := make([]float64, 4)

	dx[0] = calcWeightBicubic(originX - float64(x0-1))
	dx[1] = calcWeightBicubic(originX - float64(x0))
	dx[2] = calcWeightBicubic(originX - float64(x0+1))
	dx[3] = calcWeightBicubic(originX - float64(x0+2))
	dy[0] = calcWeightBicubic(originY - float64(y0-1))
	dy[1] = calcWeightBicubic(originY - float64(y0))
	dy[2] = calcWeightBicubic(originY - float64(y0+1))
	dy[3] = calcWeightBicubic(originY - float64(y0+2))

	weight := make([]float64, 16)

	weight[0] = dx[0] * dy[0]
	weight[1] = dx[1] * dy[0]
	weight[2] = dx[2] * dy[0]
	weight[3] = dx[3] * dy[0]

	weight[4] = dx[0] * dy[1]
	weight[5] = dx[1] * dy[1]
	weight[6] = dx[2] * dy[1]
	weight[7] = dx[3] * dy[1]

	weight[8] = dx[0] * dy[2]
	weight[9] = dx[1] * dy[2]
	weight[10] = dx[2] * dy[2]
	weight[11] = dx[3] * dy[2]

	weight[12] = dx[0] * dy[3]
	weight[13] = dx[1] * dy[3]
	weight[14] = dx[2] * dy[3]
	weight[15] = dx[3] * dy[3]

	newColor := make([]float64, 4)
	for i := 0; i < 16; i++ {
		newColor[0] = newColor[0] + float64(r[i]>>8)*weight[i]
		newColor[1] = newColor[1] + float64(g[i]>>8)*weight[i]
		newColor[2] = newColor[2] + float64(b[i]>>8)*weight[i]
		newColor[3] = newColor[3] + float64(a[i]>>8)*weight[i]
		//fmt.Printf("%d, %.30f\n", i, weight[i])
	}

	for i := 0; i < 4; i++ {
		if newColor[i] < 0 {
			newColor[i] = 0
		} else if newColor[i] > 255 {
			newColor[i] = 255
		}
	}

	return color.RGBA{uint8(newColor[0]), uint8(newColor[1]), uint8(newColor[2]), uint8(newColor[3])}
}

func bilinear(src image.Image, magnification float64, x, y int) color.Color {
	srcBounds := src.Bounds()
	maxDx := float64(srcBounds.Dx() - 1)
	maxDy := float64(srcBounds.Dy() - 1)

	originX := float64(x) / magnification
	originY := float64(y) / magnification

	for originX >= maxDx {
		originX = originX - 1
	}

	for originY >= maxDy {
		originY = originY - 1
	}

	x0, y0 := int(originX), int(originY)

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
	newG := uint8(0)
	newB := uint8(0)
	newA := uint8(0)
	for i := 0; i < 4; i++ {
		newR = newR + uint8(float64(r[i]>>8)*weight[i])
		newG = newG + uint8(float64(g[i]>>8)*weight[i])
		newB = newB + uint8(float64(b[i]>>8)*weight[i])
		newA = newA + uint8(float64(a[i]>>8)*weight[i])
	}

	return color.RGBA{newR, newG, newB, newA}
}

func nearestNeighbor(src image.Image, magnification float64, x, y int) color.Color {
	srcBounds := src.Bounds()
	srcDx := srcBounds.Dx()
	srcDy := srcBounds.Dy()

	newX := int(float64(x)/magnification + 0.5)
	newY := int(float64(y)/magnification + 0.5)

	if newX > srcDx-1 {
		newX = srcDx - 1
	}

	if newY > srcDy-1 {
		newY = srcDy - 1
	}

	return src.At(newX, newY)
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
	bc, _ := os.Create("bicubic.jpg")
	defer bc.Close()

	err = jpeg.Encode(bc, enlargement(srcImg, 2.0, BICUBIC), &jpeg.Options{100})
	if err != nil {
		fmt.Fprintln(os.Stderr, "File Write Error")
	}
}
