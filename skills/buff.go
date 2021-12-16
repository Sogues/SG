package skills

type (
	// OnOccur 创建或者层数调整时触发
	OnOccur func(buff *Buff, modifyStack int32)
	// OnRemoved 层数为0时移除前调用
	OnRemoved func(buff *Buff)
	// OnTick tick帧的回调点
	OnTick func(buff *Buff)
	// OnSkillCast 释放技能时回调
	OnSkillCast func(buff *Buff, skill *Skill, effect *Effect)
	// OnHit 命中时候的回调
	OnHit func(buff *Buff, stream *DamageStream, target Player)
	// OnBeHurt 被命中时候 受到攻击时候 受到伤害时候
	OnBeHurt func(buff *Buff, stream *DamageStream, attacker Player)
	// OnKill 确定造成击杀后
	OnKill func(buff *Buff, stream *DamageStream, target Player)
	// OnBeKilled 确定死亡后
	OnBeKilled func(buff *Buff, stream *DamageStream, attacker Player)
)

type BuffModel struct {
	// 数据id
	Id uint32

	// 优先级
	Priority uint32

	// 最大层数
	MaxStack uint32

	// tick间隔 单位毫秒
	TickInterval int32

	// 对于持有者的状态影响
	CarrierState CharacterControlState

	// todo 考虑buff本身对属性的基本修改 不要每次都走回调？

	// 不走接口因为需要方便buff效果的自由组装

	OnOccur OnOccur

	OnRemoved OnRemoved

	OnTick OnTick

	OnSkillCast OnSkillCast

	OnHit OnHit

	OnBeHurt OnBeHurt

	OnKill OnKill

	OnBeKilled OnBeKilled
}

type Buff struct {
	Model *BuffModel

	// 剩余毫秒
	Duration int64

	// 流逝
	Elapsed int64

	// OnTick 触发次数
	TickTimes uint32

	// 永久  TODO 是否通过Duration=-1 省个字段
	Permanent bool

	Stack int32

	Caster Player

	Carrier Player

	// 参数相关 先有这个东西 怎么实现再说吧
	//TODO 是否准备个struct方便策划 或者键值对直接使用 但是键值对无论性能还是内存都很拉
	Params interface{}
}
