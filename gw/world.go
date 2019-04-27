package gw

var (
	dir [4]direction
)

type World struct {
	field      [][]int
	snake      []snake
	eat        int
	lenX, lenY int
}

type direction struct {
	dx, dy int
}

func (w *World) Create(x, y, eat int) {
	w.field = make([][]int, x)
	for n := range w.field {
		w.field[n] = make([]int, y)
	}
	w.lenX = x
	w.lenY = y
	w.setWall()

	w.eat = eat
	setDir()

	w.snake = make([]snake, 0)
	w.addSnake()
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
	for n := range w.snake {
		w.snake[n].move(w)
		w.snake[n].energe--
		if w.snake[n].energe < 1 {
			w.snake[n].eatSomeself(w)
			w.snake[n].energe = energeCell
		}
	}
}
