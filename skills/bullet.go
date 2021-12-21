package skills

type BulletModel struct {
	Id uint32
	// 范围 毫米为单位
	RadiusMm        uint32
	HitTimes        uint32
	SameTargetDelay uint32

	OnCreate func()

	OnHit func()

	OnRemoved func()

	MoveType uint32

	RemoveOnObstacle bool

	HitFoe bool

	hitAlly bool
}

type BulletLauncher struct {
	Model        BulletModel
	Cast         Player
	FirePosition Vec3d
	FireDegree   uint32
	Speed        uint32
	// 持续时间
	Duration int64

	TargetFunc func()
}
