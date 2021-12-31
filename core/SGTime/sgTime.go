package SGTime

import (
	"time"
)

const (
	ticksPerMicro int64 = 10
	ticksPerMill  int64 = 1000 * ticksPerMicro
	ticksPerSec   int64 = 1000 * ticksPerMill
	tickPerMin    int64 = 60 * ticksPerSec
	tickPerHour   int64 = 60 * tickPerMin
)

var (
	spTimeZero = spanTime{}
)

type (
	spanTime struct {
		// 每100us为1tick
		ticks int64
	}
	SGTime struct {
		accElpTime          spanTime
		accFrameCountPerSec int64

		totalTime spanTime
		eplTime   spanTime

		frameCount   int64
		timePerFrame spanTime

		isRunningSlowly       bool
		incrementFrameCount   bool
		framePerSecondUpdated bool
		framePerSec           int64
	}
	SGTick struct {
		startNano int64
		lastNano  int64

		elp                spanTime
		elpWithPause       spanTime
		start              spanTime
		totalTime          spanTime
		totalTimeWithPause spanTime

		paused         bool
		pauseStartNano int64
		pausedTime     int64
	}
	FixedUpdate struct {
		updateTime SGTime

		playTick   *SGTick
		updateTick *SGTick
		tick       *SGTick
	}
)

func (t *SGTime) Update(totalTime, elpTime, elpUpdateTime spanTime, isRunningSlowly, incrementFrameCount bool) {
	t.totalTime = totalTime
	t.eplTime = elpTime
	t.isRunningSlowly = isRunningSlowly
	t.framePerSecondUpdated = false
	if incrementFrameCount {
		t.accElpTime = t.accElpTime.Add(elpTime)
		accElpInSec := t.accElpTime.ToSec()
		if t.accFrameCountPerSec > 0 && accElpInSec > 1 {
			t.timePerFrame = spanTime{t.accElpTime.ticks / t.accFrameCountPerSec}
			t.framePerSec = t.accFrameCountPerSec / accElpInSec
			t.accFrameCountPerSec = 0
			t.accElpTime = spTimeZero
			t.framePerSecondUpdated = true
		}
		t.accFrameCountPerSec++
		t.frameCount++
	}
}

func (t *SGTime) Reset(totalTime spanTime) {
	t.Update(totalTime, spTimeZero, spTimeZero, false, false)
}

func (f *FixedUpdate) Tick() {
	f.tick.Tick()
	f.playTick.Tick()
	f.updateTick.Reset()
}

func genSpanTimeFromNano(ns int64) spanTime {
	return spanTime{ns / 10}
}
func (s spanTime) ToNano() int64 {
	return s.ticks * 10
}

func (s spanTime) ToMicro() int64 {
	return s.ticks / ticksPerMicro
}

func (s spanTime) ToMill() int64 {
	return s.ticks / ticksPerMill
}

func (s spanTime) ToSec() int64 {
	return s.ticks / ticksPerSec
}

func (s spanTime) Add(st spanTime) spanTime {
	return spanTime{s.ticks + st.ticks}
}

func (s spanTime) LTZero() bool {
	return s.ticks < 0
}

func (t *SGTick) Reset() {
	t.start = spTimeZero
	t.totalTime = spTimeZero
	t.startNano = time.Now().UnixNano()
	t.lastNano = t.startNano
	t.pausedTime = 0
	t.paused = false
}

func (t *SGTick) Tick() {
	if t.paused {
		t.elp = spTimeZero
		return
	}
	rawNano := time.Now().UnixNano()
	t.totalTime = t.start.Add(genSpanTimeFromNano(rawNano - t.startNano - t.pausedTime))
	t.totalTimeWithPause = t.start.Add(genSpanTimeFromNano(rawNano - t.startNano))
	t.elp = genSpanTimeFromNano(rawNano - t.pausedTime - t.lastNano)
	t.elpWithPause = genSpanTimeFromNano(rawNano - t.lastNano)
	if t.elp.LTZero() {
		t.elp = spTimeZero
	}
	t.lastNano = rawNano
}

func (t *SGTick) Pause() {
	if t.paused {
		return
	}
	t.paused = true
	t.pauseStartNano = time.Now().UnixNano()
}

func (t *SGTick) Resume() {
	if !t.paused {
		return
	}
	t.paused = false
	t.pausedTime += time.Now().UnixNano() - t.pauseStartNano
	t.pauseStartNano = 0
}
