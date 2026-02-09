package core

import "sort"

var ePoolSizeMax int = 16

type PoolItem struct {
	key            string
	lastAccessedAt uint32
}

type evictionPool struct {
	pool   []*PoolItem
	keyset map[string]*PoolItem
}

type Byideltime []*PoolItem

func (a Byideltime) Len() int {
	return len(a)
}

func (a Byideltime) Swap(i, j int) {
	a[i], a[j] = a[j], a[i]
}

func (a Byideltime) Less(i, j int) bool {
	return getIdealTime(a[i].lastAccessedAt) > getIdealTime(a[j].lastAccessedAt)
}

func (pq *evictionPool) Push(key string, lastAccessedAt uint32) {
	_, ok := pq.keyset[key]

	if ok {
		return
	}

	item := &PoolItem{key: key, lastAccessedAt: lastAccessedAt}

	if len(pq.pool) < ePoolSizeMax {
		pq.keyset[key] = item
		pq.pool = append(pq.pool, item)

		sort.Sort(Byideltime(pq.pool))
	} else if lastAccessedAt > pq.pool[0].lastAccessedAt {
		pq.pool = pq.pool[1:]
		pq.keyset[key] = item
		pq.pool = append(pq.pool, item)
	}
}

func (pq *evictionPool) Pop() *PoolItem {
	if len(pq.pool) == 0 {
		return nil
	}

	item := pq.pool[0]
	pq.pool = pq.pool[1:]
	delete(pq.keyset, item.key)
	return item
}

func NewEvictionPool(size int) *evictionPool {
	return &evictionPool{
		pool:   make([]*PoolItem, size),
		keyset: make(map[string]*PoolItem),
	}
}

var ePool *evictionPool = NewEvictionPool(0)
