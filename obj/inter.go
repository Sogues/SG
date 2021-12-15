package obj

import "github.com/Sogues/ETForGo/types"

type (
	Entity interface {
		New() Entity
		EntityTypeId() types.EntityType
		GetUid() uint64
		SetUid(uid uint64)
		FromPool() bool
		SetFromPool(from bool)
		IsRegister() bool
		SetRegister(self Entity, register bool)
		IsComponent() bool
		SetComponent(component bool)
		IsCreate() bool
		SetCreate(create bool)
		GetParent() Entity
		// 传入self是由于base并不是完整的entity
		SetParent(self, parent Entity)
		SetComponentParent(self, parent Entity)
		GetDomain() Entity
		SetDomain(self, domain Entity)
		IsDisposed() bool

		// 方便在base类中调用
		getBase() *BaseEntity
		addToChildren(child Entity)
		removeFromChildren(child Entity)
		addToComponents(component Entity)
		removeFromComponents(component Entity)
	}

	sysInter interface {
		EntityTypeId() types.EntityType
		BaseSystemId() types.EntityType
	}
	awakeSysInter interface {
		sysInter
		Awake(component, param interface{})
	}
)
