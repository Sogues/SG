package sync_demo

import "container/list"

type (
	InputState struct {
	}
	Move struct {
		inputState InputState
		timestamp  float64
		deltaTime  float64
	}

	MoveList struct {
		lastMoveTimestamp float64
		moves             list.List
	}
)

func (m *MoveList) ClearMove() {
	m.moves.Init()
}

func (m *MoveList) HasMoves() bool {
	return 0 != m.moves.Len()
}

func (m *MoveList) GetMoveCount() int {
	return m.moves.Len()
}
