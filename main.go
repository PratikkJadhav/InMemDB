package main

import (
	"flag"
	"log"

	"github.com/PratikkJadhav/Redigo/config"
	"github.com/PratikkJadhav/Redigo/server"
)

func setupFlag() {
	flag.StringVar(&config.Host, "host", "0.0.0.0", "Host for Redigo server")
	flag.IntVar(&config.Port, "port", 7379, "port for redigo server")
	flag.Parse()
}

func main() {
	setupFlag()
	log.Printf("Rollingg the dice")
	server.RunSyncTCPServer()
}
