package core

import (
	"time"

	"github.com/PratikkJadhav/InMemDB/config"
)

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMS int64, otype uint8, oEnc uint8) *Obj {
	var expiresAt int64 = -1

	if durationMS > 0 {
		expiresAt = durationMS + time.Now().UnixMilli()
	}

	return &Obj{
		Value:        value,
		TypeEncoding: otype | oEnc,
		ExpiresAt:    expiresAt,
	}

}

func Put(k string, obj *Obj) {

	if len(store) >= config.KeysLimit {
		evict()
	}
	store[k] = obj

	if keyKeyspaceStat[0] == nil {
		keyKeyspaceStat[0] = make(map[string]int)
	}

	keyKeyspaceStat[0]["keys"]++
}

func Get(k string) *Obj {
	v := store[k]

	if v != nil {
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli() {
			Del(k)
			return nil
		}
	}

	return v
}

func Del(k string) bool {
	if _, ok := store[k]; ok {
		delete(store, k)
		return true
	}

	keyKeyspaceStat[0]["keys"]--
	return false
}
