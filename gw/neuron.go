package gw

const (
	viewRange = 10
	viewLen   = 1 + viewRange*2
)

var (
	layer = []int{viewLen * viewLen, 4}
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
					dOut = -1
				}
			}
			dx := x - x0
			dy := y - y0
			n := dx*viewLen + dy
			s.neuroNet.Layers[0][n].Out = float64(dOut)
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
