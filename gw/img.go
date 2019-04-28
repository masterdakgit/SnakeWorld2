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
	colorWall      = color.RGBA{170, 170, 170, 255}
	colorEat       = color.RGBA{0, 170, 0, 255}
	colorAid       = color.RGBA{0, 0, 255, 255}
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
		if !w.snake[n].dead {
			for c := range w.snake[n].cell {
				if c >= len(w.snake[n].cell) {
					break
				}
				x := w.snake[n].cell[c].x
				y := w.snake[n].cell[c].y
				setBar(x, y, w.snake[n].color, i)
			}
			setBar(w.snake[n].cell[0].x, w.snake[n].cell[0].y, colorHead, i)
		}
	}
}

func (w *World) imgChange() {
	for x := 0; x < w.lenX; x++ {
		for y := 0; y < w.lenY; y++ {
			switch w.field[x][y] {
			case 0:
				setBar(x, y, colorEmpty, w.image)
			case -1:
				setBar(x, y, colorWall, w.image)
			case 1:
				setBar(x, y, colorEat, w.image)
			case 2:
				setBar(x, y, colorAid, w.image)
			}
		}
	}

	w.setSnake(w.image)
}
