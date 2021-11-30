package obj

import (
	"fmt"
	"reflect"
	"time"

	"github.com/Sogues/ETForGo/types"
)

type EntityStatus uint8

const (
	StatusFromPool EntityStatus = 1 << iota
	StatusRegister
	StatusComponent
	StatusCreate
	StatusNone = 0
)

var (
	tempUid = uint64(time.Now().Unix())
)

func genUid() uint64 {
	tempUid++
	return tempUid
}

type (
	EntityImpl interface {
		EntityTypeId() types.EntityType
	}
	// todo 管理entity自身的类型id
	BaseEntity struct {
		uid    uint64
		status EntityStatus

		impl EntityImpl

		domain *BaseEntity // 作用域

		parent *BaseEntity // 父节点

		children map[uint64]*BaseEntity

		// 组件以类型为单位 禁止同类型组件重复添加
		components map[types.EntityType]*BaseEntity
	}
)

func (e *BaseEntity) GetBaseEntity() *BaseEntity { return e }
func (e *BaseEntity) GetUid() uint64             { return e.uid }
func (e *BaseEntity) SetUid(uid uint64)          { e.uid = uid }
func (e *BaseEntity) FromPool() bool {
	return 0 != e.status&StatusFromPool
}

func (e *BaseEntity) SetFromPool(from bool) {
	if from {
		e.status |= StatusFromPool
	} else {
		e.status &= ^StatusFromPool
	}
}
func (e *BaseEntity) IsRegister() bool {
	return 0 != e.status&StatusRegister
}
func (e *BaseEntity) SetRegister(register bool) {
	if register == e.IsRegister() {
		return
	}
	if register {
		e.status |= StatusRegister
	} else {
		e.status &= ^StatusRegister

	}
	// todo
	// 后续触发注册事件

	fmt.Printf("%v uid %v exec register val %v\n",
		reflect.TypeOf(e).String(), e.GetUid(), register)
}

func (e *BaseEntity) IsComponent() bool {
	return 0 != e.status&StatusComponent
}
func (e *BaseEntity) SetComponent(component bool) {
	if component {
		e.status |= StatusComponent
	} else {
		e.status &= ^StatusComponent
	}
}

func (e *BaseEntity) IsCreate() bool {
	return 0 != e.status&StatusCreate
}
func (e *BaseEntity) SetCreate(create bool) {
	if create {
		e.status |= StatusCreate
	} else {
		e.status &= ^StatusCreate
	}
}

func (e *BaseEntity) GetParent() *BaseEntity { return e.parent }
func (e *BaseEntity) SetParent(parent *BaseEntity) {
	if nil == parent || nil == parent.GetDomain() || e.GetParent() == parent {
		return
	}
	if nil != e.GetParent() {
		e.GetParent().RemoveFromChildren(e)
	}
	e.parent = parent
	e.SetComponent(false)
	e.GetParent().AddToChildren(e)
	e.SetDomain(parent.GetDomain())
}
func (e *BaseEntity) setComponentParent(parent *BaseEntity) {
	if nil == parent || nil == parent.GetDomain() || e.GetParent() == parent {
		return
	}
	if nil != e.GetParent() {
		e.GetParent().RemoveFromChildren(e)
	}
	e.parent = parent
	e.SetComponent(true)
	e.GetParent().AddToComponents(e)
	e.SetDomain(parent.GetDomain())
}
func (e *BaseEntity) AddToChildren(child *BaseEntity) {
	if nil == e.children {
		e.children = map[uint64]*BaseEntity{}
	}
	e.children[child.GetUid()] = child
}
func (e *BaseEntity) RemoveFromChildren(child *BaseEntity) {
	if nil == child {
		return
	}
	delete(e.children, child.GetUid())
	if 0 == len(e.children) {
		e.children = nil
	}
}

func (e *BaseEntity) AddToComponents(component *BaseEntity) {
	if nil == component {
		return
	}
	if nil == e.components {
		e.components = map[types.EntityType]*BaseEntity{}
	}
	e.components[component.impl.EntityTypeId()] = component
}
func (e *BaseEntity) RemoveFromComponents(component *BaseEntity) {
	if nil == component {
		return
	}
	delete(e.components, component.impl.EntityTypeId())
}

func (e *BaseEntity) GetDomain() *BaseEntity { return e.domain }
func (e *BaseEntity) SetDomain(domain *BaseEntity) {
	if nil == domain {
		return
	}
	if e.GetDomain() == domain {
		return
	}
	preDomain := e.domain
	e.domain = domain
	if nil == preDomain {
		// todo 生成uid
		e.uid = genUid()

		e.SetRegister(true)
	}
	for _, v := range e.children {
		v.SetDomain(e.GetDomain())
	}
	for _, v := range e.components {
		v.SetDomain(e.GetDomain())
	}
	if !e.IsCreate() {
		e.SetCreate(true)
	}
}

func (e *BaseEntity) IsDisposed() bool { return 0 == e.GetUid() }
