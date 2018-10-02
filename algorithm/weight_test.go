package algorithm

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestSmoothWeight(t *testing.T) {
	smoothWeight := NewSmoothWeight()

	var wvalue = map[string]int{
		"1": 1,
		"2": 2,
		"3": 3,
		"4": 9,
	}

	for k, v := range wvalue {
		smoothWeight.Put(k, v)
	}

	wg := sync.WaitGroup{}
	var addr1, addr2, addr3, addr4 int32 = 0, 0, 0, 0
	count := 10000
	for count > 0 {
		wg.Add(1)
		go func() {
			value, _ := smoothWeight.Get()
			switch value.(string) {
			case "1":
				atomic.AddInt32(&addr1, 1)
			case "2":
				atomic.AddInt32(&addr2, 1)
			case "3":
				atomic.AddInt32(&addr3, 1)
			case "4":
				atomic.AddInt32(&addr4, 1)
			default:
				t.Fatal("invalid value")
			}
			wg.Done()
		}()
		count--
	}

	wg.Wait()

	t.Log(addr1, addr2, addr3, addr4)
}