package sync_demo

import (
	"container/list"
)

type (
	InputState struct {
		r, l, f, b, s int8
	}
	Move struct {
		inputState InputState
		timestamp  float32
		deltaTime  float32
	}

	MoveList struct {
		lastMoveTimestamp float32
		moves             *list.List
	}
)

func NewMoveList() *MoveList {
	return &MoveList{
		moves: list.New(),
	}
}

func (i InputState) GetH() int8 { return i.r - i.l }
func (i InputState) GetV() int8 { return i.f - i.b }

func (i InputState) GetHF() float32 { return float32(i.r - i.l) }
func (i InputState) GetVF() float32 { return float32(i.f - i.b) }

func (m *MoveList) GetLastMoveTimestamp() float32 { return m.lastMoveTimestamp }

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

func (m *MoveList) AddMove(state InputState, inTimestamp float32) *Move {
	var deltaTime float32
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
		var deltaTime float32
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

func (m *MoveList) RemovedProcessedMoves(inLastMoveProcessedOnServerTimestamp float32) {
	for 0 != m.moves.Len() && m.moves.Front().Value.(*Move).timestamp <= inLastMoveProcessedOnServerTimestamp {
		m.moves.Remove(m.moves.Front())
	}
}
