package gw

import (
	"image/color"
	"log"
	"math/rand"
	"neuron/nr"
)

var (
	startLength = 4
	energeCell  = 10
)

type snake struct {
	cell     []cell
	color    color.RGBA
	num      int
	way      int
	energe   int
	neuroNet nr.NeuroNet
	memory   memory
	dead     bool
	genOld   int
	Age      int
	nCorrect float64
	diver    int
	test     int
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
	s.Age++
	s.way = s.neuroWay(w)

	if !s.dead {
		s.step(w)
	}
}

func (s *snake) step(w *World) {
	x := (w.lenX + s.cell[0].x + dir[s.way].dx) % w.lenX
	y := (w.lenY + s.cell[0].y + dir[s.way].dy) % w.lenY

	if w.field[x][y] == -1 {
		s.neuroCorrect(w, 0.05)
		return
	}

	if w.field[x][y] >= 1000 {
		s.neuroCorrect(w, 0.05)
		return
	}

	if w.field[x][y] > 1 {
		s.neuroCorrect(w, 0.05)
		return
	}

	if w.field[x][y] == 1 {
		s.eat(w)
		s.neuroCorrect(w, 0.95)
	} else {
		s.neuroCorrect(w, 0.5)
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
	s.Age = 0
	s.dead = true
}

func (s *snake) eat(w *World) {
	var c cell
	nLast := len(s.cell) - 1
	c.x = s.cell[nLast].x
	c.y = s.cell[nLast].y
	s.cell = append(s.cell, c)

	if len(s.cell) >= s.diver {
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

	if L < s.diver {
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
	s.energe /= 2
	newSnake.energe = s.energe

	newSnake.nCorrect = s.nCorrect
	newSnake.neuroNet = s.neuroNet
	newSnake.diver = s.diver
	newSnake.memory = s.memory
	newSnake.test = s.test

	if L != len(newSnake.cell)+len(s.cell) {
		log.Fatal("Div: Ошибка деления. Не правильно расчитана длинна.")
	}

	for n := range w.snake {
		if w.snake[n].dead {
			newSnake.num = w.snake[n].num
			w.snake[n] = newSnake
			w.snake[n].dead = false
			return
		}
	}

	newSnake.num = len(w.snake) + 1000
	w.snake = append(w.snake, newSnake)
}
