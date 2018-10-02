package main

import (
	producer "balancer-nsqd-producer"
	nsq "go-nsq"
)

func main() {
	var addrs = map[string]int{
		"192.168.1.104:4150": 2,
		"192.168.1.109:4150": 8,
	}

	var opt = Options{
		Addr:         addrs,
		Retry:        2,
		Mode:         PollingMode,
		PingInterval: 1,
		PingTimeout:  1,
	}
	bl, _ := producer.NewBalancer(opt nsq.NewConfig())
	count := 100
	for count > 0 {
		err := bl.Publish("test_polling", []byte(" "))
		if err != nil {
			t.Fatal(err)
		}
		count--
	}
}
