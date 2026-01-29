package server

import (
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/PratikkJadhav/Redigo/config"
	"github.com/PratikkJadhav/Redigo/core"
)

func readCommands(c io.ReadWriter) (*core.RedisCmd, error) {

	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	tokens, err := core.DecodeArrayString(buf[:n])
	if err != nil {
		return nil, err
	}

	return &core.RedisCmd{
		Cmd:  strings.ToUpper(tokens[0]),
		Args: tokens[1:],
	}, nil
}

func repondError(err error, c io.ReadWriter) {
	c.Write([]byte(fmt.Sprintf("-%s\r\n", err)))
}

func respond(cmd *core.RedisCmd, c io.ReadWriter) {
	err := core.EvalAndRespond(cmd, c)
	if err != nil {
		repondError(err, c)
	}
}

func RunSyncTCPServer() {
	log.Println("Starting a sync tcp server on", config.Host, config.Port)

	var con_clients = 0

	lsnr, err := net.Listen("tcp", config.Host+":"+strconv.Itoa(config.Port))
	if err != nil {
		// panic(err)
		log.Println(err)
		return
	}

	for {

		c, err := lsnr.Accept()
		if err != nil {
			// panic(err)
			log.Println(err)
		}

		con_clients += 1

		for {

			cmd, err := readCommands(c)
			if err != nil {
				c.Close()
				con_clients -= 1

				if err == io.EOF {
					break
				}
			}

			respond(cmd, c)
		}

	}
}
