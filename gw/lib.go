package gw

import (
	"image/color"
)

type World struct {
	field      [][]int
	snake      []snake
	freeCells  int
	lenX, lenY int
	Bar        int
}

type snake struct {
	cell  []cell
	color color.RGBA
}

type cell struct {
	x, y int
}

func (w *World) Create(x, y, cells int) {
	w.field = make([][]int, x)
	for n := range w.field {
		w.field[n] = make([]int, y)
	}
	w.lenX = x
	w.lenY = y
	w.Bar = 8
}
