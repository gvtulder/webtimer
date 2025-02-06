package server

import (
	"io/fs"
	"log"
	"net/http"
	"webtimer/server/model"
)

func RunServer(addr string, frontend fs.FS, watch *model.TimerWatch, timer *model.Timer) {
	log.SetFlags(0)
	hub := newHub()
	go hub.run(watch, timer)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", http.FileServerFS(frontend))
	log.Fatal(http.ListenAndServe(addr, nil))
}
