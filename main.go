package main

import (
	"SnakeWorld2/gw"
	"fmt"
	"time"
)

var w gw.World

func main() {
	w.Create(80, 60, 500, 1)
	go func() {
		for {
			w.Generation()
			time.Sleep(10000 * 1000000)
			fmt.Println(w.LiveDaedAll())
		}
	}()
	w.ListenHTTP(8080)
}
