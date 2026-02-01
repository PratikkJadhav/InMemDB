package core

import (
	"errors"
	"io"
	"log"
	"strconv"
	"time"
)

var Resp_Nil []byte = []byte("$-1\r\n")

func evalPing(args []string, c io.ReadWriter) error {
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

func evalSET(args []string, c io.ReadWriter) error {
	if len(args) <= 1 {
		return errors.New("Err wrong number of arguments for 'SET' command")
	}

	var key, value string
	key, value = args[0], args[1]

	var ExecMSec int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if len(args) <= 3 {
				return errors.New("Syntax Error")
			}

			ExecSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return err
			}

			ExecMSec = ExecSec * 1000

		default:
			return errors.New("Syntax Error")
		}
	}

	Put(key, NewObj(value, ExecMSec))
	c.Write([]byte("+OK\r\n"))
	return nil

}

func evalGET(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("Err wrong number of arguements for 'GET' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		c.Write(Resp_Nil)
		return nil
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		c.Write(Resp_Nil)
		return nil
	}

	c.Write(Encode(obj.Value, false))
	return nil

}

func evalTTL(args []string, c io.ReadWriter) error {
	if len(args) != 1 {
		return errors.New("Err wrong number of arguements for 'TTL' command")
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	if obj.ExpiresAt == -1 {
		c.Write([]byte(":-1\r\n"))
		return nil
	}

	durationMS := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMS < 0 {
		c.Write([]byte(":-2\r\n"))
		return nil
	}

	c.Write(Encode(int64(durationMS/1000), false))
	return nil

}

func evalDEL(args []string, c io.ReadWriter) error {

	var countDeleted int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}

	c.Write(Encode(countDeleted, false))

	return nil
}

func evalExpire(args []string, c io.ReadWriter) error {
	if len(args) >= 1 {
		return errors.New("Invalid number of arguements")
	}

	var key string = args[0]

	executionDur, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return errors.New("(error) ERR value is not an integer or out of range")
	}

	obj := Get(key)

	if obj == nil {
		c.Write([]byte(":0\r\n"))
		return nil
	}

	obj.ExpiresAt = time.Now().UnixMilli() + executionDur*1000

	c.Write([]byte(":1\r\n"))
	return nil
}

func EvalAndRespond(cmd *RedisCmd, c io.ReadWriter) error {
	log.Println("command: \n", cmd.Cmd)
	switch cmd.Cmd {
	case "PING":
		return evalPing(cmd.Args, c)

	case "GET":
		return evalGET(cmd.Args, c)

	case "SET":
		return evalSET(cmd.Args, c)

	case "TTL":
		return evalTTL(cmd.Args, c)

	case "DEL":
		return evalDEL(cmd.Args, c)

	case "EXPIRE":
		return evalExpire(cmd.Args, c)

	default:
		return evalPing(cmd.Args, c)
	}
}
