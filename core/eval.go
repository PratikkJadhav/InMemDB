package core

import (
	"errors"
	"log"
	"net"
)

func evalPing(args []string, c net.Conn) error {
	var b []byte

	if len(args) >= 2 {
		return errors.New("Err wrong number of arguments for 'ping' command")
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}

	_, err := c.Write(b)
	return err
}

func EvalAndRespond(cmd *RedisCmd, c net.Conn) error {
	log.Printf("command", cmd.Cmd)
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)
	}
}
