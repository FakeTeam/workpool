package workpool

import (
	"sync/atomic"
	"workpool/future"
)

type FuturePool struct {
	capacity  uint32
	workCnt   int32
	queueSize int64
	taskCh    chan *future.Future
}

func NewFuturePool(poolSize, taskChSize uint32) *FuturePool {
	return &FuturePool{
		capacity:  poolSize,
		workCnt:   0,
		queueSize: 0,
		taskCh:    make(chan *future.Future, taskChSize),
	}
}

func (pool *FuturePool) Start() {
	for {
		if atomic.LoadInt64(&pool.queueSize) == 0 {
			break
		}
		if uint32(atomic.LoadInt32(&pool.workCnt)) == pool.capacity {
			continue
		}
		fu := <-pool.taskCh
		atomic.AddInt64(&pool.queueSize, -1)
		atomic.AddInt32(&pool.workCnt, 1)
		go func(cnt *int32) {
			fu.Invoke()
			atomic.AddInt32(cnt, -1)
		}(&pool.workCnt)
	}
}

func (pool *FuturePool) Stop() {
}

func (pool *FuturePool) SubmitTask(pkt *PoolTask) interface{} {
	atomic.AddInt64(&pool.queueSize, 1)
	fu := future.GentFuture(pkt.WithFunc, pkt.args...)
	pool.taskCh <- fu
	return fu
}
