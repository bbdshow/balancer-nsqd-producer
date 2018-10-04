package algorithm

import (
	"sync"
)

// Polling 轮询算法
type Polling struct {
	objPool []interface{}
	lock    sync.Mutex
	length  int
	index   int
}

// NewPolling 轮询
func NewPolling() *Polling {
	return &Polling{
		objPool: make([]interface{}, 0),
		length:  0,
		index:   -1,
	}
}

// Get 通过轮询方式，获取对象
func (pl *Polling) Get() (obj interface{}, index int) {
	pl.lock.Lock()
	defer pl.lock.Unlock()
	if pl.length > 0 {
		pl.index++
		if pl.index < pl.length {
			obj = pl.objPool[pl.index]
		} else {
			pl.index = 0
			obj = pl.objPool[0]
		}
	}

	return obj, pl.index
}

// Put 写入到对象池
func (pl *Polling) Put(obj interface{}, weight ...int) {
	pl.lock.Lock()
	pl.objPool = append(pl.objPool, obj)
	pl.length++
	pl.lock.Unlock()
}

// Del 删除对象池 index 索引对象
func (pl *Polling) Del(index int) {
	pl.lock.Lock()
	if index < pl.length {
		pl.objPool = append(pl.objPool[:index], pl.objPool[index+1:]...)
		pl.length--
	}
	pl.lock.Unlock()
}

// GetAll 获取所有对象
func (pl *Polling) GetAll() []interface{} {
	pl.lock.Lock()
	objs := pl.objPool
	pl.lock.Unlock()
	return objs
}
