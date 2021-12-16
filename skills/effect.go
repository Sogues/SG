package skills

type EffectNode struct {
	// 时间点
	Elapsed int64

	// 效果逻辑
	Handle func()
}

type EffectModel struct {
	Nodes []*EffectNode
	// 总计时长 Nodes[-1].Elapsed <= Duration
	Duration int64
}

type Effect struct {
	Model *EffectModel

	Caster Player

	// 逻辑倍速 暴走
	TimeScale int64

	// 流逝时间
	Elapsed int64

	// 参数相关 先有这个东西 怎么实现再说吧
	//TODO 是否准备个struct方便策划 或者键值对直接使用 但是键值对无论性能还是内存都很拉
	Params interface{}
}
