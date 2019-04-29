package gw

import (
	"image/color"
	"log"
	"math/rand"
	"neuron/nr"
)

var (
	startLength = 4
	energeCell  = 20
	diver       = 8
)

type snake struct {
	cell         []cell
	color        color.RGBA
	num          int
	way          int
	energe       int
	neuroNet     nr.NeuroNet
	neuroLayer   []int
	dead         bool
	humanControl bool
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
	//s.randomWay(w)
	s.way = s.neuroWay(w)
	if s.way == 4 {
		s.div(w)
		return
	}

	if !s.dead {
		s.step(w)
	}
}

func (s *snake) randomWay(w *World) {
	d := 0
	way := rand.Intn(4)

	for {
		if d > 3 {
			//s.die(w)
			break
		}

		x := s.cell[0].x + dir[way].dx
		y := s.cell[0].y + dir[way].dy

		if x < 0 || y < 0 || x == w.lenX || y == w.lenY {
			log.Fatal("randomWay: ", x, y)
		}

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

	if w.field[x][y] == -1 || w.field[x][y] >= 1000 {
		if rand.Intn(100) == 0 {
			s.die(w)
			return
		}
		s.neuroBad(w)
		return
	}

	if w.field[x][y] == 1 {
		s.eat(w)
		s.neuroGood(w)
	} else {
		s.neuroWeak(w)
	}

	nLast := len(s.cell) - 1
	w.field[s.cell[nLast].x][s.cell[nLast].y] = 0

	for n := nLast; n > 0; n-- {
		s.cell[n] = s.cell[n-1]
	}

	s.cell[0].x = x
	s.cell[0].y = y
	w.field[x][y] = s.num

	w.field[s.cell[nLast].x][s.cell[nLast].y] = s.num
}

func (s *snake) die(w *World) {
	for n := range s.cell {
		w.field[s.cell[n].x][s.cell[n].y] = 1
	}

	s.dead = true
}

func (s *snake) eat(w *World) {
	var c cell
	nLast := len(s.cell) - 1
	c.x = s.cell[nLast].x
	c.y = s.cell[nLast].y
	s.cell = append(s.cell, c)

	if len(s.cell) >= diver {
		s.div(w)
	}
}

func (s *snake) eatSomeself(w *World) {
	nLast := len(s.cell) - 1

	if nLast < 1 {
		s.die(w)
		return
	}

	w.field[s.cell[nLast].x][s.cell[nLast].y] = 0
	s.cell = s.cell[:nLast]

	s.energe = energeCell
}

func (s *snake) div(w *World) {
	L := len(s.cell)

	if L < diver {
		s.neuroBad(w)
		return
	}

	var newSnake snake

	R, G, B, A := s.color.RGBA()
	Rr := uint8(R-10) + uint8(rand.Intn(20))
	Gr := uint8(G-10) + uint8(rand.Intn(20))
	Br := uint8(B-10) + uint8(rand.Intn(20))

	newSnake.color = color.RGBA{Rr, Gr, Br, uint8(A)}
	newSnake.cell = make([]cell, len(s.cell)/2)

	for n := range newSnake.cell {
		newSnake.cell[n].x = s.cell[len(s.cell)/2+n].x
		newSnake.cell[n].y = s.cell[len(s.cell)/2+n].y
	}

	s.cell = s.cell[:L-len(s.cell)/2]
	newSnake.num = len(w.snake) + 1000
	newSnake.energe = energeCell
	newSnake.neuroNet = s.neuroNet
	newSnake.neuroLayer = s.neuroLayer

	if L != len(newSnake.cell)+len(s.cell) {
		log.Fatal("Div: Ошибка деления. Не правильно расчитана длинна.")
	}

	for n := range w.snake {
		if w.snake[n].dead {
			w.snake[n] = newSnake
			w.snake[n].dead = false
			return
		}
	}

	w.snake = append(w.snake, newSnake)
}
