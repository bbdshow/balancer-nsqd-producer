package main

import (
	producer "balancer-nsqd-producer"
	"fmt"
	"os"

	"github.com/nsqio/go-nsq"
)

func main() {
	var addrs = map[string]int{
		"192.168.1.104:4150": 2,
		"192.168.1.109:4150": 8,
	}

	var opt = producer.Options{
		Addrs:        addrs,
		Retry:        2,
		Mode:         producer.PollingMode,
		PingInterval: 1,
		PingTimeout:  1,
	}
	bl, _ := producer.NewBalancer(opt, nsq.NewConfig())
	count := 100
	for count > 0 {
		err := bl.Publish("test_polling", []byte(" "))
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		count--
	}
}
