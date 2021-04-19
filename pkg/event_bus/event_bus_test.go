package event_bus

import (
	"fmt"
	"github.com/ClessLi/go-annotation/pkg/v2/annotation"
	"testing"
	"time"
)

type xMsg string

type yMsg xMsg

type zMsg string

type TestObserverA struct {
}

//@EventBus
func (t *TestObserverA) F1(msg xMsg) {
	fmt.Println("A.F1():", msg)
}

//@EventBus
func (t *TestObserverA) F2(msg1 *xMsg, msg2 yMsg) {
	fmt.Println("A.F2()指针:", *msg1, msg2)
}

//@EventBus
func (t *TestObserverA) F3(msg zMsg) {
	fmt.Println("A.F3():", msg)
}

type TestObserverB struct {
}

//@EventBus
func (t *TestObserverB) Fa(msg xMsg) {
	time.Sleep(time.Second)
	fmt.Println("B.Fa()", msg)
}

//@EventBus
func (t *TestObserverB) Fb(msg1 *xMsg, msg2 yMsg) {
	fmt.Println("B.Fb()指针", *msg1, msg2)
}

//@EventBus
func (t *TestObserverB) Fc(msg1 xMsg, msg2 yMsg) {
	fmt.Println("B.Fc()非指针", msg1, msg2)
}

type TestObserverProxy struct {
}

func (t TestObserverProxy) GetProxyName() string {
	return "TestObserverProxy"
}

func (t TestObserverProxy) Before(delegate annotation.AnnotatedMethod) bool {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	params := delegate.GetParams()
	fmt.Printf("before check for method %v\n", delegate.GetMethodLocation())
	if len(params) > 0 {
		for i := 0; i < len(params); i++ {
			fmt.Printf("before check: method %v, param %d, type %v, value '%v'\n", delegate.GetMethodLocation(), i, params[i].Kind(), params[i])
		}
	}
	return true
}

func (t TestObserverProxy) After(delegate annotation.AnnotatedMethod) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	result := delegate.GetResult()
	fmt.Printf("after handle for method %v\n", delegate.GetMethodLocation())
	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			fmt.Printf("after handle: method %v, result %d, type %v, value '%v'\n", delegate.GetMethodLocation(), i, result[i].Kind(), result[i])
		}
	}
}

func (t TestObserverProxy) Finally(delegate annotation.AnnotatedMethod) {
	defer func() {
		err := recover()
		if err != nil {
			fmt.Println(err)
		}
	}()
	params := delegate.GetParams()
	result := delegate.GetResult()
	fmt.Printf("finally handle for method %v\n", delegate.GetMethodLocation())
	if len(params) > 0 {
		for i := 0; i < len(params); i++ {
			fmt.Printf("finally handle: method %v, param %d, type %v, value '%v'\n", delegate.GetMethodLocation(), i, params[i].Kind(), params[i])
		}
	}
	if len(result) > 0 {
		for i := 0; i < len(result); i++ {
			fmt.Printf("finally handle: method %v, result %d, type %v, value '%v'\n", delegate.GetMethodLocation(), i, result[i].Kind(), result[i])
		}
	}
}

func TestNewEventBus(t *testing.T) {
	eb := NewEventBus()
	eb.Register(new(TestObserverA))
	eb.Register(new(TestObserverB))
	eb.RegisterProxy(new(TestObserverProxy))
	//fmt.Println(reflect.TypeOf(interface{}(msg1)).Name())

	for i := 1; i <= 20; i++ {
		msg1 := xMsg(fmt.Sprintf("%s-%d", "xMsg", i))
		msg2 := yMsg(fmt.Sprintf("%s-%d", "yMsg", i))
		msg3 := zMsg(fmt.Sprintf("%s-%d", "zMsg", i))
		eb.Post(msg1)
		eb.Post(msg2)
		eb.Post(msg3)
		eb.Post(msg1, msg2)
		eb.Post(&msg1, msg2)
	}
}

func TestNewAsyncEventBus(t *testing.T) {
	var waitFn func()
	eb := NewAsyncEventBus(20, &waitFn)
	eb.Register(new(TestObserverA))
	eb.Register(new(TestObserverB))
	eb.RegisterProxy(new(TestObserverProxy))
	//fmt.Println(reflect.TypeOf(interface{}(msg1)).Name())

	for i := 1; i <= 20; i++ {
		msg1 := xMsg(fmt.Sprintf("%s-%d", "xMsg", i))
		msg2 := yMsg(fmt.Sprintf("%s-%d", "yMsg", i))
		msg3 := zMsg(fmt.Sprintf("%s-%d", "zMsg", i))
		eb.Post(msg1)
		eb.Post(msg2)
		eb.Post(msg3)
		eb.Post(msg1, msg2)
		eb.Post(&msg1, msg2)
	}

	waitFn()
}
