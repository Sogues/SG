package battle

type (
	BuffOnOccur    func(buff *BuffObj, modifyStack int)
	BuffOnRemoved  func(buff *BuffObj)
	BuffOnTick     func(buff *BuffObj)
	BuffOnHit      func(buff *BuffObj, damageInfo *DamageInfo, target GameObject)
	BuffOnBeHurt   func(buff *BuffObj, damageInfo *DamageInfo, attacker GameObject)
	BuffOnKill     func(buff *BuffObj, damageInfo *DamageInfo, target GameObject)
	BuffOnBeKilled func(buff *BuffObj, damageInfo *DamageInfo, attacker GameObject)
	BuffOnCast     func(buff *BuffObj, skill *SkillObj, timeline *TimelineObj)

	BuffModel struct {
		id string

		name string

		priority int

		maxStack int

		tags []string

		tickTime float64

		propMod []ChaProperty

		stateMod ChaControlState

		onOccur       BuffOnOccur
		onOccurParams []object

		onTick       BuffOnTick
		onTickParams []object

		onRemoved       BuffOnRemoved
		onRemovedParams []object

		onCast       BuffOnCast
		onCastParams []object

		onHit       BuffOnHit
		onHitParams []object

		onBeHurt       BuffOnBeHurt
		onBeHurtParams []object

		onKill       BuffOnKill
		onKillParams []object

		onBeKilled       BuffOnBeKilled
		onBeKilledParams []object
	}

	BuffObj struct {
		model BuffModel

		duration float64

		permanent bool

		stack int

		caster GameObject

		carrier GameObject

		timeElapsed float64

		ticked int

		buffParam map[string]object
	}

	AddBuffInfo struct {
		caster GameObject

		target GameObject

		buffModel BuffModel

		//要添加的层数，负数则为减少
		addStack int

		//关于时间，是改变还是设置为, true代表设置为，false代表改变
		durationSetTo bool

		//是否是一个永久的buff，即便=true，时间设置也是有意义的，因为时间如果被减少到0以下，即使是永久的也会被删除
		permanent bool

		//时间值，设置为这个值，或者加上这个值，单位：秒
		duration float64

		buffParam map[string]object
	}
)
