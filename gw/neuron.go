package gw

const (
	viewRange = 5
	viewLen   = 1 + viewRange*2
)

var (
	layer = []int{viewLen * viewLen, 20, 4}
)

func (s *snake) neuroNetCreate() {
	s.neuroNet.CreateLayer(layer)
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
					dOut = -3
					if w.field[x][y] == s.num {
						dOut = -2
					}
				}
			}
			dx := x - x0
			dy := y - y0

			n := dx*viewLen + dy
			s.neuroNet.Layers[0][n].Out = float64(dOut)

			/*
				if n == 0{
					fmt.Println()
				}
				if dx % w.lenX == 0{
					fmt.Println()
				}
				switch dOut {
				case 0: fmt.Print(". ")
				case -1: fmt.Print("# ")
				case -3: fmt.Print("e ")
				case 1: fmt.Print("* ")
				case -2: fmt.Print("o ")
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
	ans := make([]float64, 4)

	for n := 0; n < 4; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 1
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}

func (s *snake) neuroBad(w *World) {
	ans := make([]float64, 4)

	for n := 0; n < 4; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 0
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}

func (s *snake) neuroWeak(w *World) {
	ans := make([]float64, 4)

	for n := 0; n < 4; n++ {
		ans[n] = s.neuroNet.Layers[len(s.neuroNet.Layers)-1][n].Out
	}

	ans[s.way] = 0.45
	s.neuroNet.SetAnswers(ans)
	s.neuroNet.Correct()
}
