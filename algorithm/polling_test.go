package algorithm

import (
	"sync"
	"sync/atomic"
	"testing"
)

func TestPolling(t *testing.T) {
	polling := NewPolling()

	addrs := []string{"1", "2", "3"}

	for _, addr := range addrs {
		value := addr
		polling.Put(value)
	}

	wg := sync.WaitGroup{}
	var addr1 int32 = 0
	var addr2 int32 = 0
	var addr3 int32 = 0

	count := 3000

	for count > 0 {
		wg.Add(1)
		go func() {
			pd, _ := polling.Get()
			switch pd.(string) {
			case "1":
				atomic.AddInt32(&addr1, 1)
			case "2":
				atomic.AddInt32(&addr2, 1)
			case "3":
				atomic.AddInt32(&addr3, 1)
			}
			wg.Done()
		}()
		count--
	}

	wg.Wait()
	value := int32(1000)
	if addr1 != value || addr2 != value || addr3 != value {
		t.Fatal(addr1, addr2, addr3)
	}

	polling.Del(0)

	addr1 = 0
	addr2 = 0
	addr3 = 0
	count = 3000
	for count > 0 {
		wg.Add(1)
		go func() {
			pd, _ := polling.Get()
			switch pd.(string) {
			case "1":
				t.Fatal("del 1 error")
			case "2":
				atomic.AddInt32(&addr2, 1)
			case "3":
				atomic.AddInt32(&addr3, 1)
			}
			wg.Done()
		}()
		count--
	}
	wg.Wait()
	value = int32(1500)
	if addr2 != value || addr3 != value {
		t.Fatal(addr1, addr2, addr3)
	}
}
