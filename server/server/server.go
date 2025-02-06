package server

import (
	"io/fs"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func addCacheHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Cache-Control", "no-cache, no-store, must-revalidate")
		next.ServeHTTP(w, r)
	})
}

func RunServer(addr string, frontend fs.FS, logger *log.Logger) {
	hub := newHub(logger)
	go hub.run()

	r := mux.NewRouter()
	srv := &http.Server{
		Addr:         addr,
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	serveFile := func(name string) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			http.ServeFileFS(w, r, frontend, name)
		})
	}

	r.StrictSlash(true)
	r.Handle("/", addCacheHeaders(serveFile("index.html")))
	r.PathPrefix("/static/").Handler(addCacheHeaders(http.StripPrefix("/static/", http.FileServerFS(frontend))))
	r.Handle("/{key:[-a-zA-Z0-9]+}/", addCacheHeaders(serveFile("timer.html")))
	r.Handle("/{key:[-a-zA-Z0-9]+}/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, mux.Vars(r)["key"], w, r)
	}))
	logger.Fatal(srv.ListenAndServe())
}
