package gw

import (
	"fmt"
	"image/png"
	"log"
	"net/http"
	"strconv"
	"sync"
)

var (
	mutex sync.Mutex
)

func (w *World) setKey(rw http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	//id := r.FormValue("id")

	//fmt.Println(key)
	switch key {
	case "-":
		w.Speed *= 3
	case " ":
		w.Speed /= 3
	case "*":
		w.balance = int(float64(w.balance) * 1.1)
	case "/":
		w.balance = int(float64(w.balance) / 1.1)
	}

	if w.Speed < 1 {
		w.Speed = 1
	}
	if w.Speed > 1000 {
		w.Speed = 1000
	}

	fmt.Println("Geneartion:", w.Gen)
	fmt.Println("Speed:", w.Speed)
	fmt.Println("Balance:", w.balance)

}

func loadHTML(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "index.html")
}

func (w *World) loadPict(rw http.ResponseWriter, r *http.Request) {
	mutex.Lock()
	w.imgChange()
	err := png.Encode(rw, w.image)
	if err != nil {
		log.Fatal("loadPict:", err)
	}
	mutex.Unlock()
}

func (w *World) ListenHTTP(port int) {
	http.HandleFunc("/pict/", w.loadPict)
	http.HandleFunc("/key/", w.setKey)
	http.HandleFunc("/", loadHTML)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenHTTP:", err)
	}
}
