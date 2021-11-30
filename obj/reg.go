package obj

import (
	"fmt"

	"github.com/Sogues/ETForGo/types"
)

var (
	global = &RegSys{}
)

type (
	RegSys struct {
		// componentId -- systemId
		systems map[types.EntityType]map[types.EntityType]BaseSystem
	}
)

func Reg(system BaseSystem) {
	if nil == global.systems {
		global.systems = map[types.EntityType]map[types.EntityType]BaseSystem{}
	}
	ctd := system.ComponentTypeId()
	group := global.systems[ctd]
	if nil == group {
		group = map[types.EntityType]BaseSystem{}
		global.systems[ctd] = group
	}
	std := system.SystemTypeId()
	if _, ok := group[std]; ok {
		panic(fmt.Sprintf("[Reg(system BaseSystem)] ctd %v std %v", ctd, std))
	}
	group[std] = system
}
