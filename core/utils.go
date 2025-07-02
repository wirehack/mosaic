package core

import "go.uber.org/dig"

var __di = dig.New()

func DI() *dig.Scope {
	return __di.Scope("root")
}

func Module[T any](name string) (proxy T, exists bool) {
	DI().Invoke(func(mp *ModuleProxy) {
		if ret, has := mp.Get(name); has {
			exists = true
			proxy = ret.(T)
		}
	})
	return
}
