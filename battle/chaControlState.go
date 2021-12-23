package battle

type ChaControlState struct {
	canMove bool

	canRotate bool

	canUseSkill bool
}

var (
	ChaControlStateOrigin = ChaControlState{true, true, true}
)
