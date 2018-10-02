package main

import (
	"balancer-nsqd-producer/algorithm"
	"fmt"
	"sync"
	"sync/atomic"
)

func raceRandmon() {
	random := algorithm.NewRandom()

	objs := []int{1, 2, 3, 4}
	for _, obj := range objs {
		random.Put(obj)
	}

	wg := sync.WaitGroup{}
	var addr1, addr2, addr3, addr4 int32 = 0, 0, 0, 0
	count := 10000
	for count > 0 {
		wg.Add(1)
		go func() {
			value, _ := random.Get()
			switch value.(int) {

			case 1:
				atomic.AddInt32(&addr1, 1)
			case 2:
				atomic.AddInt32(&addr2, 1)
			case 3:
				atomic.AddInt32(&addr3, 1)
			case 4:
				atomic.AddInt32(&addr4, 1)
			default:
				fmt.Println("value invalid")
			}
			wg.Done()
		}()
		count--
	}

	wg.Wait()

	fmt.Println(addr1, addr2, addr3, addr4)
}

func raceWeight() {
	smoothWeight := algorithm.NewSmoothWeight()

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
				fmt.Println("invalid value")
			}
			wg.Done()
		}()
		count--
	}

	wg.Wait()

	fmt.Println(addr1, addr2, addr3, addr4)
}
