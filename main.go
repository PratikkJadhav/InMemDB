package main

import (
	"flag"
	"log"

	"github.com/PratikkJadhav/InMemDB/config"
	"github.com/PratikkJadhav/InMemDB/server"
)

func setupFlag() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for InMemDB server")
	flag.IntVar(&config.Port, "port", 7379, "port for InMemDB server")
	flag.Parse()
}

func main() {
	setupFlag()
	log.Printf("Rollingg the dice")
	server.RunAsyncTCPServer()
}
