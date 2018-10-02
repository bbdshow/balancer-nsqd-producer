package algorithm

import (
	"math/rand"
	"sync"
)

type Random struct {
	objPool []interface{}
	lock    sync.RWMutex
	length  int
}

func NewRandom() *Random {
	random := Random{
		objPool: make([]interface{}, 0),
		length:  0,
	}
	return &random
}

func (rd *Random) Get() (obj interface{}, index int) {
	rd.lock.RLock()
	if rd.length < 1 {
		return obj, -1
	}

	index = rand.Intn(rd.length)
	obj = rd.objPool[index]
	rd.lock.RUnlock()

	return obj, index
}

func (rd *Random) Put(obj interface{}, weight ...int) {
	rd.lock.Lock()
	rd.objPool = append(rd.objPool, obj)
	rd.length++
	rd.lock.Unlock()
}

func (rd *Random) Del(index int) {
	rd.lock.Lock()
	if index < rd.length {
		rd.objPool = append(rd.objPool[:index], rd.objPool[index+1:])
		rd.length--
	}
	rd.lock.Unlock()
}

func (rd *Random) GetAll() []interface{} {
	rd.lock.Lock()
	objs := rd.objPool
	rd.lock.Unlock()
	return objs
}
