package producer

import (
	"testing"
	"time"

	"github.com/nsqio/go-nsq"
)

var addrs = map[string]int{"0.0.0.0:4150": 2}

var opt = Options{
	Addrs:        addrs,
	Retry:        2,
	Mode:         PollingMode,
	PingInterval: 1,
	PingTimeout:  1,
}

func TestPollingPublish(t *testing.T) {
	opt.Mode = PollingMode
	balancer, _ := NewBalancer(opt, nsq.NewConfig())
	count := 100
	for count > 0 {
		err := balancer.Publish("test_polling", []byte(" "))
		if err != nil {
			t.Fatal(err)
		}
		count--
	}
}

func TestRandomPublish(t *testing.T) {
	opt.Mode = RandomMode
	balancer, _ := NewBalancer(opt, nsq.NewConfig())
	count := 100
	for count > 0 {
		err := balancer.Publish("test_random", []byte(" "))
		if err != nil {
			t.Fatal(err)
		}
		count--
	}
}

func TestWeightPublish(t *testing.T) {
	opt.Mode = SmoothWeightMode
	balancer, _ := NewBalancer(opt, nsq.NewConfig())
	count := 100
	for count > 0 {
		err := balancer.Publish("test_weight", []byte(" "))
		if err != nil {
			t.Fatal(err)
		}
		count--
	}
}

func TestConnsErr(t *testing.T) {
	opt.Mode = PollingMode
	balancer, _ := NewBalancer(opt, nsq.NewConfig())
	go func() {
		for {
			select {
			case err := <-balancer.ErrorsChan:
				t.Log("err ping ok", err)
			}
		}
	}()

	count := 100
	for count > 0 {
		err := balancer.Publish("test_conns", []byte(" "))
		if err != nil {
			t.Log("close conns", err)
			break
		}

		if count == 20 {
			balancer.CloseAll()
		}

		time.Sleep(time.Millisecond * 100)
		count--
	}

}
