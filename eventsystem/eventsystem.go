package eventsystem

import "container/list"

type ComponentType uint64
type SystemType uint64

type Component interface {
	// 对象唯一id
	GetUid() uint64

	// 对象本身类型
	GetComponentType() ComponentType
	// 对象大类类型
	GetSystemType() SystemType
}

type System interface {
}

type Handle struct {
	group map[uint64]Component

	updates list.List

	// 第一版先随机 无序
	systems map[SystemType]map[]struct{}
}

func (h *Handle) RegisterSystem(component Component, isRegister bool) {
	if !isRegister {
		h.Remove(component.GetUid())
		return
	}
	if nil == h.group {
		h.group = map[uint64]Component{}
	}
	h.group[component.GetUid()] = component
	sysGroup := h.systems[component.GetSystemType()]
	if 0 == len(sysGroup) {
		return
	}
}

func (h *Handle) Remove(uid uint64) {

}

func (h *Handle) Update() {

}
