package battle

import (
	"fmt"
	"testing"
	"time"
)

func TestObject(t *testing.T) {
	var a []object

	a = append(a, 10)
	a = append(a, "xx")

	i, ok := a[0].(int)
	if !ok {
		panic("a[0].(int)")
	}
	fmt.Println(i)
	s, ok := a[1].(string)
	if !ok {
		panic("a[1].(string)")
	}
	fmt.Println(s)
}

func TestXXX(t *testing.T) {
	t.Run("x", func(t *testing.T) {
		type st struct {
			a int
		}
		var s []st
		s = append(s, struct{ a int }{a: 10})
		s = append(s, struct{ a int }{a: 11})
		s = append(s, struct{ a int }{a: 12})
		s[1].a = 20
		fmt.Println(s)
	})
}

type testGameObject struct {
	c *ChaState
}

func (t *testGameObject) GetChaState() *ChaState {
	return t.c
}

func TestChaState_CastSkill(t *testing.T) {
	tb := &testGameObject{}
	tb.c = NewChaState(tb)
	scene.objs = append(scene.objs, tb)

	timelineMode := NewTimelineMode(
		"技能1",
		[]TimelineNode{
			// 0s时候置状态
			NewTimelineNode(
				0,
				func(timeline *TimelineObj, args []object) {
					if nil == timeline.caster {
						return
					}
					cs := timeline.caster.GetChaState()
					if nil == cs {
						return
					}
					cs.timelineControlState.canMove = args[0].(bool)
					cs.timelineControlState.canRotate = args[0].(bool)
					cs.timelineControlState.canUseSkill = args[0].(bool)
					fmt.Printf("释放技能后0s执行node1\n")
				},
				[]object{true, true, false},
			),
			// 0.1s时候执行生效逻辑
			NewTimelineNode(
				0.1,
				func(timeline *TimelineObj, args []object) {
					if nil == timeline.caster {
						return
					}
					cs := timeline.caster.GetChaState()
					if nil == cs {
						return
					}
					// 添加一个buff吧
					buffModel := BuffModel{
						id:            "扣血",
						name:          "扣血",
						priority:      1,
						maxStack:      1,
						tags:          nil,
						tickTime:      0.1,
						propMod:       nil,
						stateMod:      ChaControlStateOrigin,
						onOccur:       nil,
						onOccurParams: nil,
						onTick: func(buff *BuffObj) {
							fmt.Printf("buff tick 扣血触发 抛出一个damgeinfo事件 \n")
							damageInfo := &DamageInfo{
								attacker: buff.caster,
								defender: buff.carrier,
								tags:     nil,
								damage: Damage{
									hp:        -1,
									bullet:    0,
									explosion: 0,
									mental:    0,
								},
								criticalRate: 0,
								hitRate:      0,
								degree:       0,
								addBuffs:     nil,
							}
							scene.damagesInfo = append(scene.damagesInfo, damageInfo)
						},
						onTickParams:     nil,
						onRemoved:        nil,
						onRemovedParams:  nil,
						onCast:           nil,
						onCastParams:     nil,
						onHit:            nil,
						onHitParams:      nil,
						onBeHurt:         nil,
						onBeHurtParams:   nil,
						onKill:           nil,
						onKillParams:     nil,
						onBeKilled:       nil,
						onBeKilledParams: nil,
					}
					addBuffInfo := AddBuffInfo{
						caster:        timeline.caster,
						target:        timeline.caster,
						buffModel:     buffModel,
						addStack:      1,
						durationSetTo: true,
						permanent:     false,
						duration:      5,
					}
					fmt.Printf("释放技能后0.1s执行node2\n")
					cs.AddBuff(addBuffInfo)
				},
				[]object{},
			),
			// 0.1s时候执行生效逻辑
			NewTimelineNode(
				0.1,
				func(timeline *TimelineObj, args []object) {
					if nil == timeline.caster {
						return
					}
					cs := timeline.caster.GetChaState()
					if nil == cs {
						return
					}
					// 添加一个buff吧
					buffModel := BuffModel{
						id:       "加攻击力",
						name:     "加攻击力",
						priority: 1,
						maxStack: 1,
						tags:     nil,
						tickTime: 0,
						propMod:  nil,
						stateMod: ChaControlStateOrigin,
						onOccur: func(buff *BuffObj, modifyStack int) {
							fmt.Printf("加20点攻击力\n")
						},
						onOccurParams: nil,
						onTick:        nil,
						onTickParams:  nil,
						onRemoved: func(buff *BuffObj) {
							fmt.Printf("减少20点攻击力\n")
						},
						onRemovedParams:  nil,
						onCast:           nil,
						onCastParams:     nil,
						onHit:            nil,
						onHitParams:      nil,
						onBeHurt:         nil,
						onBeHurtParams:   nil,
						onKill:           nil,
						onKillParams:     nil,
						onBeKilled:       nil,
						onBeKilledParams: nil,
					}
					addBuffInfo := AddBuffInfo{
						caster:        timeline.caster,
						target:        timeline.caster,
						buffModel:     buffModel,
						addStack:      1,
						durationSetTo: true,
						permanent:     false,
						duration:      1,
					}
					fmt.Printf("释放技能0.1s执行node3\n")
					cs.AddBuff(addBuffInfo)
				},
				[]object{},
			),
			// 0.1s后可使用其他技能
			NewTimelineNode(
				0.1,
				func(timeline *TimelineObj, args []object) {
					if nil == timeline.caster {
						return
					}
					cs := timeline.caster.GetChaState()
					if nil == cs {
						return
					}
					cs.timelineControlState.canMove = args[0].(bool)
					cs.timelineControlState.canRotate = args[0].(bool)
					cs.timelineControlState.canUseSkill = args[0].(bool)
					fmt.Printf("执行node1.3\n")
				},
				[]object{true, true, true},
			),
		},
		0.1,
		TimelineGoTo{},
	)
	id := "每0.1秒扣除1hp持续1s，增加20attack, cd2s"
	skillModel := NewSkillModel(
		id,
		&ChaResource{},
		&ChaResource{},
		timelineMode,
		nil,
	)
	skillObj := NewSkillObj(
		skillModel,
		1,
	)
	skillObj.coolDownStatic = 10

	tb.c.skills = append(tb.c.skills, skillObj)

	tk := time.NewTicker(time.Second)
	const interval float64 = 0.033
	var tick float64
	for {
		select {
		case <-tk.C:
			fmt.Printf("当前[%v]帧，持续时间[%v]\n",
				tick, interval*tick)
			if tb.c.CastSkill(id) {
				fmt.Printf("技能释放成功[%v]\n", id)
			}
			scene.GameObjectTick(interval)
			scene.TimelineTick(interval)
			scene.DamageInfoTick(interval)
			tick++
		}

	}

}
