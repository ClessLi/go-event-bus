package executor

import (
	"github.com/ClessLi/go-event-bus/pkg/observer"
	"github.com/panjf2000/ants/v2"
	"sync"
)

type Executor interface {
	Execute(action observer.Action, event ...interface{})
}

func NewExecutor() Executor {
	exec := new(executor)
	return Executor(exec)
}

type executor struct {
}

func (e executor) Execute(action observer.Action, event ...interface{}) {
	action.Execute(event...)
}

type AsyncExecutor interface {
	Executor
	Wait()
}

type asyncExecutor struct {
	pool      *ants.Pool
	locker    sync.Mutex
	waitGroup *sync.WaitGroup
}

func NewAsyncExecutor(poolCap uint) AsyncExecutor {
	pool, err := ants.NewPool(int(poolCap))
	if err != nil {
		panic(err)
	}
	exec := &asyncExecutor{
		pool:      pool,
		locker:    sync.Mutex{},
		waitGroup: new(sync.WaitGroup),
	}
	return AsyncExecutor(exec)
}

func (a *asyncExecutor) Execute(action observer.Action, event ...interface{}) {
	a.locker.Lock()
	defer a.locker.Unlock()
	a.waitGroup.Add(1)
	_ = a.pool.Submit(func() {
		action.Execute(event...)
		a.waitGroup.Done()
	})
}

func (a *asyncExecutor) Wait() {
	a.waitGroup.Wait()
}
