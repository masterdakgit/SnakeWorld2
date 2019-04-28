package gw

import (
	"image"
	"log"
	"math/rand"
	"sync"
)

var (
	dir    [4]direction
	mutex2 sync.Mutex
)

type World struct {
	field      [][]int
	snake      []snake
	lenX, lenY int
	image      *image.NRGBA
}

type direction struct {
	dx, dy int
}

func (w *World) Create(x, y, nEat, nSnake int) {
	w.field = make([][]int, x)
	for n := range w.field {
		w.field[n] = make([]int, y)
	}
	w.lenX = x
	w.lenY = y
	w.setWall()

	setDir()

	w.snake = make([]snake, 0)
	w.addEat(nEat)

	for n := 0; n < nSnake; n++ {
		w.addSnake()
	}

	w.image = image.NewNRGBA(image.Rect(0, 0, bar*w.lenX+1, bar*w.lenY+1))
	w.imgChange()
}

func (w *World) setWall() {
	for x := range w.field {
		w.field[x][0] = -1
		w.field[x][w.lenY-1] = -1
	}
	for y := range w.field[0] {
		w.field[0][y] = -1
		w.field[w.lenX-1][y] = -1
	}
}

func setDir() {
	dir[0].dx = -1
	dir[0].dy = 0

	dir[1].dx = 1
	dir[1].dy = 0

	dir[2].dx = 0
	dir[2].dy = -1

	dir[3].dx = 0
	dir[3].dy = 1

}

func (w *World) Generation() {
	mutex2.Lock()
	for n := range w.snake {
		w.snake[n].move(w)
		w.snake[n].energe--
		if w.snake[n].energe < 1 {
			w.snake[n].eatSomeself(w)
			w.snake[n].energe = energeCell
		}
	}
	mutex2.Unlock()
}

func (w *World) addEat(n int) {
	r := 0
	f := 0
	for {
		if r >= n || f > 1000 {
			break
		}
		x := 1 + rand.Intn(w.lenX-3)
		y := 1 + rand.Intn(w.lenY-2)

		if w.field[x][y] == 0 {
			w.field[x][y] = 1
			r++
		} else {
			f++
		}
	}
}

func (w *World) delEat(n int) {
	r := 0
	f := 0

	x := rand.Intn(w.lenX)
	y := rand.Intn(w.lenY)
	sy := y

	for {
		if r >= n {
			break
		}
		if f > w.lenX*w.lenY {
			log.Fatal("Невозможно удалить еду.")
		}
		if w.field[x%w.lenX][y%w.lenY] == 1 {
			w.field[x%w.lenX][y%w.lenY] = 0
			r++
			f = 0
		} else {
			x++
			y = sy + x/w.lenX
			f++
		}
	}
}
