package event_bus

import (
	"github.com/ClessLi/go-event-bus/pkg/executor"
	"github.com/ClessLi/go-event-bus/pkg/observer"
)

type EventBus interface {
	Register(object interface{})
	//Unregister(object interface{})
	Post(event ...interface{})
	annotationGenerate()
}

func NewEventBus() EventBus {
	bus := &eventBusImp{
		registry: observer.NewRegistry(),
		executor: executor.NewExecutor(),
	}
	return EventBus(bus)
}

func NewEventBusWithExecutor(executor executor.Executor) EventBus {
	bus := &eventBusImp{
		registry: observer.NewRegistry(),
		executor: executor,
	}
	return EventBus(bus)
}

func NewAsyncEventBus(poolCap uint, waitFn *func()) EventBus {
	exec := executor.NewAsyncExecutor(poolCap)
	*waitFn = exec.Wait
	return NewEventBusWithExecutor(exec)
}

type eventBusImp struct {
	//annotation annotation.Annotation
	registry observer.Registry
	executor executor.Executor
}

func (e *eventBusImp) Register(object interface{}) {
	e.registry.Register(object)
}

//func (e *eventBusImp) Unregister(object interface{}) {
//	panic("implement me")
//}

func (e *eventBusImp) Post(event ...interface{}) {
	observerActions := e.registry.GetMatchedObserverActions(event...)
	if observerActions != nil {
		for _, observerAction := range observerActions {
			e.executor.Execute(observerAction, event...)
		}
	}
}

func (e *eventBusImp) annotationGenerate() {
	panic("implement me")
}
