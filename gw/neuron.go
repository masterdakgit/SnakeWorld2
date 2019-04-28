package gw

import (
	"fmt"
	"math"
)

const (
	viewRange = 8
	viewLen   = 1 + viewRange*2
	dirWay    = 5
)

var (
	layer = []int{viewLen * viewLen, 10, dirWay}
)

func (s *snake) neuroNetCreate() {
	s.neuroNet.CreateLayer(layer)
}

func (s *snake) relative(w *World, num int) bool {
	num -= 1000
	dR := math.Abs(float64(s.color.R - w.snake[num].color.R))
	dG := math.Abs(float64(s.color.G - w.snake[num].color.G))
	dB := math.Abs(float64(s.color.B - w.snake[num].color.B))

	if dR >= 245 {
		dR = 255 - dR
	}

	if dG >= 245 {
		dG = 255 - dG
	}

	if dB >= 245 {
		dB = 255 - dB
	}

	//fmt.Println(num, dR, dG, dB)
	if dR <= 10 && dG <= 10 && dB <= 10 {
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

	dOut := 0

	for y := y0; y <= y1; y++ {
		for x := x0; x <= x1; x++ {
			if y < 0 || y >= w.lenY || x < 0 || x >= w.lenX {
				dOut = -1
			} else {
				dOut = w.field[x][y]
				if dOut >= 1000 {
					dOut = -4
					if s.relative(w, w.field[x][y]) {
						dOut = -3
					}
					if w.field[x][y] == s.num {
						dOut = -2
					}
				}
			}
			dx := x - x0
			dy := y - y0

			n := dx*viewLen + dy
			s.neuroNet.Layers[0][n].Out = float64(dOut)

			if n == 0 {
				fmt.Println()
			}
			if dx%w.lenX == 0 {
				fmt.Println()
			}
			switch dOut {
			case 0:
				fmt.Print(". ")
			case -1:
				fmt.Print("# ")
			case -4:
				fmt.Print("e ")
			case 1:
				fmt.Print("* ")
			case -2:
				fmt.Print("o ")
			case -3:
				fmt.Print("r")
			}

		}
	}
}

func (s *snake) neuroWay(w *World) int {
	s.neuroSetIn(w)
	s.neuroNet.Calc()
	return s.neuroNet.MaxOutputNumber(0)
}

func (s *snake) neuroGood(w *World) {
	ans := make([]float64, dirWay)

	for n := 0; n < dirWay; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 1
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}

func (s *snake) neuroBad(w *World) {
	ans := make([]float64, dirWay)

	for n := 0; n < dirWay; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 0
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}

func (s *snake) neuroWeak(w *World) {
	ans := make([]float64, dirWay)

	for n := 0; n < dirWay; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 0.2
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}
