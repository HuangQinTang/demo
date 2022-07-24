package Injector

import (
	"reflect"
)

// BeanMapper 容器
type BeanMapper map[reflect.Type]reflect.Value

// add 加入容器
func (b BeanMapper) add(bean interface{}) {
	t := reflect.TypeOf(bean)
	if t.Kind() != reflect.Ptr { //todo 不是指针不处理(暂时)
		panic("require ptr object")
	}
	b[t] = reflect.ValueOf(bean)
}

// get 从容器中获取值
func (b BeanMapper) get(bean interface{}) reflect.Value {
	var t reflect.Type
	if bt, ok := bean.(reflect.Type); ok {
		t = bt
	} else {
		t = reflect.TypeOf(bean)
	}
	if v, ok := b[t]; ok {
		return v
	}

	//处理接口 继承
	for k, v := range b {
		if k.Implements(t) {
			return v
		}
	}
	return reflect.Value{}
}

//
//func (b BeanMapper)
