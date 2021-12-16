package skills

type SkillModel struct {
	Id uint32

	// 条件以及消耗先设计成函数便于扩展
	// 比如传送技能 普通传送+小飞鞋

	Condition func() bool

	Cost func()

	Effect *EffectModel
}

type Skill struct {
	Model *SkillModel

	Level uint32

	CoolDown int64
}
