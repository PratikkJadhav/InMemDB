package core

import "github.com/PratikkJadhav/InMemDB/config"

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
func evict() {
	switch config.EvictionStrategy {
	case "simple-first":
		evictFirst()

	case "allkeys-random":
		evictAllkeysRandom()
	}
}
