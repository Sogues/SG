package battle

type MoveType int8

const (
	ground MoveType = iota
	fly
)

type DamageInfoTag int32

const (
	directDamage  DamageInfoTag = 0  //直接伤害
	periodDamage  DamageInfoTag = 1  //间歇性伤害
	reflectDamage DamageInfoTag = 2  //反噬伤害
	directHeal    DamageInfoTag = 10 //直接治疗
	periodHeal    DamageInfoTag = 11 //间歇性治疗
)
