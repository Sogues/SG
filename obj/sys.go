package obj

import "github.com/Sogues/ETForGo/types"

type (
	BaseSystem interface {
		ComponentTypeId() types.EntityType
		SystemTypeId() types.EntityType
		Run(component, param interface{})
	}

	AwakeBaseSystemImpl interface {
		ComponentTypeId() types.EntityType
		Awake(component, param interface{})
	}

	AwakeBaseSystem struct {
		impl AwakeBaseSystemImpl
	}
)

func NewAwakeSystemWithImpl(impl AwakeBaseSystemImpl) BaseSystem {
	s := &AwakeBaseSystem{}
	s.impl = impl
	return s
}

func (s *AwakeBaseSystem) ComponentTypeId() types.EntityType {
	return s.impl.ComponentTypeId()
}
func (s *AwakeBaseSystem) SystemTypeId() types.EntityType {
	return types.EntityTypeSystemAwake
}

func (s *AwakeBaseSystem) Run(component, param interface{}) {
	s.impl.Awake(component, param)
}
