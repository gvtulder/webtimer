package main

import (
	"embed"
	"flag"
	"io/fs"
	"log"
	"net"
	"os"
	"strings"
	"webtimer/server/server"
)

var addr string
var frontendPath string

//go:embed dist/frontend/index.html
//go:embed dist/frontend/timer.html
//go:embed dist/frontend/normalize.css
//go:embed dist/frontend/style.css
//go:embed dist/frontend/main.js
//go:embed dist/frontend/inter.woff2
var embedFrontend embed.FS

func init() {
	flag.StringVar(&addr, "addr", ":8000", "http service address")
	flag.StringVar(&frontendPath, "web", "", "frontend resources path")
	flag.Parse()
}

func listAddresses(logger *log.Logger) []net.IP {
	var ipAddresses []net.IP

	interfaces, err := net.Interfaces()
	if err != nil {
		log.Fatal(err)
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			logger.Println("error getting addresses for interface", iface.Name, ":", err)
			continue
		}

		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err == nil {
				ipAddresses = append(ipAddresses, ip)
			}
		}
	}

	return ipAddresses
}

func main() {
	logger := log.New(os.Stdout, "", log.LstdFlags)

	frontend, err := fs.Sub(embedFrontend, "dist/frontend")
	if err != nil {
		logger.Fatal(err)
	}

	if frontendPath != "" {
		frontend = os.DirFS(frontendPath)
	}

	parts := strings.Split(addr, ":")
	if len(parts) == 1 {
		log.Fatal("Invalid listening address.")
	} else if len(parts) == 2 && parts[0] == "" {
		for _, a := range listAddresses(logger) {
			logger.Printf("Running at http://%s:%s/", a, parts[len(parts)-1])
		}
	} else {
		logger.Printf("Running at http://%s/", addr)
	}

	server.RunServer(addr, frontend, logger)
}
