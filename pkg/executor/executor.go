package executor

import (
	"fmt"
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
	//a.waitGroup.Add(1)
	//go func() {
	//	action.Execute(event...)
	//	a.waitGroup.Done()
	//}()
	a.locker.Lock()
	defer a.locker.Unlock()
	for a.pool.Free() <= 0 {
		//time.Sleep(time.Millisecond)
	}
	err := a.pool.Submit(func() {
		action.Execute(event...)
	})
	if err != nil {
		fmt.Println("AsyncExecutor.Execute error,", err)
	}
}

func (a *asyncExecutor) Wait() {
	//a.waitGroup.Wait()
	for a.pool.Running() > 0 {
		//time.Sleep(time.Millisecond)
	}
}
