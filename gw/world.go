package gw

import (
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"sync"
)

var (
	dir    [dirWay]direction
	mutex2 sync.Mutex
	wg     sync.WaitGroup
	Core   = 1
	rAid   = 2
)

type World struct {
	field      [][]int
	snake      []snake
	lenX, lenY int
	balance    int
	image      *image.NRGBA
	Gen        int
	ageEra     int
	Speed      float64
	minSnake   int
}

type direction struct {
	dx, dy int
}

func (w *World) Create(x, y, nEat, minSnake, rWall int) {
	w.field = make([][]int, x)
	for n := range w.field {
		w.field[n] = make([]int, y)
	}
	w.lenX = x
	w.lenY = y
	w.setWall()
	w.addRandomWall(rWall)
	w.Speed = 100

	setDir()

	w.balance = nEat

	w.snake = make([]snake, 0)
	w.addEat(nEat)
	w.minSnake = minSnake

	for n := 0; n < minSnake; n++ {
		w.addSnake()
	}

	w.image = image.NewNRGBA(image.Rect(0, 0, bar*w.lenX+1, bar*w.lenY+1+infoPanelY))
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

func (w *World) addRandomWall(n int) {
	r := 0
	f := 0
	for {
		if r >= n || f > 1000 {
			break
		}
		x := 1 + rand.Intn(w.lenX-3)
		y := 1 + rand.Intn(w.lenY-2)

		if w.field[x][y] == 0 {
			w.field[x][y] = -1
			r++
		} else {
			f++
		}
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

	for n := 0; n < len(w.snake)/Core; n++ {
		wg.Add(Core)
		for cr := 0; cr < Core; cr++ {
			nc := n*Core + cr
			go func() {
				if !w.snake[nc].dead {
					w.snake[nc].move(w)
					w.snake[nc].energe--
					if w.snake[nc].energe < 1 {
						w.snake[nc].eatSomeself(w)
					}
				}
				wg.Done()
			}()
		}
		wg.Wait()
	}

	ns := len(w.snake) / Core
	ns *= Core

	for nc := ns; nc < len(w.snake); nc++ {
		if !w.snake[nc].dead {
			w.snake[nc].move(w)
			w.snake[nc].energe--
			if w.snake[nc].energe < 1 {
				w.snake[nc].eatSomeself(w)
			}
		}
	}

	l, _ := w.liveDeadSnakes()
	if l < w.minSnake {
		fmt.Println("Добавляем новую змейку, поколение:", w.Gen)
		w.ageEra = w.Gen
		for n := range w.snake {
			if w.snake[n].dead {
				var s snake
				s.cell = make([]cell, startLength)
				s.energe = energeCell

				x, y := w.findEmptyXY()

				for n := range s.cell {
					s.cell[n].x = x
					s.cell[n].y = y
				}

				s.num = w.snake[n].num
				w.field[x][y] = s.num

				R := uint8(rand.Intn(255))
				G := uint8(rand.Intn(255))
				B := uint8(rand.Intn(255))
				s.color = color.RGBA{R, G, B, 255}

				s.neuroNetCreate()
				w.snake[n] = s
				break
			}
		}
	}

	w.setBalanceEat()
	w.Gen++
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
			fmt.Println("Еды нет.")
			break
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

func (w *World) removeAid(x, y int) {
	w.field[x][y]--
}

func (w *World) calcCell() int {
	result := 0
	for x := range w.field {
		for y := range w.field[0] {
			if w.field[x][y] >= 1 {
				result++
			}
			if w.field[x][y] > 1 && w.field[x][y] <= 1+rAid*energeCell {
				w.removeAid(x, y)
			}
		}
	}
	for n := range w.snake {
		if !w.snake[n].dead {
			result += len(w.snake[n].cell)
		}
	}
	return result
}

func (w *World) setBalanceEat() {
	n := w.balance - w.calcCell()

	if n > 0 {
		w.addEat(n)
	} else {
		w.delEat(w.calcCell() - w.balance)
	}
}

func (w *World) liveDeadSnakes() (l, d int) {
	l = 0
	d = 0
	for n := range w.snake {
		if !w.snake[n].dead {
			l++
		} else {
			d++
		}
	}
	return
}

func (w *World) LiveDaedAll() (l, d, a int) {
	l, d = w.liveDeadSnakes()
	a = len(w.snake)
	return
}

func (w *World) avergeAge() int {
	sumAge := 0
	nSnake := 0
	for n := range w.snake {
		if w.snake[n].dead {
			continue
		}
		nSnake++
		sumAge += w.snake[n].Age
	}
	return sumAge / nSnake
}
