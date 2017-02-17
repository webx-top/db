package factory

import (
	"math/rand"
	"sync/atomic"
)

var DefaultSelecter = &RoundRobin{}

type Selecter interface {
	Select(length int) int
}

type Random struct{}

func (r *Random) Select(length int) int {
	return rand.Intn(length - 1)
}

type RoundRobin struct {
	count uint64
}

func (r *RoundRobin) Select(length int) int {
	return int((atomic.AddUint64(&r.count, 1) % uint64(length)))
}
