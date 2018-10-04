package algorithm

import (
	"math/rand"
	"sync"
)

// Random 随机算法
type Random struct {
	objPool []interface{}
	lock    sync.RWMutex
	length  int
}

// NewRandom 随机方式
func NewRandom() *Random {
	random := Random{
		objPool: make([]interface{}, 0),
		length:  0,
	}
	return &random
}

// Get 获取通过随机算法的对象
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

// Put 写入到对象池
func (rd *Random) Put(obj interface{}, weight ...int) {
	rd.lock.Lock()
	rd.objPool = append(rd.objPool, obj)
	rd.length++
	rd.lock.Unlock()
}

// Del 删除对象池 index 索引对象
func (rd *Random) Del(index int) {
	rd.lock.Lock()
	if index < rd.length {
		rd.objPool = append(rd.objPool[:index], rd.objPool[index+1:]...)
		rd.length--
	}
	rd.lock.Unlock()
}

// GetAll 获取所有对象
func (rd *Random) GetAll() []interface{} {
	rd.lock.Lock()
	objs := rd.objPool
	rd.lock.Unlock()
	return objs
}
