package main

import (
	"flag"
)

func main() {
	port := flag.Uint("Port", 7000, "Port for Executive server")
	app := NewExecutiveNode(uint16(*port))
	app.Run()
}
