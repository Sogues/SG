package battle

// todo

func NewSkillObj(
	model SkillModel,

	level int,
) *SkillObj {
	s := &SkillObj{
		model:    model,
		level:    level,
		coolDown: 0,
	}
	return s
}

func NewSkillModel(
	id string,

	condition *ChaResource,

	cost *ChaResource,

	effect TimelineModel,

	buff []AddBuffInfo,
) SkillModel {

	s := SkillModel{
		id:        id,
		condition: condition,
		cost:      cost,
		effect:    effect,
		buff:      buff,
	}
	return s
}

type (
	SkillModel struct {
		id string

		condition *ChaResource

		cost *ChaResource

		effect TimelineModel

		buff []AddBuffInfo
	}

	SkillObj struct {
		model SkillModel

		level int

		coolDown float64

		coolDownStatic float64
	}
)
