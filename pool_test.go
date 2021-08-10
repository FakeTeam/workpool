package workpool

import (
	"fmt"
	"sync"
	"testing"
	"time"
	"workpool/future"

	"github.com/stretchr/testify/assert"
)

func TestWorkPool_Start(t *testing.T) {
	wg := sync.WaitGroup{}
	wp := GenWorkPool(ForeverType, 2, 100)
	wp.Start()
	lenth := 100
	wg.Add(lenth)
	for i := 0; i < lenth; i++ {
		wp.SubmitTask(&PoolTask{withoutFunc: func(args ...interface{}) {
			defer wg.Done()
			fmt.Printf("%d ", args[0].(int))
		}, args: []interface{}{i, "23", 4}})
	}
	wg.Wait()
}

func BenchmarkWorkPool(b *testing.B) {
	wp := GenWorkPool(ForeverType, 200, 10000)
	wp.Start()
	for i := 0; i < b.N; i++ {
		wp.SubmitTask(&PoolTask{withoutFunc: func(args ...interface{}) {
			time.Sleep(time.Microsecond * time.Duration(200))
		}, args: []interface{}{b.N}})
	}
}

func TestFuPool(t *testing.T) {
	wp := GenWorkPool(FutureType, 2, 100)
	var reulst []interface{}
	lenth := 100
	except := 5050 * 2
	for i := 1; i <= lenth; i++ {
		fu := wp.SubmitTask(&PoolTask{WithFunc: func(args ...interface{}) interface{} {
			return args[0].(int) * 2
		}, args: []interface{}{i}})
		reulst = append(reulst, fu)
	}
	wp.Start()
	sum := 0
	for _, v := range reulst {
		sum += v.(*future.Future).Get().(int)
	}
	assert.True(t, sum == except)
}
