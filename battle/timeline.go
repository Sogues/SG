package battle

// todo

func NewTimelineObj(
	model TimelineModel,

	caster GameObject,
	param object,

) *TimelineObj {
	t := &TimelineObj{
		model:     model,
		caster:    caster,
		timeScale: 1,
		param:     param,
		values:    nil,
	}
	if nil != caster {
		cs := caster.GetChaState()
		if nil != cs {
			t.values = map[string]object{}
			t.values["faceDegree"] = cs.faceDegree
			t.values["moveDegree"] = cs.moveDegree
			t.setTimeScale(cs.actionSpeed)
		}
	}
	return t
}

func NewTimelineMode(
	id string,
	nodes []TimelineNode,

	//Timeline一共多长时间（到时间了就丢掉了），单位秒
	duration float64,

	chargeGoBack TimelineGoTo,
) TimelineModel {
	t := TimelineModel{
		id:           id,
		nodes:        nodes,
		duration:     duration,
		chargeGoBack: chargeGoBack,
	}
	return t
}

func NewTimelineNode(
	timeElapsed float64,

	doEvent TimelineEvent,
	eveParams []object,
) TimelineNode {
	t := TimelineNode{
		timeElapsed: timeElapsed,
		doEvent:     doEvent,
		eveParams:   eveParams,
	}
	return t
}

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

func (t *TimelineObj) setTimeScale(timeScale float64) {
	if timeScale < 0.1 {
		timeScale = 0.1
	}
	t.timeScale = timeScale
}
