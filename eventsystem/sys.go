package eventsystem

const (
	SystemTypeNone SystemType = iota
	SystemTypeAwake
)

var (
	ComponentTypeBaseId ComponentType = 0
)

func GenComponentTypeId() ComponentType {
	ComponentTypeBaseId++
	return ComponentTypeBaseId
}

type (
	SystemType    uint32
	ComponentType uint32
)

type (
	BaseSystem interface {
		ComponentTypeId() ComponentType
		SystemTypeId() SystemType
		Run(component, param interface{})
	}

	AwakeBaseSystemImpl interface {
		ComponentTypeId() ComponentType
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

func (s *AwakeBaseSystem) ComponentTypeId() ComponentType {
	return s.impl.ComponentTypeId()
}
func (s *AwakeBaseSystem) SystemTypeId() SystemType {
	return SystemTypeAwake
}

func (s *AwakeBaseSystem) Run(component, param interface{}) {
	s.impl.Awake(component, param)
}
