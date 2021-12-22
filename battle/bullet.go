package battle

type (
	BulletOnCreate  func(bullet GameObject)
	BulletOnHit     func(bullet GameObject, target GameObject)
	BulletOnRemoved func(bullet GameObject)
	// BulletTween 子弹的轨迹函数，传入一个时间点，返回出一个Vector3，作为这个时间点的速度和方向，这是个相对于正在飞行的方向的一个偏移（*speed的）
	//正在飞行的方向按照z轴，来算，也就是说，当你只需要子弹匀速行动的时候，你可以让这个函数只做一件事情——return Vector3.forward。
	//t 子弹飞行了多久的时间点，单位秒。
	BulletTween func(t float64, bullet GameObject, target GameObject) Vector3

	// BulletTargettingFunction /子弹在发射瞬间，可以捕捉一个GameObject作为目标，并且将这个目标传递给BulletTween，作为移动参数
	///<param name="bullet">是当前的子弹GameObject，不建议公式中用到这个</param>
	///<param name="targets">所有可以被选作目标的对象，这里是GameManager的逻辑决定的传递过来谁，比如这个游戏子弹只能捕捉角色作为对象，那就是只有角色的GameObject，当然如果需要，加入子弹也不麻烦</param>
	///<return>在创建子弹的瞬间，根据这个函数获得一个GameObject作为followingTarget</return>
	BulletTargettingFunction func(bullet GameObject, targets []GameObject)

	BulletModel struct {
		id string

		prefab string

		//子弹的碰撞半径，单位：米。这个游戏里子弹在逻辑世界都是圆形的，当然是这个游戏设定如此，实际策划的需求未必只能是圆形
		radius float64

		//子弹可以碰触的次数，每次碰到合理目标-1，到0的时候子弹就结束了
		hitTimes int

		//子弹碰触同一个目标的延迟，单位：秒，最小值是Time.fixedDeltaTime（每帧发生一次）
		sameTargetDelay float64

		onCreate      BulletOnCreate
		onCreateParam []object

		onHit       BulletOnHit
		onHitParams []object

		onRemoved       BulletOnRemoved
		onRemovedParams []object

		moveType MoveType

		removeOnObstacle bool

		//子弹是否会命中敌人
		hitFoe bool

		//子弹是否会命中盟军
		hitAlly bool
	}

	BulletHitRecord struct {
		target GameObject

		//多久之后还能再次命中，单位秒
		timeToCanHit float64
	}

	BulletLauncher struct {
		module BulletModel

		caster GameObject

		//发射的坐标，y轴是无效的
		firePosition Vector3

		fireDegree float64

		speed float64

		duration float64

		targetFunc BulletTargettingFunction

		tween BulletTween

		///子弹的移动轨迹是否严格遵循发射出来的角度
		///如果是true，则子弹每一帧Tween返回的角度是按照fireDegree来偏移的
		///如果是false，则会根据子弹正在飞的角度(transform.rotation)来算下一帧的角度
		useFireDegreeForever bool

		//子弹创建后多久是没有碰撞的，这样比如子母弹之类的，不会在创建后立即命中目标，但绝大多子弹还应该是0的
		canHitAfterCreated float64

		//子弹的一些特殊逻辑使用的参数，可以在创建子的时候传递给子弹
		param []object
	}
)
