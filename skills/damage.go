package skills

type DamageTag int8

const (
	DamageTagDirectDamage = 1 << iota
	DamageTagPeriodDamage
	DamageTagReflectDamage
	DamageTagDirectHeal
	DamageTagPeriodHeal
)

// DamageStream 伤害流
type DamageStream struct {
	Attacker Player
	Defender Player
	Tag      DamageTag
	Damage   *Damage
}

// Damage 先定义个结构对象吧
type Damage struct {
	// 扣血 物理伤害 魔法伤害等等
	// 窃魔
	// 窃能量
}
