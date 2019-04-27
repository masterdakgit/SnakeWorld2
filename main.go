package main

import (
	"SnakeWorld2/gw"
)

var w gw.World

func main() {
	w.Create(80, 60, 1000)
	w.ListenHTTP(8080)

}
