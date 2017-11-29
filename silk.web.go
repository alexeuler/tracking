package main

import (
	"github.com/up-finder/silk.web/app/server"
	_ "github.com/up-finder/silk.web/app/signal"
)

func main() {
	server.Start()
}
