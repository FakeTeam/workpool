package future

type FutrueFunc func(args ...interface{}) interface{}

type Future struct {
	f    FutrueFunc
	args []interface{}
	c    chan interface{}
}

func GentFuture(f FutrueFunc, args ...interface{}) *Future {
	return &Future{f: f, args: args, c: make(chan interface{}, 1)}
}

func (future *Future) Invoke() {
	go func() {
		future.c <- future.f(future.args...)
	}()
}

func (future *Future) Get() interface{} {
	i := <-future.c
	return i
}
