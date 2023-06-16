package httputils

import (
	"math/rand"
	"time"
)

type BackoffImpl struct {
	Initial, Max, cur time.Duration
}

func (b *BackoffImpl) Pause() time.Duration {
	if b.cur == 0 {
		b.cur = b.Initial
	}
	interval := time.Duration(1 + rand.Int63n(int64(b.cur)))
	b.cur = time.Duration(b.cur * 2)
	if b.cur > b.Max {
		b.cur = b.Max
	}
	return interval
}

func NewBackoffImpl() *BackoffImpl {
	return &BackoffImpl{Initial: 100 * time.Millisecond}
}
