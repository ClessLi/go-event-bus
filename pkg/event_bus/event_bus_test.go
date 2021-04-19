package event_bus

import (
	"fmt"
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

func TestNewEventBus(t *testing.T) {
	eb := NewEventBus()
	eb.Register(new(TestObserverA))
	eb.Register(new(TestObserverB))
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
