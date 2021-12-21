package skills

import "testing"

type tPlayer struct {
	state CharacterControlState
}

func (t tPlayer) DoNothing() {
}

func (t *tPlayer) CharacterControlState() CharacterControlState {
	return t.state
}

func (t tPlayer) GetSkill(id uint32) *Skill {
	return nil
}

func (t tPlayer) GetBuffs() []*Buff {
	return nil
}

func TestCastSkill(t *testing.T) {

}
