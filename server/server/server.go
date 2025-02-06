package server

import (
	"io/fs"
	"log"
	"net/http"
	"webtimer/server/model"
)

func addCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func RunServer(addr string, frontend fs.FS, watch *model.TimerWatch, timer *model.Timer) {
	log.SetFlags(0)
	hub := newHub()
	go hub.run(watch, timer)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, w, r)
	})
	http.Handle("/", addCacheHeaders(http.FileServerFS(frontend)))
	log.Fatal(http.ListenAndServe(addr, nil))
}
