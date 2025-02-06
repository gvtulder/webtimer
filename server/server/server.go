package server

import (
	"log"
	"net/http"
	"webtimer/server/model"
)

func RunServer(addr string, watch *model.TimerWatch, timer *model.Timer) {
	log.SetFlags(0)
	hub := newHub()
	go hub.run(watch, timer)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	log.Fatal(http.ListenAndServe(addr, nil))
}
