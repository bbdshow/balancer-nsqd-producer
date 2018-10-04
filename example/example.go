package main

import (
	producer "balancer-nsqd-producer"
	"fmt"
	"os"

	"github.com/nsqio/go-nsq"
)

func main() {
	var addrs = map[string]int{"127.0.0.1:4150": 2}

	var opt = producer.Options{
		Addrs:        addrs,
		Retry:        2,                    // 如果遇到当前链接发送失败，重试次数，建议与 链接地址数量一致
		Mode:         producer.PollingMode, // 算法方式 PollingMode RandomMode SmoothWeightMode
		PingInterval: 1,                    // 暂时失效链接进入 ping 队列， ping 的间隔时间
		PingTimeout:  1,                    // ping 此时间后，balancer 通过 ErrorsChan   chan error  返回 nsqd 链接错误， 使用者应该消费 ErrorsChan
	}
	// nsq.NewConfig() go-nsq 官方包 生产者配置

	bl, err := producer.NewBalancer(opt, nsq.NewConfig())
	if err != nil { // 初始化链接时，出错
		fmt.Println(err)
		os.Exit(1)
	}
	count := 100
	for count > 0 {
		err := bl.Publish("test_polling", []byte(" ")) // 使用方法与 go-nsq 官方包一致
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		count--
	}

	bl.CloseAll() // 关闭所有 nsqd 链接
}
