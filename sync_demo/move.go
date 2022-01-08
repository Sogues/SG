package sync_demo

import (
	"container/list"
)

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

func (m *MoveList) GetLastMoveTimestamp() float64 { return m.lastMoveTimestamp }

func (m *MoveList) GetLatestMove() *Move {
	return m.moves.Back().Value.(*Move)
}

func (m *MoveList) Clear() {
	m.moves.Init()
}

func (m *MoveList) HasMoves() bool {
	return 0 != m.moves.Len()
}

func (m *MoveList) GetMoveCount() int {
	return m.moves.Len()
}

func (m *MoveList) AddMove(state InputState, inTimestamp float64) *Move {
	var deltaTime float64
	if m.lastMoveTimestamp > 0 {
		deltaTime = inTimestamp - m.lastMoveTimestamp
	}
	move := &Move{
		inputState: state,
		timestamp:  inTimestamp,
		deltaTime:  deltaTime,
	}
	m.moves.PushBack(move)
	m.lastMoveTimestamp = inTimestamp
	return move
}

func (m *MoveList) AddMoveIfNew(inMove *Move) bool {
	timestamp := inMove.timestamp
	if timestamp > m.lastMoveTimestamp {
		var deltaTime float64
		if m.lastMoveTimestamp > 0 {
			deltaTime = timestamp - m.lastMoveTimestamp
		}
		m.moves.PushBack(&Move{
			inputState: inMove.inputState,
			timestamp:  timestamp,
			deltaTime:  deltaTime,
		})
		m.lastMoveTimestamp = timestamp
		return true
	}
	return false
}

func (m *MoveList) RemovedProcessedMoves(inLastMoveProcessedOnServerTimestamp float64) {
	for 0 != m.moves.Len() && m.moves.Front().Value.(*Move).timestamp <= inLastMoveProcessedOnServerTimestamp {
		m.moves.Remove(m.moves.Front())
	}
}
