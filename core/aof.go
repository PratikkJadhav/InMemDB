package core

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/PratikkJadhav/Redigo/config"
)

func dumpKey(fp *os.File, key string, obj *Obj) {
	cmd := fmt.Sprintf("SET %s %s", key, obj.Value)
	tokens := strings.Split(cmd, " ")

	fp.Write(Encode(tokens, false))
}

func DumpAllAOF() {
	fp, err := os.OpenFile(config.AOFFile, os.O_CREATE|os.O_WRONLY, os.ModeAppend)

	if err != nil {
		fmt.Print(err)
		return
	}
	log.Println("rewriting AOF file at", config.AOFFile)
	for key, obj := range store {
		dumpKey(fp, key, obj)
	}
	log.Println("AOF file rewrite complete")
}
