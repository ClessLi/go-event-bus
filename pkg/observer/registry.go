package observer

import (
	"fmt"
	"github.com/ClessLi/go-annotation/pkg/v2/annotation"
	"reflect"
	"strings"
)

type Registry interface {
	Register(object interface{})
	RegisterProxy(proxy annotation.AnnotatedMethodProxy)
	GetMatchedObserverActions(event ...interface{}) []Action
	findAllObserverActions(observer interface{}) map[string][]Action
	getAnnotatedMethods(object interface{}) (methods []reflect.Method)
}

type registry struct {
	annotation  annotation.Annotation
	registryMap map[string][]Action
}

func NewRegistry() Registry {
	a := annotation.NewAnnotation("EventBus")
	return newRegistry(a)
}

func newRegistry(annotation annotation.Annotation) Registry {
	r := &registry{
		annotation:  annotation,
		registryMap: make(map[string][]Action),
	}
	return Registry(r)
}

func (r *registry) Register(object interface{}) {
	//r.annotation.RegisterAnnotatedObject(object)
	observerActions := r.findAllObserverActions(object)
	for eventType, eventActions := range observerActions {
		registeredEventActions, has := r.registryMap[eventType]
		if !has || registeredEventActions == nil {
			r.registryMap[eventType] = make([]Action, 0)
			registeredEventActions = r.registryMap[eventType]
		}
		registeredEventActions = append(registeredEventActions, eventActions...)
		r.registryMap[eventType] = registeredEventActions
	}
}

func (r *registry) RegisterProxy(proxy annotation.AnnotatedMethodProxy) {
	err := r.annotation.RegisterAnnotatedObjectProxy(proxy)
	if err != nil {
		fmt.Println(err)
	}
}

func (r *registry) GetMatchedObserverActions(event ...interface{}) []Action {
	if event == nil || len(event) == 0 {
		return nil
	}
	matchedObservers := make([]Action, 0)
	eventType := ""
	for _, param := range event {
		paramType := reflect.TypeOf(param)
		for paramType.Kind() == reflect.Ptr {
			paramType = paramType.Elem()
			eventType += "*"
		}
		eventType += paramType.Name() + ", "
	}
	eventType = strings.TrimSuffix(eventType, ", ")
	if _, has := r.registryMap[eventType]; has {
		matchedObservers = append(matchedObservers, r.registryMap[eventType]...)
	}
	return matchedObservers
}

func (r *registry) findAllObserverActions(observer interface{}) map[string][]Action {
	observerActions := make(map[string][]Action)
	for _, method := range r.getAnnotatedMethods(observer) {
		eventType := ""
		for i := 1; i < method.Type.NumIn(); i++ {
			paramType := method.Type.In(i)
			for paramType.Kind() == reflect.Ptr {
				paramType = paramType.Elem()
				eventType += "*"
			}
			eventType += paramType.Name() + ", "
		}
		eventType = strings.TrimSuffix(eventType, ", ")

		if _, has := observerActions[eventType]; !has {
			observerActions[eventType] = make([]Action, 0)
		}
		observerActions[eventType] = append(observerActions[eventType], NewAction(observer, method))
	}
	return observerActions
}

func (r *registry) getAnnotatedMethods(object interface{}) (methods []reflect.Method) {
	r.annotation.RegisterAnnotatedObject(object)
	objectType := reflect.TypeOf(object)
	receiveType := objectType
	pkgPth := receiveType.PkgPath()
	receiverName := objectType.Name()
	for receiveType.Kind() == reflect.Ptr {
		receiveType = receiveType.Elem()
		pkgPth = receiveType.PkgPath()
		receiverName = receiveType.Name()
	}
	annotatedMethods := make([]reflect.Method, 0)
	for i := 0; i < objectType.NumMethod(); i++ {
		method := objectType.Method(i)
		pkgList := strings.Split(pkgPth, "/")
		methodLocation := fmt.Sprintf("%s.%s.%s", pkgList[len(pkgList)-1], receiverName, method.Name)
		if method.Type.NumIn() == 0 {
			fmt.Printf("method(%s) need no parameter", methodLocation)
			continue
		}
		if _, has := r.annotation.GetAnnotatedMethodInfos()[methodLocation]; has {
			annotatedMethods = append(annotatedMethods, method)
		}
	}
	return annotatedMethods
}
