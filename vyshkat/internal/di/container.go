
package di

import (
	"fmt"
	"reflect"
)

type Container struct {
	factories map[reflect.Type]func(*Container) any
	instances map[reflect.Type]any
}

func New() *Container {
	return &Container{
		factories: make(map[reflect.Type]func(*Container) any),
		instances: make(map[reflect.Type]any),
	}
}

func typeOf[T any]() reflect.Type {
	return reflect.TypeOf((*T)(nil)).Elem()
}

func Register[T any](c *Container, factory func(c *Container) T) {
	t := typeOf[T]()
	c.factories[t] = func(c *Container) any { return factory(c) }
	delete(c.instances, t)
}

func RegisterInstance[T any](c *Container, v T) {
	Register(c, func(*Container) T { return v })
}

func Resolve[T any](c *Container) T {
	t := typeOf[T]()
	if inst, ok := c.instances[t]; ok {
		return inst.(T)
	}
	factory, ok := c.factories[t]
	if !ok {
		panic(fmt.Sprintf("di: тип %v не зарегистрирован", t))
	}
	inst := factory(c)
	c.instances[t] = inst
	return inst.(T)
}
