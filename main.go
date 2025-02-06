package main

import (
	"embed"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"webtimer/server/model"
	"webtimer/server/server"
)

var addr string
var frontendPath string

//go:embed dist/frontend/index.html
//go:embed dist/frontend/style.css
//go:embed dist/frontend/main.js
//go:embed dist/frontend/inter.woff2
var embedFrontend embed.FS

func init() {
	flag.StringVar(&addr, "addr", ":8000", "http service address")
	flag.StringVar(&frontendPath, "web", "", "frontend resources path")
	flag.Parse()
}

func main() {
	frontend, err := fs.Sub(embedFrontend, "dist/frontend")
	if err != nil {
		log.Fatal(err)
	}

	if frontendPath != "" {
		frontend = os.DirFS(frontendPath)
	}

	timer := model.NewTimer()
	watch := model.NewTimerWatch(timer)

	fmt.Println("Running at " + addr)
	watch.Start()
	server.RunServer(addr, frontend, watch, timer)
}
