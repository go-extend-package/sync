package sync

import (
	"sync"
	"sync/atomic"
)

type Once struct {
	done uint32
	m    sync.Mutex
}

func (o *Once) Do(f func() error) (err error) {
	if atomic.LoadUint32(&o.done) == 0 {
		return o.doSlow(f)
	}
	return
}

func (o *Once) doSlow(f func() error) (err error) {
	o.m.Lock()
	defer o.m.Unlock()
	if o.done == 0 {
		if err = f(); err == nil {
			defer atomic.StoreUint32(&o.done, 1)
		}
	}
	return
}
