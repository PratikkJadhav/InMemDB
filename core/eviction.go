package core

import (
	"time"

	"github.com/PratikkJadhav/InMemDB/config"
)

func evictFirst() {
	for key := range store {
		Del(key)
		return
	}
}

func evictAllkeysRandom() {
	evictnum := int64(config.EvictionRatio * float64(config.KeysLimit))

	for key := range store {
		Del(key)
		evictnum--

		if evictnum <= 0 {
			break
		}
	}
}

func getCurrentClock() uint32 {
	return uint32(time.Now().Unix()) & 0x00FFFFFF
}

func getIdealTime(lastAccessedAt uint32) uint32 {
	c := getCurrentClock()

	if c >= lastAccessedAt {
		return c - lastAccessedAt
	}

	return (0x00FFFFFF - lastAccessedAt) + c
}

func populateEvictionPool() {
	sampleSize := 5

	for k := range store {
		ePool.Push(k, store[k].lastAccessedAt)
		sampleSize--
		if sampleSize == 0 {
			break
		}
	}
}

func evictAllkeysLRU() {
	populateEvictionPool()

	evictCount := int64(config.EvictionRatio * float64(config.KeysLimit))

	for i := 0; i < int(evictCount) && len(ePool.pool) > 0; i++ {
		item := ePool.Pop()
		if item == nil {
			return
		}

		Del(item.key)
	}
}
func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()

	case "allkeys-random":
		evictAllkeysRandom()

	case "allkeys-lru":
		evictAllkeysLRU()
	}

}
