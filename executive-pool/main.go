package main

import (
	"flag"
	"fmt"
)

func main() {
	fmt.Println("yellowwwww")
	port := flag.Uint("port", 8080, "TCP Port Number for Executive-pool")
	gateway := flag.String("gateway", "http://127.0.0.1:7000", "Runnable Gateway")
	app := NewExecutivePoolServer(uint16(*port), *gateway)
	app.Run()
}
