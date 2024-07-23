package service

import (
	"reflect"
)

// Locator is a simple service-locator PoC to get us off the ground and in to
// development nice and quick.  It's not been fully tested or hardened.
//
// Service-Locator isn't always the most appropriate way of solving dependency
// scenarios as it hides dependencies and so can become an issue in larger
// scale systems.
type Locator struct {
	store map[reflect.Type]reflect.Value
}

func NewLocator() *Locator { return new(Locator) }

var DefaultLocator = &defaultLocator

var defaultLocator Locator

func (s *Locator) Register(value interface{}) {
	if s.store == nil {
		s.store = make(map[reflect.Type]reflect.Value)
	}
	k := reflect.TypeOf(value)
	v := reflect.ValueOf(value)
	s.store[k] = v
}

func (s *Locator) RegisterFunc(value func() interface{}) {
	s.Register(value())
}

func (s *Locator) Resolve(value interface{}) bool {
	k := reflect.TypeOf(value).Elem()
	kind := k.Kind()
	if kind == reflect.Ptr {
		k = k.Elem()
		kind = k.Kind()
	}
	for t, v := range s.store {
		if kind == reflect.Interface && t.Implements(k) {
			reflect.Indirect(reflect.ValueOf(value)).Set(v)
			return true
		} else if kind == reflect.Struct && k.AssignableTo(t.Elem()) {
			reflect.ValueOf(value).Elem().Set(v)
			return true
		}
	}
	return false
}

func Register(value interface{}) {
	DefaultLocator.Register(value)
}

func RegisterFunc(value func() interface{}) {
	DefaultLocator.RegisterFunc(value)
}

func Resolve(value interface{}) bool {
	return DefaultLocator.Resolve(value)
}
