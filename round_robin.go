package zcache

import (
	"sync"

	"github.com/zhiqiangxu/util/concurrent"
)

type RoundRobin[K comparable, V any] struct {
	sync.RWMutex
	cache          map[K]value[V]
	slots          []K
	slotsAllocated uint64
	allAssigned    bool
}

type value[V any] struct {
	v V
	i uint64
}

func NewRoundRobin[K comparable, V any](slots int) *RoundRobin[K, V] {

	if slots <= 0 {
		panic("slots <= 0")
	}

	return &RoundRobin[K, V]{
		cache: make(map[K]value[V]),
		slots: make([]K, slots),
	}

}

func (r *RoundRobin[K, V]) Has(k K) (ok bool) {

	r.RLock()

	_, ok = r.cache[k]

	r.RUnlock()

	return

}

func (r *RoundRobin[K, V]) Get(k K) (v V, ok bool) {

	r.RLock()

	value, ok := r.cache[k]
	if ok {
		v = value.v
	}

	r.RUnlock()

	return

}

func (r *RoundRobin[K, V]) Set(k K, v V) (existed bool) {

	r.Lock()

	ovalue, existed := r.cache[k]

	if existed {
		r.cache[k] = value[V]{v: v, i: ovalue.i}
	} else {
		slot := r.slotsAllocated % uint64(len(r.slots))
		if r.allAssigned {
			ok := r.slots[slot]
			delete(r.cache, ok)
		}

		r.cache[k] = value[V]{v: v, i: slot}
		r.slots[slot] = k
		r.slotsAllocated += 1
		if !r.allAssigned && int(slot) == len(r.slots)-1 {
			r.allAssigned = true
		}
	}

	r.Unlock()

	return

}

func (r *RoundRobin[K, V]) SlotsAllocated() (ret uint64) {

	r.RLock()

	ret = r.slotsAllocated

	r.RUnlock()

	return

}

func NewBucketRoundRobin[K comparable, V any](slots, shards int, hash func(K) uint32) *concurrent.Bucket[K, *RoundRobin[K, V]] {
	return concurrent.NewBucket(shards, func() *RoundRobin[K, V] { return NewRoundRobin[K, V](slots / shards) }, hash)
}
