package gofm

import (
	"sync"
	"testing"
)

func Test_audience_getCookie(t *testing.T) {
	a := &audience{
		roomID: 0,
	}

	a.getCookie()
}

func Test_audience_Connect(t *testing.T) {
	wg := &sync.WaitGroup{}
	wg.Add(1)
	for i := 0; i < 10; i++ {
		a := &audience{
			roomID: 176719728,
		}
		a.Connect()
	}

	wg.Wait()
}
