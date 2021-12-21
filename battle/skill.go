package battle

// todo

type (
	SkillModel struct {
		id string

		condition ChaResource

		cost ChaResource

		effect TimelineModel

		buff []AddBuffInfo
	}

	SSkillObj struct {
		model SkillModel

		level int

		coolDown float64
	}
	SkillObj *SSkillObj
)
