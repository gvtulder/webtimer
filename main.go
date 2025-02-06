package main

import (
	"flag"
	"fmt"
	"webtimer/server/model"
	"webtimer/server/server"
)

var addr = flag.String("addr", ":8000", "http service address")

func main() {
	timer := model.NewTimer()
	watch := model.NewTimerWatch(timer)

	fmt.Println("Running at " + *addr)
	watch.Start()
	server.RunServer(*addr, watch, timer)
}
