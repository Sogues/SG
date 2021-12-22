package battle

type (
	AoeOnCreate         func(aoe GameObject)
	AoeOnTick           func(aoe GameObject)
	AoeOnRemoved        func(aoe GameObject)
	AoeOnCharacterEnter func(aoe GameObject, cha map[GameObject]struct{})
	AoeOnCharacterLeave func(aoe GameObject, cha map[GameObject]struct{})
	AoeOnBulletEnter    func(aoe GameObject, bullet map[GameObject]struct{})
	AoeOnBulletLeave    func(aoe GameObject, bullet map[GameObject]struct{})
	//aoe 要执行的aoeObj
	//t 这个tween在aoe中运行了多久了，单位：秒
	AoeTween func(aoe GameObject, t float64) *AoeMoveInfo

	AoeModel struct {
		id string

		prefab string
		//aoe是否碰撞到阻挡就摧毁了（removed），如果不是，移动就是smooth的，如果移动的话
		removeOnObstacle bool

		tags []string

		//aoe每一跳的时间，单位：秒
		//如果这个时间小于等于0，或者没有onTick，则不会执行aoe的onTick事件
		tickTime float64

		onCreate       AoeOnCreate
		onCreateParams []object

		onTick       AoeOnTick
		onTickParams []object

		onRemoved       AoeOnRemoved
		onRemovedParams []object

		onChaEnter       AoeOnCharacterEnter
		onChaEnterParams []object

		onChaLeave       AoeOnCharacterLeave
		onChaLeaveParams []object

		onBulletEnter       AoeOnBulletEnter
		onBulletEnterParams []object

		onBulletLeave       AoeOnBulletLeave
		onBulletLeaveParams []object
	}

	AoeMoveInfo struct {
		moveType MoveType

		velocity Vector3

		//aoe的角度变成这个值
		rotateToDegree float64
	}

	AoeLauncher struct {
		model AoeModel

		//释放的中心坐标
		position Vector3

		caster GameObject

		//目前这游戏的设计中，aoe只有圆形，所以只有一个半径，也不存在角度一说，如果需要可以扩展
		radius float64

		//aoe存在的时间，单位：秒
		duration float64

		degree float64

		//aoe移动轨迹函数
		tween      AoeTween
		tweenParam []object

		//aoe的传入参数，比如可以吸收次数之类的
		param map[string]object
	}
)
