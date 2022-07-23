package Injector

var BeanFactory *BeanFactoryImpl

func init() {
	BeanFactory = NewBeanFactory()
}

// BeanFactoryImpl 容器工厂对象，负责管理容器
type BeanFactoryImpl struct {
	beanMapper BeanMapper
}

func NewBeanFactory() *BeanFactoryImpl {
	return &BeanFactoryImpl{beanMapper: make(BeanMapper)}
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
