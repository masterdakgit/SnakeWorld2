package gw

import (
	"image/color"
	"math"
	"math/rand"
	"strconv"
)

const (
	viewRange = 4
	viewLen   = 1 + viewRange*2
	dirWay    = 4
)

func (s *snake) neuroNetCreate() {
	hidenLayer := rand.Intn(3)
	s.neuroLayer = make([]int, 2+hidenLayer)

	s.neuroLayer[0] = viewLen * viewLen
	s.neuroLayer[len(s.neuroLayer)-1] = dirWay

	if hidenLayer > 0 {
		s.neuroLayer[len(s.neuroLayer)-2] = dirWay + rand.Intn(viewLen*dirWay)
	}

	if hidenLayer > 1 {
		s.neuroLayer[len(s.neuroLayer)-3] = viewLen*viewLen + rand.Intn(viewLen*viewLen*dirWay)
	}

	//fmt.Println(s.neuroLayer)
	s.neuroNet.CreateLayer(s.neuroLayer)
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

	//fmt.Println(num, dR, dG, dB)
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
				if dOut > 1 && dOut <= 1+rAid*energeCell {
					dOut = -5
				}
			}
			dx := x - x0
			dy := y - y0

			n := dx*viewLen + dy
			s.neuroNet.Layers[0][n].Out = float64(dOut)
			/*
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
					fmt.Print("r ")
				case -5:
					fmt.Print("A ")
				}
			*/
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

	ans[s.way] = 0.45
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}

func (w *World) bestNeuroLayer() (bestLayerStr string, color color.RGBA, age int) {
	liveLayer := make(map[string]int)
	bestLayer := 0
	bestLayerStr = ""
	for n := range w.snake {
		if w.snake[n].dead {
			continue
		}
		str := ""
		for s := range w.snake[n].neuroLayer {
			str += strconv.Itoa(w.snake[n].neuroLayer[s]) + " "
		}
		liveLayer[str]++
		if bestLayer < liveLayer[str] {
			bestLayer = liveLayer[str]
			bestLayerStr = str
			color = w.snake[n].color
			age = n
		}
	}
	return
}
