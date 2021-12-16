package skills

var (
	// EffectorGroup 随便写个组
	EffectorGroup []*Effect
)

// Player 谁都可以是对象
type Player interface {
	// DoNothing 我就是啥也不干
	DoNothing()
	CharacterControlState() CharacterControlState
	GetSkill(id uint32) *Skill
	GetBuffs() []*Buff
}

func hasMask(val, mask CharacterControlState) bool {
	ok := val & mask
	return 0 != ok
}

func CastSkill(skillId uint32, p Player) int32 {
	if !hasMask(p.CharacterControlState(), CharacterControlStateCastSkill) {
		return -1
	}
	skill := p.GetSkill(skillId)
	if nil == skill || 0 != skill.CoolDown {
		return -2
	}
	if nil != skill.Model.Condition {
		if !skill.Model.Condition() {
			return -3
		}
	}
	effect := &Effect{
		Model:     skill.Model.Effect,
		Caster:    p,
		Skill:     skill,
		TimeScale: 0,
		Elapsed:   0,
		Params:    nil,
	}

	for _, v := range p.GetBuffs() {
		if nil != v.Model.OnSkillCast {
			// 现有buff可能会对技能效果进行修改
			effect = v.Model.OnSkillCast(v, skill, effect)
		}
	}
	// 考虑同步问题 再说吧 todo
	// gcd 100ms
	skill.CoolDown = 100
	if nil == effect {
		return -4
	}
	if nil != skill.Model.Cost {
		skill.Model.Cost()
	}
	EffectorGroup = append(EffectorGroup, effect)
	return 0
}

func UpdateEffector(interval int64) {
	if 0 == len(EffectorGroup) {
		return
	}
	for i, effect := range EffectorGroup {
		lastElapsed := effect.Elapsed
		effect.Elapsed += interval
		for _, node := range effect.Model.Nodes {
			// 在这区间的
			if node.Elapsed < effect.Elapsed &&
				node.Elapsed >= lastElapsed {
				node.Handle(effect, node.Params)
			}
		}
		if effect.Model.Duration <= effect.Elapsed {
			EffectorGroup = append(EffectorGroup[:i], EffectorGroup[i+1:]...)
		}
	}
}
