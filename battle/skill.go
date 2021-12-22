package battle

// todo

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
	}
)
