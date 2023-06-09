package sync_test

import (
	"errors"
	"testing"

	. "github.com/go-extend-package/sync"
)

type one uint32

func (o *one) Increment() error {
	*o++
	if *o == 1 {
		return errors.New("want again increment")
	}
	return nil
}

func run(t *testing.T, once *Once, o *one, c chan struct{}) {
	err := once.Do(func() error {
		err := o.Increment()
		v := *o
		if err != nil {
			if v != 1 {
				t.Errorf("once failed inside run: %d is not 1", v)
			}
		} else {
			if v != 2 {
				t.Errorf("once failed inside run: %d is not 2", v)
			}
		}
		return err
	})
	if err == nil {
		if v := *o; v != 2 {
			t.Errorf("once failed inside run: %d is not 2", v)
		}
	}
	c <- struct{}{}
}

func TestOnce(t *testing.T) {
	o := new(one)
	once := new(Once)
	c := make(chan struct{})
	const N = 1000
	for i := 0; i < N; i++ {
		go run(t, once, o, c)
	}
	for i := 0; i < N; i++ {
		<-c
	}
	if *o != 2 {
		t.Errorf("once failed outside run: %d is not 2", *o)
	}
}
