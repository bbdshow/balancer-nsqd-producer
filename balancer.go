package producer

import (
	"balancer-nsqd-producer/algorithm"
	"errors"
	"fmt"
	"time"

	"github.com/nsqio/go-nsq"
)

type (
	Options struct {
		Addrs        map[string]int
		Retry        int
		Mode         BalanceMode
		PingInterval int // s
		PingTimeout  int
	}

	ProducerConn struct {
		producer *nsq.Producer
		connAddr string
		errCount int
		weight   int
	}

	BalanceMode int

	algorithmMode interface {
		GetAll() (objPool []interface{})
		Get() (obj interface{}, index int)
		Put(obj interface{}, weight ...int)
		Del(int)
	}
)

const (
	PollingMode BalanceMode = iota
	RandomMode
	SmoothWeightMode
)

var (
	ErrNoAvailableConns = errors.New("no available conns")
)

type Balancer struct {
	addrs        map[string]int // string nsqd 地址  int 权重 非权重 int = 0
	balanceWay   algorithmMode
	errConns     chan *ProducerConn // 存在异常的链接， 等待 retryConn 重试链接
	ErrorsChan   chan error         // nsqd 链接错误
	pingInterval int                // retryConn 间隔
	pingTimeout  int                // pingTimeout ／ retryConn 连续 ping 这么多次，则返回  ErrorsChan
	retry        int                //如果连续未成功则返回 error
}

// 非权重 int = 0
func NewBalancer(opt Options, config *nsq.Config) (*Balancer, error) {
	if err := Validate(opt); err != nil {
		return nil, err
	}

	bl := Balancer{
		addrs:        opt.Addrs,
		pingInterval: opt.PingInterval,
		pingTimeout:  opt.PingTimeout,
		errConns:     make(chan *ProducerConn, 100),
		ErrorsChan:   make(chan error, 100),
		retry:        opt.Retry,
	}

	bl.setAlgorithm(opt.Mode)
	if config == nil {
		return nil, fmt.Errorf("config nil")
	}

	for addr, wright := range opt.Addrs {
		pd, err := nsq.NewProducer(addr, config)
		if err != nil {
			return nil, err
		}

		pdConn := ProducerConn{
			producer: pd,
			connAddr: addr,
			errCount: 0,
			weight:   0,
		}

		bl.balanceWay.Put(&pdConn, wright)
	}

	bl.retryConns()

	return &bl, nil
}

func (bl *Balancer) setAlgorithm(mode BalanceMode) {
	switch mode {
	case PollingMode:
		bl.balanceWay = algorithm.NewPolling()
	case RandomMode:
		bl.balanceWay = algorithm.NewRandom()
	case SmoothWeightMode:
		bl.balanceWay = algorithm.NewSmoothWeight()
	default:
		bl.balanceWay = algorithm.NewPolling()
	}
}

func Validate(opt Options) error {
	if len(opt.Addrs) == 0 {
		return errors.New("invalid addr")
	}
	for addr, weight := range opt.Addrs {
		if weight < 0 {
			return errors.New(addr + " invalid weight")
		}
	}

	if opt.Retry < 1 || opt.Retry > 50 {
		return errors.New("invalid retry")
	}

	if opt.PingInterval > 5 || opt.PingInterval < 1 {
		return errors.New("invalid pingInterval")
	}

	if opt.PingTimeout/opt.PingInterval <= 0 {
		return errors.New("invalid pingTimeout")
	}

	return nil
}

func (bl *Balancer) Publish(topic string, body []byte) error {
	retry := bl.retry
	var err error
get:
	retry--
	pd, index := bl.balanceWay.Get()
	if index <= -1 {
		return ErrNoAvailableConns
	}

	if pd == nil {
		if retry < 0 {
			return ErrNoAvailableConns
		}
		goto get
	}

	conn := pd.(*ProducerConn)
	if conn == nil {
		if retry < 0 {
			return ErrNoAvailableConns
		}
		bl.errConns <- conn
		bl.balanceWay.Del(index)
		goto get
	}

	err = conn.producer.Publish(topic, body)
	if err != nil {
		bl.errConns <- conn
		bl.balanceWay.Del(index)
		if retry < 0 {
			return err
		}
	}

	return nil
}

func (bl *Balancer) MultiPublish(topic string, body [][]byte) error {
	retry := bl.retry
	var err error
get:
	retry--
	pd, index := bl.balanceWay.Get()
	if index <= -1 {
		return ErrNoAvailableConns
	}

	if pd == nil {
		if retry < 0 {
			return ErrNoAvailableConns
		}
		goto get
	}

	conn := pd.(*ProducerConn)
	if conn == nil {
		if retry < 0 {
			return ErrNoAvailableConns
		}
		bl.errConns <- conn
		bl.balanceWay.Del(index)
		goto get
	}

	err = conn.producer.MultiPublish(topic, body)
	if err != nil {
		bl.errConns <- conn
		bl.balanceWay.Del(index)
		if retry < 0 {
			return err
		}
	}

	return nil
}

func (bl *Balancer) retryConns() {
	go func() {
		for {
			select {
			case conn := <-bl.errConns:
				if err := conn.producer.Ping(); err != nil {
					conn.errCount++
					if conn.errCount >= (bl.pingTimeout / bl.pingInterval) {
						bl.ErrorsChan <- err
						conn.errCount = 0
					}

					time.Sleep(time.Second * time.Duration(bl.pingInterval))
					bl.errConns <- conn
				} else {
					bl.balanceWay.Put(conn, conn.weight)
				}
			}
		}
	}()
}

func (bl *Balancer) CloseAll() {
	objs := bl.balanceWay.GetAll()
	for _, obj := range objs {
		conn := obj.(*ProducerConn)
		if conn != nil {
			conn.producer.Stop()
		}
	}
}
