package gw

import (
	"image"
	"image/color"
)

var (
	bar        int = 8
	colorEmpty     = color.RGBA{255, 255, 255, 255}
	colorGreed     = color.RGBA{220, 220, 220, 255}
	colorHead      = color.RGBA{0, 0, 0, 255}
	colorWall      = color.RGBA{0, 0, 0, 255}
	colorEat       = color.RGBA{0, 170, 0, 255}
)

func setBar(x, y int, c color.RGBA, i *image.NRGBA) {
	for bx := 0; bx < bar+1; bx++ {
		for by := 0; by < bar+1; by++ {
			i.Set(x*bar+bx, y*bar+by, c)
			if bx%bar == 0 || by%bar == 0 {
				i.Set(x*bar+bx, y*bar+by, colorGreed)
			}
		}
	}
}

func (w *World) setSnake(i *image.NRGBA) {
	for n := range w.snake {
		for c := range w.snake[n].cell {
			x := w.snake[n].cell[c].x
			y := w.snake[n].cell[c].y
			if c == 0 {
				setBar(x, y, colorHead, i)
			} else {
				setBar(x, y, w.snake[n].color, i)
			}

		}
	}
}

func (w *World) img() *image.NRGBA {
	i := image.NewNRGBA(image.Rect(0, 0, bar*w.lenX+1, bar*w.lenY+1))
	for x := 0; x < w.lenX; x++ {
		for y := 0; y < w.lenY; y++ {
			switch w.field[x][y] {
			case 0:
				setBar(x, y, colorEmpty, i)
			case -1:
				setBar(x, y, colorWall, i)
			case 1:
				setBar(x, y, colorEat, i)
			}
		}
	}

	w.setSnake(i)

	return i
}
