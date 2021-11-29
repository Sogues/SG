package obj

import "sync"

var componentFactory []sync.Pool

func RegComponentFactory(idx uint32, fn func()) {

}

type BaseEntity struct {
}

func (m *BaseEntity) AddComponent(componentType ComponentType, param interface{}) {
}
