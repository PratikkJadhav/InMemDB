package core

import (
	"time"

	"github.com/PratikkJadhav/Redigo/config"
)

type Obj struct {
	Value     interface{}
	ExpiresAt int64
}

var store map[string]*Obj

func init() {
	store = make(map[string]*Obj)
}

func NewObj(value interface{}, durationMS int64) *Obj {
	var expiresAt int64 = -1

	if durationMS > 0 {
		expiresAt = durationMS + time.Now().UnixMilli()
	}

	return &Obj{
		Value:     value,
		ExpiresAt: expiresAt,
	}

}

func Put(k string, obj *Obj) {

	if len(store) >= config.KeysLimit {
		evict()
	}
	store[k] = obj
}

func Get(k string) *Obj {
	v := store[k]

	if v != nil {
		if v.ExpiresAt != -1 && v.ExpiresAt <= time.Now().UnixMilli() {
			delete(store, k)
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

	return false
}
