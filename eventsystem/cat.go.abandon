package eventsystem

type (
	componentTemp interface {
		ComponentTypeId() ComponentType
	}
	SysI interface {
		SystemTypeId() SystemType
		ComponentTypeId() ComponentType
		Run(temp componentTemp)
	}
	AwakeImpl interface {
		ComponentTypeId() ComponentType
		Awake(temp componentTemp)
	}
	AwakeBase struct {
		impl AwakeImpl
	}
)

func (AwakeBase) SystemTypeId() SystemType {
	return 1
}

func (a *AwakeBase) ComponentTypeId() ComponentType {
	return a.impl.ComponentTypeId()
}
func (a *AwakeBase) Run(temp componentTemp) {
	a.impl.Awake(temp)
}

type SpecifyComponent struct {
	id uint32
}

func (SpecifyComponent) ComponentTypeId() ComponentType { return 1 }

type SpecifyComponentSystem struct {
	AwakeBase
}

func NewSpecifyComponentSystem() SysI {
	s := &SpecifyComponentSystem{}
	s.impl = s
	return s
}

func (SpecifyComponentSystem) ComponentTypeId() ComponentType {
	return SpecifyComponent{}.ComponentTypeId()
}

func (s *SpecifyComponentSystem) Awake(temp componentTemp) {
	cpt, ok := temp.(*SpecifyComponent)
	if !ok {
		return
	}
	_ = cpt
}
