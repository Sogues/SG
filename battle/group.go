package battle

type Group struct {
	timelines   []*TimelineObj
	damagesInfo []*DamageInfo
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
					attackerChaState.buffs[i].model.onHit(attackerChaState.buffs[i], dInfo, dInfo.defender)
				}
			}
		}
	}

	for i := range defenderChaState.buffs {
		if nil != defenderChaState.buffs[i].model.onBeHurt {
			defenderChaState.buffs[i].model.onBeHurt(defenderChaState.buffs[i], dInfo, dInfo.attacker)
		}
	}
	//defenderChaState.CanBeKilledByDamageInfo()
}
