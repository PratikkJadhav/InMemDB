package core

import (
	"bytes"
	"errors"
	"fmt"
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
	oType, oEnc := deduceTypeEncoding(value)

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
				return Encode(errors.New("ERR value is not an integer or out of range"), false)
			}

			ExecMSec = ExecSec * 1000

		default:
			return Encode(errors.New("Syntax Error"), false)
		}
	}

	Put(key, NewObj(value, ExecMSec, oType, oEnc))
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

	if hasExpired(obj) {
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

	exp, isExpiredSet := getExpiry(obj)
	if !isExpiredSet {
		return RESP_MINUS_ONE
	}

	if exp < uint64(time.Now().UnixMilli()) {
		return RESP_MINUS_TWO
	}
	durationMs := exp - uint64(time.Now().UnixMilli())
	return Encode(int64(durationMs/1000), false)

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
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}

	obj := Get(key)

	if obj == nil {
		return RESP_ZERO
	}

	SetExpiry(obj, uint64(executionDur*1000))

	return RESP_ONE
}

func evalBGREWRITEAOF(args []string) []byte {
	DumpAllAOF()
	return RESP_OK
}

func evalINCR(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("Err wrong number of arguments for 'incr' command"), false)
	}

	var key string = args[0]

	obj := Get(key)

	if obj == nil {
		obj = NewObj("0", -1, OBJ_TYPE_STRING, OBJ_ENCODING_INT)
		Put(key, obj)
	}

	if err := assertType(obj.TypeEncoding, OBJ_TYPE_STRING); err != nil {
		return Encode(err, false)
	}

	if err := assertEncoding(obj.TypeEncoding, OBJ_ENCODING_INT); err != nil {
		return Encode(err, false)
	}

	i, _ := strconv.ParseInt(obj.Value.(string), 10, 64)
	i++
	obj.Value = strconv.FormatInt(i, 10)

	return Encode(i, false)
}

func evalINFO(args []string) []byte {
	var info []byte

	buf := bytes.NewBuffer(info)
	buf.WriteString("# Keyspace\r\n")

	for i := range keyKeyspaceStat {
		buf.WriteString(fmt.Sprintf("db%dkeys=%d,expired=0,avg_ttl=0\r\n", i, keyKeyspaceStat[i]["keys"]))
	}

	return Encode(buf.String(), false)
}
func evalCLIENT(args []string) []byte {
	return RESP_OK
}

func evalLATENCY(args []string) []byte {
	return Encode([]string{}, false)
}

func evalLRU(args []string) []byte {
	evictAllkeysLRU()
	return RESP_OK
}

func evalSLEEP(args []string) []byte {
	if len(args) != 1 {
		return Encode(errors.New("ERR wrong number of arguments for 'SLEEP' command"), false)
	}

	durationSec, err := strconv.ParseInt(args[0], 10, 64)

	if err != nil {
		return Encode(errors.New("ERR value is not an integer or out of range"), false)
	}

	time.Sleep(time.Duration(durationSec) * time.Second)
	return RESP_OK
}
func EvalAndRespond(cmds RedisCmds, c io.ReadWriter) {
	var response []byte
	buf := bytes.NewBuffer(response)

	for _, cmd := range cmds {
		switch cmd.Cmd {
		case "PING":
			buf.Write(evalPing(cmd.Args))
		case "SET":
			buf.Write(evalSET(cmd.Args))
		case "GET":
			buf.Write(evalGET(cmd.Args))
		case "TTL":
			buf.Write(evalTTL(cmd.Args))
		case "DEL":
			buf.Write(evalDEL(cmd.Args))
		case "EXPIRE":
			buf.Write(evalExpire(cmd.Args))

		case "BGREWRITEAOF":
			buf.Write(evalBGREWRITEAOF(cmd.Args))

		case "INCR":
			buf.Write(evalINCR(cmd.Args))

		case "INFO":
			buf.Write(evalINFO(cmd.Args))

		case "CLIENT":
			buf.Write(evalCLIENT(cmd.Args))

		case "LATENCY":
			buf.Write(evalLATENCY(cmd.Args))

		case "LRU":
			buf.Write(evalLRU(cmd.Args))

		case "SLEEP":
			buf.Write(evalSLEEP(cmd.Args))
		default:
			buf.Write(evalPing(cmd.Args))
		}
	}
	c.Write(buf.Bytes())
}
