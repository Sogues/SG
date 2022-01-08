package sync_demo

import (
	"time"
)

var (
	TimingIst = NewTiming()
)

type Timing struct {
	mDeltaTime          float32
	mDeltaTick          uint64
	mLastFrameStartTime float64
	mFrameStartTime     float32
	mPerfCountDuration  float64

	baseTime int64
}

func NewTiming() *Timing {
	t := &Timing{}
	t.baseTime = time.Now().UnixNano()
	t.mLastFrameStartTime = t.GetTime()
	return t
}

func (t *Timing) GetTime() float64 {
	return float64(time.Now().UnixNano()-t.baseTime) / float64(time.Second)
}

func (t *Timing) GetTimeF() float32 {
	return float32(t.GetTime())
}

func (t *Timing) Update() {
	current := t.GetTime()

	t.mDeltaTime = float32(current - t.mLastFrameStartTime)
	t.mLastFrameStartTime = current
	t.mFrameStartTime = float32(current)
}

func (t *Timing) GetDeltaTime() float32 { return t.mDeltaTime }
