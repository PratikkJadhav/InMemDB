package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

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

	var sigs chan os.Signal = make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)
	var wg sync.WaitGroup
	wg.Add(2)

	go server.RunAsyncTCPServer(&wg)
	go server.WaitforSignal(&wg, sigs)

	wg.Wait()
}
