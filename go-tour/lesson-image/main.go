package main

import (
	"fmt"
	"golang.org/x/tour/pic"
	"image"
	"image/color"
)

type Image struct {
	w, h int
}

func (i Image) Bounds() image.Rectangle {
	return image.Rect(0, 0, i.w, i.h)
}

func (i Image) ColorModel() color.Model {
	return color.RGBAModel
}

func (i Image) At(x, y int) color.Color {
	return color.RGBA{uint8(x), uint8(y), uint8((x + y) / 2), 255}
}

// 展示图片
func (i Image) Show() {
	bounds := i.Bounds()
	for x := bounds.Min.X; x <= bounds.Max.X; x++ {
		for y := bounds.Min.Y; y <= bounds.Max.Y; y++ {
			fmt.Printf("%.2f ", Color2Num(i.At(x, y)))
		}
		fmt.Println()
	}
}

// 将颜色转为数值
func Color2Num(c color.Color) float64 {
	r, g, b, a := c.RGBA()
	return float64(r+g+b+a) / 4.0
}

func main() {
	m := Image{100, 100}
	m.Show()
	pic.ShowImage(m)
}
