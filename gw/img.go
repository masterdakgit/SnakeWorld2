package gw

import (
	"golang.org/x/image/font"
	"golang.org/x/image/font/basicfont"
	"golang.org/x/image/math/fixed"
	"image"
	"image/color"
	"strconv"
)

const (
	infoPanelY = 64
)

var (
	bar        int = 8
	colorEmpty     = color.RGBA{255, 255, 255, 255}
	colorGreed     = color.RGBA{220, 220, 220, 255}
	colorHead      = color.RGBA{0, 0, 0, 255}
	colorWall      = color.RGBA{170, 170, 170, 255}
	colorEat       = color.RGBA{0, 170, 0, 255}
	colorAid       = color.RGBA{255, 0, 0, 255}
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
			default:
				if w.field[x][y] <= 1+rAid*energeCell {
					setBar(x, y, colorAid, w.image)
				}
			}
		}
	}

	w.setSnake(w.image)

	w.infoPanelClear()
	addLabel(w.image, 10, bar*w.lenY+20, "Pause: "+strconv.Itoa(int(w.Speed))+"ms")
	addLabel(w.image, bar*w.lenX/2, bar*w.lenY+20, "Generation: "+strconv.Itoa(int(w.Gen)))
	addLabel(w.image, bar*w.lenX/2, bar*w.lenY+40, "Balance: "+strconv.Itoa(int(w.balance)))
	s, c, _ := w.bestNeuroLayer()
	addLabel(w.image, 22, bar*w.lenY+40, "Best neuron layer: "+s)
	addLabel(w.image, bar*w.lenX/2, bar*w.lenY+60, "Averge age: "+strconv.Itoa(w.avergeAge()))
	w.bestColorToInfoPanel(c)
}

func (w *World) bestColorToInfoPanel(c color.RGBA) {
	for x := 10; x < 10+bar; x++ {
		for y := bar*w.lenY + 31; y < bar*w.lenY+31+bar; y++ {
			w.image.Set(x, y, c)
		}
	}
}

func (w *World) infoPanelClear() {
	for x := 0; x < bar*w.lenX+1; x++ {
		for y := bar*w.lenY + 1; y < bar*w.lenY+1+infoPanelY; y++ {
			w.image.Set(x, y, color.RGBA{255, 255, 255, 255})
		}
	}
}

func addLabel(img *image.NRGBA, x, y int, label string) {
	col := color.RGBA{0, 0, 0, 255}
	point := fixed.Point26_6{fixed.Int26_6(x * 64), fixed.Int26_6(y * 64)}

	d := &font.Drawer{
		Dst:  img,
		Src:  image.NewUniform(col),
		Face: basicfont.Face7x13,
		Dot:  point,
	}
	d.DrawString(label)
}
