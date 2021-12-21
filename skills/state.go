package skills

type CharacterControlState int8

// 基础效果
const (
	CharacterControlStateNone CharacterControlState = 1 << iota
	// CharacterControlStateMove 移动
	CharacterControlStateMove
	// CharacterControlStateCastSkill 技能使用
	CharacterControlStateCastSkill

	//// CharacterControlStateSkillCast1 瞬发技能
	//CharacterControlStateSkillCast1
	//// CharacterControlStateSkillCast2 吟唱技能
	//CharacterControlStateSkillCast2
)

type (
	Vec2d [2]uint32
	Vec3d [3]uint32
)
