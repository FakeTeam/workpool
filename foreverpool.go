package workpool

import (
	"context"
	"sync"
)

type Task struct {
	f    PoolWithoutReturnFunc
	args []interface{}
}

func (t *Task) execute() {
	t.f(t.args...)
}

type ForeverPool struct {
	tasks         chan *Task
	size          uint32
	stopCtx       context.Context
	stopCtxCancel context.CancelFunc
	wg            sync.WaitGroup
}

func NewForeverPool(size, len uint32) *ForeverPool {
	return &ForeverPool{
		size:  size,
		tasks: make(chan *Task, len),
	}
}

func (pool *ForeverPool) SubmitTask(pkt *PoolTask) interface{} {
	pool.tasks <- &Task{
		f:    pkt.withoutFunc,
		args: pkt.args,
	}
	return nil
}

func (pool *ForeverPool) work() {
	for {
		select {
		case <-pool.stopCtx.Done():
			pool.wg.Done()
			return
		case t := <-pool.tasks:
			t.execute()
		}
	}
}

func (pool *ForeverPool) Start() {
	pool.wg.Add(int(pool.size))
	pool.stopCtx, pool.stopCtxCancel = context.WithCancel(context.Background())
	for i := uint32(0); i < pool.size; i++ {
		go pool.work()
	}
}

func (pool *ForeverPool) Stop() {
	pool.stopCtxCancel()
	pool.wg.Wait()
}
