package battle

type (
	ChaState struct {
		controlState ChaControlState

		timelineControlState ChaControlState

		//角色的无敌状态持续时间，如果在无敌状态中，子弹不会碰撞，DamageInfo处理无效化
		immuneTime float64

		charging bool

		moveDegree float64

		faceDegree float64

		dead bool

		moveOrder Vector3

		forceMove []*MovePreorder

		aimOrder []string

		rotateToOrder float64

		forceRotate []float64

		resource *ChaResource

		side int

		tags []string

		property ChaProperty

		moveSpeed float64

		actionSpeed float64

		baseProp ChaProperty

		buffProp []ChaProperty

		equipmentProp ChaProperty

		skills []*SkillObj

		buffs []*BuffObj
	}
)

func (c *ChaState) CanBeKilledByDamageInfo(damageInfo *DamageInfo) bool {
	if c.immuneTime > 0 || damageInfo.IsHeal() {
		return false
	}
	return true
}
