package gw

import (
	"log"
	"math"
	"math/rand"
	"neuron/nr"
)

const (
	viewRange = 4
	viewLen   = 1 + viewRange*2
	lenMomory = viewRange * 2
	dirWay    = 4
)

type memory struct {
	data [lenMomory][]nr.Neuron
	way  [lenMomory]int
	pos  int
}

func (s *snake) neuroNetCreate() {
	s.diver = 8
	s.nCorrect = 0.2
	s.test = 1 + rand.Intn(2)

	neuroLayer := make([]int, 2)
	neuroLayer[0] = viewLen * viewLen
	neuroLayer[1] = dirWay

	for n := range s.memory.data {
		s.memory.data[n] = make([]nr.Neuron, viewLen*viewLen)
	}

	s.neuroNet.CreateLayer(neuroLayer)
}

func (s *snake) relative(w *World, num int) bool {
	num -= 1000

	if num >= len(w.snake) {
		return false
	}

	dR := math.Abs(float64(s.color.R - w.snake[num].color.R))
	dG := math.Abs(float64(s.color.G - w.snake[num].color.G))
	dB := math.Abs(float64(s.color.B - w.snake[num].color.B))

	if dR >= 235 {
		dR = 255 - dR
	}

	if dG >= 235 {
		dG = 255 - dG
	}

	if dB >= 235 {
		dB = 255 - dB
	}

	if dR <= 20 && dG <= 20 && dB <= 20 {
		return true
	}

	return false
}

func (s *snake) neuroSetIn(w *World) {
	x := s.cell[0].x
	y := s.cell[0].y
	x0 := x - viewRange
	x1 := x + viewRange
	y0 := y - viewRange
	y1 := y + viewRange

	dOut := float64(0.01)
	//str := "  "

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if y < 0 || y >= w.lenY || x < 0 || x >= w.lenX {
				dOut = 0.01 //Выход за край карты
				//str = "##"
			} else {
				dOut, _ = s.dataToOut(w, w.field[x][y])
			}

			dx := x - x0
			dy := y - y0

			n := dx*viewLen + dy
			s.neuroNet.Layers[0][n].Out = dOut
			/*
				if n == 0 {
					fmt.Println()
				}
				if dx%w.lenX == 0 {
					fmt.Println()
				}
				fmt.Print(str)
			*/
		}
	}

	s.memory.data[s.memory.pos] = s.neuroNet.Layers[0]

}

func (s *snake) neuroWay(w *World) int {
	s.neuroSetIn(w)
	s.neuroNet.Calc()
	mo := s.neuroNet.MaxOutputNumber(0)
	s.memory.way[s.memory.pos] = mo
	s.memory.pos = (s.memory.pos + 1) % lenMomory
	return mo
}

func (s *snake) neuroCorrect(w *World, a float64) {
	ans := make([]float64, dirWay)
	n := float64(lenMomory)
	way := 0

	switch s.test {
	case 0:
		log.Fatal("s.test: Не установлено значение.")
	case 1:
		for pos := s.memory.pos + lenMomory; pos > s.memory.pos; pos-- {
			p := pos % lenMomory

			s.neuroNet.NCorrect = 0.1 + 0.4*n/lenMomory
			n--
			s.neuroNet.Layers[0] = s.memory.data[p]
			s.neuroNet.Calc()

			for n := 0; n < dirWay; n++ {
				ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
			}

			way = s.neuroNet.MaxOutputNumber(0)

			ans[way] = a

			s.neuroNet.SetAnswers(ans)
			s.neuroNet.Correct()
		}
	case 2:
		for pos := s.memory.pos + lenMomory; pos > s.memory.pos; pos-- {
			p := pos % lenMomory

			s.neuroNet.NCorrect = 0.1 + 0.4*n/lenMomory
			n--
			s.neuroNet.Layers[0] = s.memory.data[p]
			s.neuroNet.Calc()

			for n := 0; n < dirWay; n++ {
				ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
			}

			way = s.neuroNet.MaxOutputNumber(0)

			ans[way] = a

			s.neuroNet.SetAnswers(ans)
			s.neuroNet.Correct()
		}
	}
}

func (w *World) bestTest() (t1, t2 int) {
	for n := range w.snake {
		if !w.snake[n].dead {
			switch w.snake[n].test {
			case 1:
				t1++
			case 2:
				t2++
			}
		}
	}
	return
}

func (s *snake) dataToOut(w *World, data int) (d float64, str string) {
	switch data {
	case -1:
		return -0.5, "# " //Стена
	case 0:
		return 0.01, ". " //Пусто
	case 1:
		return 0.99, "* " //Еда
	}

	if data > 1 && data < 1000 {
		return -0.9, "A " //Яд
	}

	if data >= 1000 {
		if s.relative(w, data) {
			return -0.1, "o " //Я или родственник
		} else {
			return -0.3, "e " //Чужой
		}
	}

	log.Fatal("dataToOut: Пустое значение.")
	return 0, "  " //
}
