package skills

type AoeModel struct {
	Id uint32

	RemoveOnObstacle bool

	TickInterval int64

	OnCreate func()

	OnTick func()

	OnRemoved func()

	OnChaEnter func()

	OnChaLeave func()

	OnBulletEnter func()

	OnBulletLevel func()
}

type AoeLauncher struct {
	Model AoeModel

	Position Vec3d

	Caster Player

	RadiusMm uint32

	Duration uint32

	Degree uint32
}
