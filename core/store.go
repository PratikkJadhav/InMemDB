package core

import (
	"time"

	"github.com/PratikkJadhav/InMemDB/config"
)

var store map[string]*Obj
var expires map[*Obj]uint64

func init() {
	store = make(map[string]*Obj)
	expires = make(map[*Obj]uint64)

}

func SetExpiry(obj *Obj, expDurationMS uint64) {
	expires[obj] = uint64(time.Now().UnixMilli()) + uint64(expDurationMS)
}

func NewObj(value interface{}, expdurationMS int64, otype uint8, oEnc uint8) *Obj {

	obj := &Obj{
		Value:          value,
		TypeEncoding:   otype | oEnc,
		lastAccessedAt: getCurrentClock(),
	}

	if expdurationMS > 0 {
		SetExpiry(obj, uint64(expdurationMS))
	}

	return obj

}

func Put(k string, obj *Obj) {

	if len(store) >= config.KeysLimit {
		evict()
	}
	obj.lastAccessedAt = getCurrentClock()
	store[k] = obj

	if keyKeyspaceStat[0] == nil {
		keyKeyspaceStat[0] = make(map[string]int)
	}

	keyKeyspaceStat[0]["keys"]++
}

func Get(k string) *Obj {
	v := store[k]

	if v != nil {
		if hasExpired(v) {
			Del(k)
			return nil
		}
	}

	v.lastAccessedAt = getCurrentClock()

	return v
}

func Del(k string) bool {
	if obj, ok := store[k]; ok {
		delete(expires, obj)
		return true
	}

	keyKeyspaceStat[0]["keys"]--
	return false
}
