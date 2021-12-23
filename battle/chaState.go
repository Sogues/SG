package battle

import (
	"fmt"
)

func NewChaState(obj GameObject) *ChaState {
	c := &ChaState{}
	c.GameObject = obj
	c._controlState = ChaControlState{
		canMove:     true,
		canRotate:   true,
		canUseSkill: true,
	}
	c.timelineControlState = ChaControlState{
		canMove:     true,
		canRotate:   true,
		canUseSkill: true,
	}
	c.resource = &ChaResource{
		hp: 100,
	}
	c.moveSpeed = 3
	c.actionSpeed = 1
	c.baseProp = ChaProperty{
		hp:          100,
		attack:      20,
		moveSpeed:   3,
		actionSpeed: 1,
		ammo:        0,
		bodyRadius:  0.25,
		hitRadius:   0.25,
		moveType:    ground,
	}
	c.buffProp = []ChaProperty{{}, {}}
	return c
}

type (
	ChaState struct {
		GameObject

		_controlState ChaControlState

		timelineControlState ChaControlState

		//角色的无敌状态持续时间，如果在无敌状态中，子弹不会碰撞，DamageInfo处理无效化
		immuneTime float64

		charging bool

		moveDegree float64

		faceDegree float64

		dead bool

		moveOrder Vector3

		forceMove []*MovePreorder

		aimOrder []string

		rotateToOrder float64

		forceRotate []float64

		resource *ChaResource

		side int

		tags []string

		property ChaProperty

		moveSpeed float64

		actionSpeed float64

		baseProp ChaProperty

		buffProp []ChaProperty

		equipmentProp ChaProperty

		skills []*SkillObj

		buffs []*BuffObj
	}
)

func (c *ChaState) Tick(interval float64) {
	if c.dead {
		return
	}

	if c.immuneTime > 0 {
		c.immuneTime -= interval
		if c.immuneTime < 0 {
			c.immuneTime = 0
		}
	}
	for i := range c.skills {
		if c.skills[i].coolDown > 0 {
			c.skills[i].coolDown -= interval
			if c.skills[i].coolDown < 0 {
				c.skills[i].coolDown = 0
			}
		}
	}
	for i := 0; i < len(c.buffs); {
		if !c.buffs[i].permanent {
			c.buffs[i].duration -= interval
		}
		c.buffs[i].timeElapsed += interval
		if c.buffs[i].model.tickTime > 0 &&
			nil != c.buffs[i].model.onTick {

			tick := 0 == c.buffs[i].tickTimeElapsed
			if c.buffs[i].tickTimeElapsed >= c.buffs[i].model.tickTime {
				c.buffs[i].tickTimeElapsed -= c.buffs[i].model.tickTime
				tick = true
			}

			if tick {
				c.buffs[i].model.onTick(
					c.buffs[i],
				)
				c.buffs[i].ticked += 1
			}
			c.buffs[i].tickTimeElapsed += interval

		}
		if c.buffs[i].duration <= 0 || c.buffs[i].stack <= 0 {
			if nil != c.buffs[i].model.onRemoved {
				c.buffs[i].model.onRemoved(
					c.buffs[i],
				)
			}
			c.buffs = append(c.buffs[:i], c.buffs[i+1:]...)
		} else {
			i++
		}
	}
	// 移动数据
}

func (c *ChaState) CanBeKilledByDamageInfo(damageInfo *DamageInfo) bool {
	if c.immuneTime > 0 || damageInfo.IsHeal() {
		return false
	}
	val := damageInfo.DamageValue(false)
	ok := val >= c.resource.hp
	return ok
}

func (c *ChaState) ModResource(value *ChaResource) {
	c.resource.hp += value.hp
	c.resource.ammo += value.ammo
	c.resource.stamina += value.stamina

	c.resource.hp = clamp(c.resource.hp, 0, c.baseProp.hp)
	c.resource.ammo = clamp(c.resource.ammo, 0, c.baseProp.ammo)
	c.resource.stamina = clamp(c.resource.stamina, 0, 100)

	if 0 == c.resource.hp {
		c.Kill()
	}
}

func (c *ChaState) Kill() {
	c.dead = true
	fmt.Printf("todo dead \n")
}

func (c *ChaState) AddBuff(buff AddBuffInfo) {
	hasOne := c.GetBuffId(buff.buffModel.id, buff.caster)
	modStack := buff.addStack
	if modStack > buff.buffModel.maxStack {
		modStack = buff.buffModel.maxStack
	}
	toRemove := false
	var toAddBuff *BuffObj
	if nil != hasOne {
		hasOne.buffParam = map[string]object{}
		for k, v := range buff.buffParam {
			hasOne.buffParam[k] = v
		}
		// 调整buff剩余持续时间
		if buff.durationSetTo {
			hasOne.duration = buff.duration
		} else {
			hasOne.duration += buff.duration
		}
		affAdd := hasOne.stack + modStack
		if affAdd > hasOne.model.maxStack {
			modStack = hasOne.model.maxStack - hasOne.stack
		} else if affAdd <= 0 {
			affAdd = -hasOne.stack
		}

		hasOne.stack += modStack
		hasOne.permanent = buff.permanent
		toAddBuff = hasOne
		toRemove = hasOne.stack <= 0
	} else {
		toAddBuff = NewBuffObj(
			buff.buffModel,
			buff.duration,
			buff.permanent,
			buff.addStack,
			buff.caster,
			c.GameObject,
			buff.buffParam,
		)
		// todo 调整优先级 刷新
		c.buffs = append(c.buffs, toAddBuff)
	}
	if !toRemove && nil != buff.buffModel.onOccur {
		buff.buffModel.onOccur(toAddBuff, modStack)
	}
	// todo 重新计算属性
}

func (c *ChaState) GetBuffId(id string, caster GameObject) *BuffObj {
	for i := range c.buffs {
		if id != c.buffs[i].model.id {
			continue
		}
		// 可以存在不同对象添加同一个类型的buff
		if nil == caster || caster == c.buffs[i].caster {
			return c.buffs[i]
		}
	}
	return nil
}

func (c *ChaState) CastSkill(id string) bool {
	if !c._controlState.canUseSkill {
		return false
	}
	skillObj := c.GetSkillById(id)
	if nil == skillObj || skillObj.coolDown > 0 {
		return false
	}
	castSuccess := false
	if c.resource.Enough(skillObj.model.condition) {
		timeline := NewTimelineObj(
			skillObj.model.effect, c, skillObj,
		)
		for i := range c.buffs {
			if nil != c.buffs[i].model.onCast {
				timeline = c.buffs[i].model.onCast(
					c.buffs[i], skillObj, timeline,
				)
			}
		}
		if nil != timeline {
			c.ModResource(&ChaResource{
				hp:      -skillObj.model.cost.hp,
				ammo:    -skillObj.model.cost.ammo,
				stamina: -skillObj.model.cost.stamina,
			})
			scene.AddTimeline(timeline)
			skillObj.coolDown = skillObj.coolDownStatic
			castSuccess = true
		}
	} else {
		skillObj.coolDown = 0.1
	}
	return castSuccess
}

func (c *ChaState) GetSkillById(id string) *SkillObj {
	for _, v := range c.skills {
		if id == v.model.id {
			return v
		}
	}
	return nil
}

func clamp(val int, a, b int) int {
	if val < a {
		return a
	} else if val > b {
		return b
	}
	return val
}
