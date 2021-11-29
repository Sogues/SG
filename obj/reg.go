package obj

import (
	"fmt"
)

var (
	global = &RegSys{}
)

type (
	RegSys struct {
		systems map[ComponentType]map[SystemType]BaseSystem
	}
)

func Reg(system BaseSystem) {
	if nil == global.systems {
		global.systems = map[ComponentType]map[SystemType]BaseSystem{}
	}
	ctd := system.ComponentTypeId()
	group := global.systems[ctd]
	if nil == group {
		group = map[SystemType]BaseSystem{}
		global.systems[ctd] = group
	}
	std := system.SystemTypeId()
	if _, ok := group[std]; ok {
		panic(fmt.Sprintf("[Reg(system BaseSystem)] ctd %v std %v", ctd, std))
	}
	group[std] = system
}
