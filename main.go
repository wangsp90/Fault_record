package main

import (
	"cfg"
	"httpserver"
)

func main() {
	config := cfg.Getcfg()
	httpserver.Server(config)
}
