package core

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"time"
)

var RESP_Nil []byte = []byte("$-1\r\n")
var RESP_OK []byte = []byte("+OK\r\n")
var RESP_ZERO []byte = []byte(":0\r\n")
var RESP_ONE []byte = []byte(":1\r\n")
var RESP_MINUS_ONE []byte = []byte(":-1\r\n")
var RESP_MINUS_TWO []byte = []byte(":-2\r\n")

func evalPing(args []string) []byte {
	var b []byte

	if len(args) >= 2 {
		return Encode(errors.New("Err wrong number of arguments for 'ping' command"), false)
	}

	if len(args) == 0 {
		b = Encode("PONG", true)
	} else {
		b = Encode(args[0], false)
	}
	return b
}

func evalSET(args []string) []byte {
	if len(args) <= 1 {
		return Encode(errors.New("Err wrong number of arguments for 'SET' command"), false)
	}

	var key, value string
	key, value = args[0], args[1]

	var ExecMSec int64 = -1

	for i := 2; i < len(args); i++ {
		switch args[i] {
		case "EX", "ex":
			i++
			if len(args) <= 3 {
				return Encode(errors.New("Syntax Error"), false)
			}

			ExecSec, err := strconv.ParseInt(args[3], 10, 64)
			if err != nil {
				return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
			}

			ExecMSec = ExecSec * 1000

		default:
			return Encode(errors.New("Syntax Error"), false)
		}
	}

	Put(key, NewObj(value, ExecMSec))
	return RESP_OK

}

func evalGET(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("Err wrong number of arguements for 'GET' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		return RESP_Nil
	}

	if obj.ExpiresAt != -1 && obj.ExpiresAt <= time.Now().UnixMilli() {
		return RESP_Nil
	}

	return Encode(obj.Value, false)

}

func evalTTL(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("Err wrong number of arguements for 'TTL' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		return RESP_MINUS_TWO
	}

	if obj.ExpiresAt == -1 {
		return RESP_MINUS_ONE
	}

	durationMS := obj.ExpiresAt - time.Now().UnixMilli()

	if durationMS < 0 {
		return RESP_MINUS_TWO
	}

	return Encode(int64(durationMS/1000), false)

}

func evalDEL(args []string) []byte {

	var countDeleted int = 0

	for _, key := range args {
		if ok := Del(key); ok {
			countDeleted++
		}
	}

	return Encode(countDeleted, false)
}

func evalExpire(args []string) []byte {
	if len(args) >= 1 {
		return Encode(errors.New("Invalid number of arguements"), false)
	}

	var key string = args[0]

	executionDur, err := strconv.ParseInt(args[1], 10, 64)
	if err != nil {
		return Encode(errors.New("(error) ERR value is not an integer or out of range"), false)
	}

	obj := Get(key)

	if obj == nil {
		return RESP_ZERO
	}

	obj.ExpiresAt = time.Now().UnixMilli() + executionDur*1000

	return RESP_ONE
}

func EvalAndRespond(cmds RedisCmds, c io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)

	for _, cmd := range cmds {
		switch cmd.Cmd {
		case "PING":
			buf.Write(evalPing(cmd.Args))

		case "GET":
			buf.Write(evalGET(cmd.Args))

		case "SET":
			buf.Write(evalSET(cmd.Args))

		case "TTL":
			buf.Write(evalTTL(cmd.Args))

		case "DEL":
			buf.Write(evalDEL(cmd.Args))

		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args))

		default:
			buf.Write(evalPing(cmd.Args))
		}
	}

	c.Write(buf.Bytes())

}
