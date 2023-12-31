package httputils

import (
	"math/rand"
	"time"
)

type BackoffImpl struct {
	Initial, Max, cur time.Duration
}

func (b *BackoffImpl) Pause() time.Duration {
	d := time.Duration(1 + rand.Int63n(int64(b.cur)))
	b.cur = time.Duration(float64(b.cur) * 2)
	if b.cur > b.Max {
		b.cur = b.Max
	}
	return d
}

func NewBackoffImpl() *BackoffImpl {
	return &BackoffImpl{Initial: 100 * time.Millisecond}
}
