package battle

// todo

type (
	TimelineEvent func(timeline *TimelineObj, args []object)

	TimelineNode struct {
		timeElapsed float64

		doEvent   TimelineEvent
		eveParams []object
	}

	TimelineGoTo struct {
		//自身处于时间点
		atDuration float64

		//跳转到时间点
		gotoDuration float64
	}

	TimelineModel struct {
		id string

		nodes []TimelineNode

		//Timeline一共多长时间（到时间了就丢掉了），单位秒
		duration float64

		chargeGoBack TimelineGoTo
	}

	TimelineObj struct {
		model TimelineModel

		caster GameObject

		timeScale float64

		//Timeline的创建参数，如果是一个技能，这就是一个skillObj
		param object

		//Timeline已经运行了多少秒了
		timeElapsed float64

		values map[string]object
	}
)
