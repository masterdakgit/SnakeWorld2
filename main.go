package main

import (
	"SnakeWorld2/gw"
	"time"
)

var w gw.World

func main() {
	w.Create(80, 60, 1000, 10)
	go func() {
		for {
			w.Generation()
			time.Sleep(1 * 1000000)
		}
	}()
	w.ListenHTTP(8080)
}
