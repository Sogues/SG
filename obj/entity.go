package obj

import (
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
	// todo 管理entity自身的类型id
	Entity interface {
		EntityTypeId() types.EntityType
		GetUid() uint64
		SetUid(uid uint64)
		FromPool() bool
		SetFromPool(from bool)
		IsRegister() bool
		SetRegister(register bool)
		IsComponent() bool
		SetComponent(component bool)
		IsCreate() bool
		SetCreate(create bool)
		GetParent() Entity
		SetParent(self, parent Entity)
		SetComponentParent(self, parent Entity)
		GetDomain() Entity
		SetDomain(domain Entity)
		IsDisposed() bool

		//
		addToChildren(child Entity)
		removeFromChildren(child Entity)
		addToComponents(component Entity)
		removeFromComponents(component Entity)

		// base部分驱动
		getBase() *BaseEntity
	}
	BaseEntity struct {
		uid    uint64
		status EntityStatus

		domain Entity // 作用域

		parent Entity // 父节点

		children map[uint64]Entity

		// 组件以类型为单位 禁止同类型组件重复添加
		components map[types.EntityType]Entity
	}
)

func (e *BaseEntity) getBase() *BaseEntity {
	if nil == e {
		return nil
	}
	return e
}

func (e *BaseEntity) checkSelf(self Entity) bool {
	if nil == e || nil == self.getBase() {
		return false
	}
	ok := e == self.getBase()
	return ok
}

func (e *BaseEntity) GetUid() uint64    { return e.uid }
func (e *BaseEntity) SetUid(uid uint64) { e.uid = uid }
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

func (e *BaseEntity) GetParent() Entity { return e.parent }
func (e *BaseEntity) SetParent(self, parent Entity) {
	// 必须是自己
	if !e.checkSelf(self) {
		return
	}
	if nil == parent || nil == parent.GetDomain() || e.GetParent() == parent {
		return
	}
	if nil != e.GetParent() {
		e.GetParent().removeFromChildren(self)
	}
	e.parent = parent
	e.SetComponent(false)
	e.GetParent().addToChildren(self)
	e.SetDomain(parent.GetDomain())
}
func (e *BaseEntity) SetComponentParent(self, parent Entity) {
	// 必须是自己
	if !e.checkSelf(self) {
		return
	}
	if nil == parent || nil == parent.GetDomain() || e.GetParent() == parent {
		return
	}
	if nil != e.GetParent() {
		e.GetParent().removeFromChildren(self)
	}
	e.parent = parent
	e.SetComponent(true)
	e.GetParent().addToComponents(self)
	e.SetDomain(parent.GetDomain())
}
func (e *BaseEntity) addToChildren(child Entity) {
	if nil == e.children {
		e.children = map[uint64]Entity{}
	}
	e.children[child.GetUid()] = child
}
func (e *BaseEntity) removeFromChildren(child Entity) {
	if nil == child {
		return
	}
	delete(e.children, child.GetUid())
	if 0 == len(e.children) {
		e.children = nil
	}
}

func (e *BaseEntity) addToComponents(component Entity) {
	if nil == component {
		return
	}
	if nil == e.components {
		e.components = map[types.EntityType]Entity{}
	}
	e.components[component.EntityTypeId()] = component
}
func (e *BaseEntity) removeFromComponents(component Entity) {
	if nil == component {
		return
	}
	delete(e.components, component.EntityTypeId())
}

func (e *BaseEntity) GetDomain() Entity { return e.domain }
func (e *BaseEntity) SetDomain(domain Entity) {
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
