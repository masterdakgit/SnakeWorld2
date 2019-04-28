package gw

import (
	"image/color"
	"log"
	"math/rand"
	"neuron/nr"
)

var (
	startLength = 1
	energeCell  = 20
)

type snake struct {
	cell     []cell
	color    color.RGBA
	num      int
	way      int
	energe   int
	dead     bool
	neuroNet nr.NeuroNet
}

type cell struct {
	x, y int
}

func (w *World) findEmptyXY() (x, y int) {
	x = 1 + rand.Intn(w.lenX-3)
	y = 1 + rand.Intn(w.lenY-2)
	r := 0

	for {
		if r > 100 {
			log.Fatal("addSnake: Нет места для новой змейки.")
		}
		if w.field[x][y] == 0 {
			break
		}
		x = (x + 1) % w.lenX
		y = (y + 1) % w.lenY
		r++
	}
	return
}

func (w *World) addSnake() {
	var s snake
	s.cell = make([]cell, startLength)
	s.energe = energeCell

	x, y := w.findEmptyXY()

	for n := range s.cell {
		s.cell[n].x = x
		s.cell[n].y = y
	}

	num := len(w.snake) + 1000
	s.num = num
	w.field[x][y] = num

	R := uint8(rand.Intn(255))
	G := uint8(rand.Intn(255))
	B := uint8(rand.Intn(255))
	s.color = color.RGBA{R, G, B, 255}

	s.neuroNetCreate()

	w.snake = append(w.snake, s)
}

func (s *snake) move(w *World) {
	s.randomWay(w)
	if !s.dead {
		s.step(w)
	}
}

func (s *snake) randomWay(w *World) {
	d := 0
	way := rand.Intn(4)

	for {
		if d > 3 {
			s.die(w)
			break
		}

		x := s.cell[0].x + dir[way].dx
		y := s.cell[0].y + dir[way].dy
		if w.field[x][y] >= 0 && w.field[x][y] < 1000 {
			s.way = way
			break
		}
		way = (way + 1) % 4
		d++
	}
}

func (s *snake) step(w *World) {
	x := s.cell[0].x + dir[s.way].dx
	y := s.cell[0].y + dir[s.way].dy

	if w.field[x][y] == 1 {
		s.eat(w)
	}

	nLast := len(s.cell) - 1
	w.field[s.cell[nLast].x][s.cell[nLast].y] = 0

	for n := nLast; n > 0; n-- {
		s.cell[n] = s.cell[n-1]
	}

	s.cell[0].x = x
	s.cell[0].y = y
	w.field[x][y] = s.num
}

func (s *snake) die(w *World) {
	for n := range s.cell {
		w.field[s.cell[n].x][s.cell[n].y] = 1
	}

	w.delEat(startLength - len(s.cell))
	s.cell = s.cell[:startLength]

	x, y := w.findEmptyXY()

	for n := range s.cell {
		s.cell[n].x = x
		s.cell[n].y = y
	}
}

func (s *snake) eat(w *World) {
	var c cell
	nLast := len(s.cell) - 1
	c.x = s.cell[nLast].x
	c.y = s.cell[nLast].y
	s.cell = append(s.cell, c)
}

func (s *snake) eatSomeself(w *World) {
	nLast := len(s.cell) - 1
	w.addEat(1)

	if nLast < 1 {
		s.die(w)
		return
	}

	w.field[s.cell[nLast].x][s.cell[nLast].y] = 0
	s.cell = s.cell[:nLast]
}
