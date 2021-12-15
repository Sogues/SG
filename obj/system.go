package obj

import (
	"fmt"

	"github.com/Sogues/ETForGo/types"
)

// 通过 component.EntityTypeId()
// 获取属于该component的所有 system
// 通过 systems[BaseSystemId] 获取 system
// 执行对应函数

var (
	SystemProcessor = NewProcessor()
)

func RegSystem(sys sysInter) {
	if nil == SystemProcessor.types {
		SystemProcessor.types = map[types.EntityType][]sysInter{}
	}
	sid, eid := sys.BaseSystemId(), sys.EntityTypeId()
	l := SystemProcessor.types[eid]
	for _, v := range l {
		if v.BaseSystemId() == sid {
			panic(fmt.Sprintf("same EntityId %v with Sys Id %v", eid, sid))
		}
	}
	SystemProcessor.types[eid] = append(l, sys)
}

type (
	BaseAwakeSystem struct{}
)

func (*BaseAwakeSystem) BaseSystemId() types.EntityType { return types.EntityTypeSystemAwake }

func NewProcessor() *Processor {
	p := &Processor{}
	return p
}

type (
	Processor struct {
		types map[types.EntityType][]sysInter

		entities map[uint64]Entity

		updates map[uint64]struct{}
	}
)

func (p *Processor) Register(entity Entity, register bool) {
	uid := entity.GetUid()

	if register {
		if nil == p.entities {
			p.entities = map[uint64]Entity{}
		}
		p.entities[uid] = entity
		l, ok := p.types[entity.EntityTypeId()]
		if !ok {
			return
		}
		for _, v := range l {
			if v.BaseSystemId() == types.EntityTypeSystemUpdate {
				if nil == p.updates {
					p.updates = map[uint64]struct{}{}
				}
				p.updates[uid] = struct{}{}
				break
			}
		}
	} else {
		delete(p.entities, uid)
	}
}

func (p *Processor) Awake(entity Entity, param interface{}) {
	l, ok := p.types[entity.EntityTypeId()]
	if !ok {
		return
	}
	for _, v := range l {
		if v.BaseSystemId() != types.EntityTypeSystemAwake {
			continue
		}
		v.(awakeSysInter).Awake(entity, param)
	}
}
