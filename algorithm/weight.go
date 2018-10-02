package algorithm

import (
	"sync"
)

type SmoothWeight struct {
	objPool   []*weighted
	lock      sync.Mutex
	length    int
	lastIndex int
}

type weighted struct {
	obj             interface{}
	weight          int
	currentWeight   int
	effectiveWeight int
}

func NewSmoothWeight() *SmoothWeight {
	return &SmoothWeight{
		objPool:   make([]*weighted, 0),
		length:    0,
		lastIndex: 0,
	}
}

// 使用过权重值后，当前分就要减掉使用完的权重值，剩下的权重影响留下一次处理
func (sw *SmoothWeight) Get() (obj interface{}, index int) {
	sw.lock.Lock()
	defer sw.lock.Unlock()
	if sw.length > 0 {
		total := 0
		for i := 0; i < sw.length; i++ {
			w := sw.objPool[i]

			if w == nil {
				continue
			}

			w.currentWeight += w.effectiveWeight
			total += w.effectiveWeight
			if w.effectiveWeight < w.weight {
				w.effectiveWeight++
			}

			if w.currentWeight > sw.objPool[sw.lastIndex].currentWeight {
				sw.lastIndex = i
				return w.obj, i
			}
		}
		sw.objPool[sw.lastIndex].currentWeight -= total

		return sw.objPool[sw.lastIndex].obj, sw.lastIndex
	}

	return obj, -1
}

func (sw *SmoothWeight) Put(obj interface{}, weight ...int) {
	sw.lock.Lock()
	if len(weight) > 0 && weight[0] > 0 {
		w := weighted{
			obj:             obj,
			weight:          weight[0],
			currentWeight:   0,
			effectiveWeight: 0,
		}
		sw.objPool = append(sw.objPool, &w)
		sw.length++
	}
	sw.lock.Unlock()
}

func (sw *SmoothWeight) Del(index int) {
	sw.lock.Lock()
	if index < sw.length {
		sw.objPool = append(sw.objPool[:index], sw.objPool[index+1:]...)
		sw.lastIndex = 0
		sw.length--
	}
	sw.lock.Unlock()
}

func (sw *SmoothWeight) GetAll() []interface{} {
	sw.lock.Lock()
	objs := []interface{}{}
	for _, obj := range sw.objPool {
		objs = append(objs, obj.obj)
	}
	sw.lock.Unlock()
	return objs
}
