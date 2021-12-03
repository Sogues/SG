package es

import (
	"github.com/Sogues/ETForGo/types"
)

// 通过 component.EntityTypeId()
// 获取属于该component的所有 system
// 通过 systems[BaseSystemId] 获取 system
// 执行对应函数

type (
	Type interface {
		EntityTypeId() types.EntityType
		BaseSystemId() types.EntityType
	}
	AwakeEs interface {
		Awake(component, param interface{})
	}
	BaseAwakeSystem struct {
	}
)

func (*BaseAwakeSystem) BaseSystemId() types.EntityType { return types.EntityTypeSystemAwake }

func (*BaseAwakeSystem) Awake(component, param interface{}) {
}
