package server

import (
	"io"
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/PratikkJadhav/InMemDB/config"
	"github.com/PratikkJadhav/InMemDB/core"
)

func toArrayString(ai []interface{}) ([]string, error) {
	as := make([]string, len(ai))
	for i := range ai {
		as[i] = ai[i].(string)
	}

	return as, nil
}
func readCommands(c io.ReadWriter) ([]*core.RedisCmd, error) {

	var buf []byte = make([]byte, 512)

	n, err := c.Read(buf[:])

	if err != nil {
		return nil, err
	}

	values, err := core.Decode(buf[:n])
	if err != nil {
		return nil, err
	}

	var cmds []*core.RedisCmd = make([]*core.RedisCmd, 0)
	for _, value := range values {
		tokens, err := toArrayString(value.([]interface{}))

		if err != nil {
			return nil, err
		}

		cmds = append(cmds, &core.RedisCmd{
			Cmd:  strings.ToUpper(tokens[0]),
			Args: tokens[1:],
		})
	}

	return cmds, nil
}

func respond(cmds core.RedisCmds, c io.ReadWriter) {
	core.EvalAndRespond(cmds, c)
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

			cmds, err := readCommands(c)
			if err != nil {
				c.Close()
				con_clients -= 1

				if err == io.EOF {
					break
				}
			}

			respond(cmds, c)
		}

	}
}
