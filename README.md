# balancer-nsqd-producer

[![Build Status](https://travis-ci.org/hopingtop/balancer-nsqd-producer.svg?branch=master)](https://travis-ci.org/hopingtop/balancer-nsqd-producer)
[![Go Report Card](https://goreportcard.com/badge/github.com/hopingtop/balancer-nsqd-producer)](https://goreportcard.com/report/github.com/hopingtop/balancer-nsqd-producer)
[![codecov](https://codecov.io/gh/hopingtop/balancer-nsqd-producer/branch/master/graph/badge.svg)](https://codecov.io/gh/hopingtop/balancer-nsqd-producer)

解决 nsq 消息队列，官方 go-nsq 生产者没有负载均衡的问题。

一般 nsqd 都会以集群的方式部署多台，通过负载可以将信息分散的写入 nsqd，有如下考虑：

1. 减少因为 nsqd 异常退出，产生内存数据未及时落盘，导致数据的丢失。
2. 平滑权重可有效解决 nsqd 机器性能不一致的缺点，提高机器使用率。
3. 生产数据 io 吞吐会更高，充分利用集群整体性能。

## 特性

1. 实现了 轮询、随机、平滑权重 3种常用负载算法，抽象为接口，可快速添加其他（有状态的负载）算法。
2. 算法性能较高，4 core cpu  30-50ns/op 不会存在瓶颈。
3. 链接失效，摘除负载队列，进入 ping 队列，恢复后自动加入负载。

#### 注意

因为 go-nsq 包的特性，在下不才，只实现了同步生产消息，支持方法 MultiPublish  Publish，一般情况够用，balancer 包支持方法可再扩充，
同时通过多个 goroutine 写入也可解决异步并发写入问题。

## 使用

	go get github.com/hopingtop/balancer-nsqd-producer

``` go
func main() {
	
	var addrs = map[string]int{
		"192.168.1.104:4150": 2,
		"192.168.1.109:4150": 8,
	}

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

```



