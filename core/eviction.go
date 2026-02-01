package core

import "github.com/PratikkJadhav/Redigo/config"

func evictFirst() {
	for key := range store {
		delete(store, key)
		return
	}
}

func evict() {
	switch config.EvictionStratergy {
	case "simple-first":
		evictFirst()
	}
}
