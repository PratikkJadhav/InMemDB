package core

const (
	OBJ_STRING uint8 = 0
	OBJ_LIST   uint8 = 1
	OBJ_HASH   uint8 = 2
)

type Obj struct {
	Type           uint8
	TypeEncoding   uint8
	Value          interface{}
	lastAccessedAt uint32
}

var OBJ_TYPE_STRING uint8 = 0 << 4
var OBJ_TYPE_HASH uint8 = 2 << 4

var OBJ_ENCODING_RAW uint8 = 0
var OBJ_ENCODING_INT uint8 = 1
var OBJ_ENCODING_EMBSTR uint8 = 8
var OBJ_ENCODING_HT uint8 = 2
