package workpool

import "workpool/future"

type PoolWithoutReturnFunc func(args ...interface{})
type PoolTask struct {
	withoutFunc PoolWithoutReturnFunc
	WithFunc    future.FutrueFunc
	args        []interface{}
}

type PoolType uint8

const (
	ForeverType PoolType = 0
	FutureType  PoolType = 1
)

type WorkPool interface {
	Start()
	Stop()
	SubmitTask(pkt *PoolTask) interface{}
}

func GenWorkPool(which PoolType, poolSize, taskChSize uint32) WorkPool {
	var pool WorkPool
	switch which {
	case ForeverType:
		pool = NewForeverPool(poolSize, taskChSize)
	case FutureType:
		pool = NewFuturePool(poolSize, taskChSize)
	}
	return pool
}
