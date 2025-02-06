// Copyright (C) 2025 Gijs van Tulder
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <https://www.gnu.org/licenses/>.

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

func RunServer(addr string, frontend fs.FS, logger *log.Logger, version string) {
	hub := newHub(logger, version)
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
	r.Handle("/{key:[-a-zA-Z0-9]+}/", addCacheHeaders(serveFile("index.html")))
	r.Handle("/{key:[-a-zA-Z0-9]+}/ws", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		serveWs(hub, mux.Vars(r)["key"], w, r)
	}))
	logger.Fatal(srv.ListenAndServe())
}
