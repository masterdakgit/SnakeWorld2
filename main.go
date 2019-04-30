package main

import (
	"SnakeWorld2/gw"
	"fmt"
	"time"
)

var w gw.World

func main() {
	w.Create(80, 60, 1000, 40, 300)
	go func() {
		for {
			w.Generation()
			time.Sleep(time.Duration(w.Speed) * 1000000)
			_, d, a := w.LiveDaedAll()
			if d == a {
				fmt.Println("Все змейки умерли, поколение:", w.Gen)
				break
			}
		}
	}()
	w.ListenHTTP(8080)
}
