package skills

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
	CarrierState int32

	// todo 考虑buff本身对属性的基本修改 不要每次都走回调？

	// 不走接口因为需要方便buff效果的自由组装

	// OnOccur 创建或者层数调整时触发
	OnOccur func()
	// OnRemoved 层数为0时移除前调用
	OnRemoved func()

	// OnTick tick帧的回调点
	OnTick func()

	// OnSkillCast 释放技能时回调
	OnSkillCast func()

	// OnHit 命中时候的回调
	OnHit func()
	// OnBeHurt 被命中时候 受到攻击时候 受到伤害时候
	OnBeHurt func()

	// OnKill 确定造成击杀后
	OnKill func()
	// OnBeKilled 确定死亡后
	OnBeKilled func()
}
