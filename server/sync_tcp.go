package server

import (
	"io"
	"log"
	"net"
	"strconv"

	"github.com/PratikkJadhav/Redigo/config"
)

func readCommands(c net.Conn) (string, error) {

	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])

	if err != nil {
		return "", err
	}

	return string(buf[:n]), nil
}

func respond(cmd string, c net.Conn) error {

	if _, err := c.Write([]byte(cmd)); err != nil {
		return err
	}

	return nil
}

func RunSyncTCPServer() {
	log.Printf("Starting a sync tcp server on %s:%d", config.Host, config.Port)

	var con_clients = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		panic(err)
	}

	for {

		c, err := lsnr.Accept()
		if err != nil {
			panic(err)
		}

		con_clients += 1
		log.Printf("client connected with address: %s, concurrent client: %d", c.RemoteAddr(), con_clients)

		for {

			cmd, err := readCommands(c)
			if err != nil {
				c.Close()
				con_clients -= 1
				log.Printf("Client disconnected with address %s, concurrent client: %d", c.RemoteAddr(), con_clients)

				if err == io.EOF {
					break
				}
				log.Printf("err: %s", err)
			}

			log.Printf("command %s", cmd)

			if err = respond(cmd, c); err != nil {
				log.Printf("err write: %s", err)
			}
		}

	}
}
