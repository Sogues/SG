package obj

import "github.com/Sogues/ETForGo/types"

type (
	Entity interface {
		New() Entity
		EntityTypeId() types.EntityType
		GetInstanceId() uint64
		SetInstanceId(instanceId uint64)
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
		GetDomain() Entity
		SetDomain(self, domain Entity)
		IsDisposed() bool
		AddComponent(self Entity, entityType types.EntityType, param interface{}) Entity
		AddChild(self Entity, entityType types.EntityType, param interface{}) Entity
		AddChildWithId(self Entity, id uint64, entityType types.EntityType, param interface{}) Entity
		Dispose(self Entity)

		// 非对外
		// SetParent 传入self是由于base并不是完整的entity
		setParent(self, parent Entity)
		setComponentParent(self, parent Entity)
		addToChildren(child Entity)
		removeFromChildren(child Entity)
		addToComponents(component Entity)
		removeFromComponents(component Entity)

		// 方便在base类中调用
		getBase() *BaseEntity
	}

	sysInter interface {
		EntityTypeId() types.EntityType
		BaseSystemId() types.EntityType
	}
	awakeSysInter interface {
		sysInter
		Awake(component, param interface{})
	}
	destroyInter interface {
		sysInter
		Destroy(component Entity)
	}
)
