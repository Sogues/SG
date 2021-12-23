package battle

import "fmt"

var scene = &Group{}

type Group struct {
	timelines   []*TimelineObj // 每个人只能存在一调timeline
	damagesInfo []*DamageInfo
	objs        []GameObject
}

func (g *Group) GameObjectTick(interval float64) {
	for _, v := range g.objs {
		v.GetChaState().Tick(interval)
	}
}

func (g *Group) AddTimeline(timeline *TimelineObj) {
	if nil != timeline.caster && g.CasterHasTimeline(timeline.caster) {
		return
	}
	g.timelines = append(g.timelines, timeline)
}

func (g *Group) CasterHasTimeline(caster GameObject) bool {
	for _, v := range g.timelines {
		if caster == v.caster {
			return true
		}
	}
	return false
}

func (g *Group) TimelineTick(interval float64) {
	ori := len(g.timelines)

	if ori <= 0 {
		return
	}

	for idx := 0; idx < ori; {
		wasTimeElapsed := g.timelines[idx].timeElapsed
		g.timelines[idx].timeElapsed += interval * g.timelines[idx].timeScale
		if g.timelines[idx].model.chargeGoBack.atDuration < g.timelines[idx].timeScale &&
			g.timelines[idx].model.chargeGoBack.atDuration >= wasTimeElapsed {

			// todo

		}

		for i := range g.timelines[idx].model.nodes {
			if g.timelines[idx].model.nodes[i].timeElapsed < g.timelines[idx].timeElapsed &&
				g.timelines[idx].model.nodes[i].timeElapsed >= wasTimeElapsed {
				g.timelines[idx].model.nodes[i].doEvent(
					g.timelines[idx],
					g.timelines[idx].model.nodes[i].eveParams,
				)
			}
		}
		if g.timelines[idx].model.duration <= g.timelines[idx].timeElapsed {
			g.timelines = append(g.timelines[:idx], g.timelines[idx+1:]...)
			ori--
		} else {
			idx++
		}
	}
}

func (g *Group) DamageInfoTick(interval float64) {
	ori := len(g.damagesInfo)
	if 0 == ori {
		return
	}
	for _, v := range g.damagesInfo {
		g.DealWithDamage(v)
	}
	g.damagesInfo = nil
}

func (g *Group) DealWithDamage(dInfo *DamageInfo) {
	if nil == dInfo.defender {
		return
	}
	defenderChaState := dInfo.defender.GetChaState()
	if nil == defenderChaState {
		return
	}

	var attackerChaState *ChaState

	if defenderChaState.dead {
		return
	}
	if nil != dInfo.attacker {
		attackerChaState = dInfo.attacker.GetChaState()
		if nil != attackerChaState {
			for i := range attackerChaState.buffs {
				if nil != attackerChaState.buffs[i].model.onHit {
					attackerChaState.buffs[i].model.onHit(
						attackerChaState.buffs[i], dInfo, dInfo.defender,
					)
				}
			}
		}
	}

	for i := range defenderChaState.buffs {
		if nil != defenderChaState.buffs[i].model.onBeHurt {
			defenderChaState.buffs[i].model.onBeHurt(
				defenderChaState.buffs[i], dInfo, dInfo.attacker,
			)
		}
	}
	if defenderChaState.CanBeKilledByDamageInfo(dInfo) {
		if nil != attackerChaState {
			for i := range attackerChaState.buffs {
				if nil != attackerChaState.buffs[i].model.onKill {
					attackerChaState.buffs[i].model.onKill(
						attackerChaState.buffs[i], dInfo, dInfo.defender,
					)
				}
			}
		}
		for i := range defenderChaState.buffs {
			if nil != defenderChaState.buffs[i].model.onKill {
				defenderChaState.buffs[i].model.onKill(
					defenderChaState.buffs[i], dInfo, dInfo.attacker,
				)
			}
		}
	}

	isHeal := dInfo.IsHeal()
	dval := dInfo.DamageValue(isHeal)
	if isHeal || defenderChaState.immuneTime <= 0 {
		defenderChaState.ModResource(&ChaResource{
			hp:      -dval,
			ammo:    0,
			stamina: 0,
		})

		fmt.Printf("hp modify %v\n", dval)
	}
	for i := range dInfo.addBuffs {
		toCha := dInfo.addBuffs[i].target
		var toChaState *ChaState
		if toCha == dInfo.attacker {
			toChaState = attackerChaState
		} else {
			toChaState = defenderChaState
		}
		if nil != toChaState && !toChaState.dead {
			toChaState.AddBuff(dInfo.addBuffs[i])
		}
	}
}
