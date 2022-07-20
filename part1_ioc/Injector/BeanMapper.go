package Injector

import "reflect"

type BeanMapper map[reflect.Type]reflect.Value

func (b BeanMapper) add(bean interface{}) {
	t := reflect.TypeOf(bean)
	if t.Kind() != reflect.Ptr {
		panic("require ptr object")
	}
	b[t] = reflect.ValueOf(bean)
}
