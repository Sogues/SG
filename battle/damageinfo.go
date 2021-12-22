package battle

type (
	Damage struct {
		bullet    int
		explosion int
		mental    int
	}

	DamageInfo struct {
		attacker GameObject

		defender GameObject

		tags []DamageInfoTag

		damage Damage

		criticalRate float64

		hitRate float64

		//伤害的角度，作为伤害打向角色的入射角度，比如子弹，就是它当前的飞行角度
		//潘森抵抗正面的伤害
		degree float64

		addBuffs []*AddBuffInfo
	}
)

func (d *DamageInfo) IsHeal() bool {
	for _, v := range d.tags {
		if v == directHeal || v == periodHeal {
			return true
		}
	}
	return false
}

func (d *DamageInfo) DamageValue(asHeal bool) int {
	// todo 只计算一次 每次计算后加入缓存
	// 后续找个数值给点方案
	return 10
}
