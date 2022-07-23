package Injector

import (
	"github.com/shenyisyn/goft-expr/src/expr"
	"reflect"
)

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

// BeanFactoryImpl 容器工厂对象，负责管理容器
type BeanFactoryImpl struct {
	beanMapper BeanMapper //容器
	ExprMap    map[string]interface{}
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper), ExprMap: make(map[string]interface{})}
}

func (b *BeanFactoryImpl) Set(vlist ...interface{}) {
	if vlist == nil || len(vlist) == 0 {
		return
	}
	for _, v := range vlist {
		b.beanMapper.add(v)
	}
}

func (b *BeanFactoryImpl) Get(v interface{}) interface{} {
	if v == nil {
		return nil
	}
	get_v := b.beanMapper.get(v)
	if get_v.IsValid() { //值不为空时转为interface
		return get_v.Interface()
	}
	return nil
}

// Apply 处理依赖注入
func (b *BeanFactoryImpl) Apply(bean interface{}) {
	if bean == nil {
		return
	}
	v := reflect.ValueOf(bean) //获取反射值对象
	if v.Kind() == reflect.Ptr || v.Kind() == reflect.Interface {
		v = v.Elem() //通过反射对象获取接口的值或者该指针所指向的值
	}
	if v.Kind() != reflect.Struct {
		return
	}
	for i := 0; i < v.NumField(); i++ { //遍历结构体字段
		field := v.Type().Field(i)
		if !v.Field(i).CanSet() || field.Tag.Get("inject") == "" { //字段不能访问的(首字母非大写)，不存在指定tag(约定为inject)，不进行依赖注入
			continue
		}

		//通过类型从容器中取值，如果容器中存在该类型的值，把该值反射赋予
		if get_v := b.Get(field.Type); get_v != nil {
			v.Field(i).Set(reflect.ValueOf(get_v))
			continue
		}
		//表达式注入,依赖goft-expr包(https://github.com/shenyisyn/goft-expr)
		if field.Tag.Get("inject") != "-" {
			ret := expr.BeanExpr(field.Tag.Get("inject"), b.ExprMap) //通过tag标签填写的表达式从 b.ExprMap获取该表达是对应的方法(定义在Config下)
			if ret == nil && ret.IsEmpty() {                         //ExprMap取值为空不处理
				continue
			}
			retValue := ret[0]   //约定，ExprMap里对应的方法只有一个放回置
			if retValue == nil { //值为空不处理
				continue
			}
			b.Set(retValue)                           //把依赖也存入容器，二次获取时直接从容器取
			v.Field(i).Set(reflect.ValueOf(retValue)) //反射赋值
		}
	}
}
