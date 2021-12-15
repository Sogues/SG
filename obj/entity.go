package obj

import (
	"fmt"

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

func RegEntity(entity Entity) {
	if nil == EntityFactory.types {
		EntityFactory.types = map[types.EntityType]func() Entity{}
	}
	_, ok := EntityFactory.types[entity.EntityTypeId()]
	if ok {
		panic(fmt.Sprintf("RegEntity duplicate type id %v", entity.EntityTypeId()))
	}
	EntityFactory.types[entity.EntityTypeId()] = entity.New
}

var (
	EntityFactory = &Factory{}
)

type (
	// todo 管理entity自身的类型id

	BaseEntity struct {
		instanceId uint64
		uid        uint64
		status     EntityStatus

		domain Entity // 作用域

		parent Entity // 父节点

		children map[uint64]Entity

		// 组件以类型为单位 禁止同类型组件重复添加
		components map[types.EntityType]Entity
	}

	Factory struct {
		types map[types.EntityType]func() Entity
	}
)

func (e *BaseEntity) GetInstanceId() uint64           { return e.instanceId }
func (e *BaseEntity) SetInstanceId(instanceId uint64) { e.instanceId = instanceId }
func (e *BaseEntity) GetUid() uint64                  { return e.uid }
func (e *BaseEntity) SetUid(uid uint64)               { e.uid = uid }
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
func (e *BaseEntity) SetRegister(self Entity, register bool) {
	if !e.checkSelf(self) {
		return
	}
	if register == e.IsRegister() {
		return
	}
	if register {
		e.status |= StatusRegister
	} else {
		e.status &= ^StatusRegister

	}
	// 后续触发注册事件
	SystemProcessor.Register(self, register)
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
func (e *BaseEntity) setParent(self, parent Entity) {
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
	e.SetDomain(self, parent.GetDomain())
}
func (e *BaseEntity) GetComponentParent() Entity { return e.parent }
func (e *BaseEntity) setComponentParent(self, parent Entity) {
	// 必须是自己
	if !e.checkSelf(self) {
		return
	}
	if nil == parent || nil == parent.GetDomain() || e.GetComponentParent() == parent {
		return
	}
	if nil != e.GetComponentParent() {
		e.GetComponentParent().removeFromChildren(self)
	}
	e.parent = parent
	e.SetComponent(true)
	e.GetParent().addToComponents(self)
	e.SetDomain(self, parent.GetDomain())
}

func (e *BaseEntity) GetDomain() Entity { return e.domain }
func (e *BaseEntity) SetDomain(self, domain Entity) {
	if !e.checkSelf(self) || nil == domain {
		return
	}
	if e.GetDomain() == domain {
		return
	}
	preDomain := e.domain
	e.domain = domain
	if nil == preDomain {
		e.instanceId = IdGen.GenId()

		e.SetRegister(self, true)
	}
	for _, v := range e.children {
		v.SetDomain(v, e.GetDomain())
	}
	for _, v := range e.components {
		v.SetDomain(v, e.GetDomain())
	}
	if !e.IsCreate() {
		e.SetCreate(true)
	}
}

func (e *BaseEntity) IsDisposed() bool { return 0 == e.GetInstanceId() }

func (e *BaseEntity) GetComponent(entityType types.EntityType) Entity {
	v, ok := e.components[entityType]
	if !ok {
		return nil
	}
	return v
}
func (e *BaseEntity) AddComponent(self Entity, entityType types.EntityType, param interface{}) Entity {
	// 必须是自己
	if !e.checkSelf(self) {
		return nil
	}
	if _, ok := e.components[entityType]; ok {
		// todo
		return nil
	}
	et := e.create(entityType)
	if nil == et {
		return nil
	}
	// 不一致点
	et.SetUid(self.GetUid())
	et.setComponentParent(et, self)
	SystemProcessor.Awake(et, param)
	return et
}

func (e *BaseEntity) AddChild(self Entity, entityType types.EntityType, param interface{}) Entity {
	return e.AddChildWithId(self, IdGen.GenId(), entityType, param)
}

func (e *BaseEntity) AddChildWithId(self Entity, id uint64, entityType types.EntityType, param interface{}) Entity {
	// 必须是自己
	if !e.checkSelf(self) {
		return nil
	}
	if _, ok := e.components[entityType]; ok {
		// todo
		return nil
	}
	et := e.create(entityType)
	if nil == et {
		return nil
	}
	// 不一致点
	et.SetUid(id)
	et.setParent(et, self)
	SystemProcessor.Awake(et, param)
	return et
}

func (e *BaseEntity) Dispose(self Entity) {
	if !e.checkSelf(self) {
		return
	}
	if e.IsDisposed() {
		return
	}
	e.SetRegister(self, false)
	e.SetInstanceId(0)
	for _, v := range e.components {
		v.Dispose(v)
	}
	e.components = nil
	for _, v := range e.children {
		v.Dispose(v)
	}
	e.children = nil
	SystemProcessor.Destroy(self)
	e.domain = nil
	if nil != e.GetParent() && !e.GetParent().IsDisposed() {
		if e.IsComponent() {
			e.GetComponentParent().removeFromComponents(self)
		} else {
			e.GetParent().removeFromChildren(self)
		}
	}
	e.status = 0
}

func (*BaseEntity) create(entityType types.EntityType) Entity {
	et := EntityFactory.Create(entityType)
	if nil == et {
		return nil
	}
	et.SetCreate(true)
	et.SetUid(0)
	return et
}

func (e *BaseEntity) getBase() *BaseEntity {
	if nil == e {
		return nil
	}
	return e
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
	if nil == e.components {
		e.components = nil
	}
}

func (e *BaseEntity) checkSelf(self Entity) bool {
	if nil == e || nil == self.getBase() {
		return false
	}
	ok := e == self.getBase()
	return ok
}

func (f *Factory) Create(eid types.EntityType) Entity {
	if nil == f.types {
		return nil
	}
	fn := f.types[eid]
	if nil == fn {
		return nil
	}
	return fn()
}

func (f *Factory) Release(entity Entity) {
	return
}
