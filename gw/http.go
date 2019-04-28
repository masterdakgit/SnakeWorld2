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

func setKey(w http.ResponseWriter, r *http.Request) {
	key := r.FormValue("key")
	id := r.FormValue("id")

	fmt.Println(key, id)

	switch key {
	case "ArrowUp":
	case "ArrowDown":
	case "ArrowLeft":
	case "ArrowRight":
	}
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
	http.HandleFunc("/key/", setKey)
	http.HandleFunc("/", loadHTML)
	err := http.ListenAndServe(":"+strconv.Itoa(port), nil)
	if err != nil {
		log.Fatal("ListenHTTP:", err)
	}
}
