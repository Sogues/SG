package SGTime

import (
	"fmt"
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

		playTimer                   *SGTick
		updateTimer                 *SGTick
		timer                       *SGTick
		lastUpdateCount             []int64
		updateCountAverageSlowLimit int64

		singleFrameUpdateTime      spanTime
		totalUpdateTime            spanTime
		maximumElapsedTime         spanTime
		accumulatedElapsedGameTime spanTime
		lastFrameElapsedGameTime   spanTime
		nextLastUpdateCountIndex   int64
		drawRunningSlowly          bool
		forceElapsedTimeToZero     bool

		TargetElapsedTime spanTime

		UpdateCallback func()
	}
)

func (f *FixedUpdate) Tick() {
	f.timer.Tick()
	f.playTimer.Tick()
	f.updateTimer.Reset()
	elpAdjustedTime := f.timer.elpWithPause
	if f.forceElapsedTimeToZero {
		elpAdjustedTime = spTimeZero
		f.forceElapsedTimeToZero = false
	}

	if elpAdjustedTime.GreatThan(f.maximumElapsedTime) {
		elpAdjustedTime = f.maximumElapsedTime
	}
	updateCount := int64(1)
	if elpAdjustedTime.ticks > f.TargetElapsedTime.ticks {
		if elpAdjustedTime.ticks-f.TargetElapsedTime.ticks < f.TargetElapsedTime.ticks>>6 {
			elpAdjustedTime = f.TargetElapsedTime
		}
	} else {
		if f.TargetElapsedTime.ticks-elpAdjustedTime.ticks < f.TargetElapsedTime.ticks>>6 {
			elpAdjustedTime = f.TargetElapsedTime
		}
	}
	f.accumulatedElapsedGameTime = f.accumulatedElapsedGameTime.Add(elpAdjustedTime)

	updateCount = f.accumulatedElapsedGameTime.ticks / f.TargetElapsedTime.ticks

	if 0 == updateCount {
		return
	}
	f.lastUpdateCount[f.nextLastUpdateCountIndex] = updateCount

	var updateCountMean int64
	for _, v := range f.lastUpdateCount {
		updateCountMean += v
	}

	updateCountMean = updateCountMean * 100 / int64(len(f.lastUpdateCount))

	f.nextLastUpdateCountIndex = (f.nextLastUpdateCountIndex + 1) % int64(len(f.lastUpdateCount))

	f.drawRunningSlowly = updateCountMean > f.updateCountAverageSlowLimit

	f.accumulatedElapsedGameTime = spanTime{f.accumulatedElapsedGameTime.ticks - (updateCount * f.TargetElapsedTime.ticks)}
	singleFrameEplTime := f.TargetElapsedTime

	for f.lastFrameElapsedGameTime = spTimeZero; updateCount > 0; updateCount-- {
		f.updateTime.Update(
			f.totalUpdateTime,
			singleFrameEplTime,
			f.singleFrameUpdateTime,
			f.drawRunningSlowly,
			true,
		)
		f.UpdateAndProfile(f.updateTime)
		f.lastFrameElapsedGameTime = f.lastFrameElapsedGameTime.Add(singleFrameEplTime)
		f.totalUpdateTime = f.totalUpdateTime.Add(singleFrameEplTime)

	}
	f.updateTimer.Tick()
	f.singleFrameUpdateTime = spTimeZero
}

func (f *FixedUpdate) UpdateAndProfile(sgTime SGTime) {
	defer func() {
		if err := recover(); nil != err {
			fmt.Println(err)
		}
	}()
	f.updateTimer.Reset()
	if nil != f.UpdateCallback {
		f.UpdateCallback()
	}
	f.lastFrameElapsedGameTime = spTimeZero
	f.updateTimer.Tick()
	f.singleFrameUpdateTime = f.singleFrameUpdateTime.Add(f.updateTimer.elp)
	f.lastFrameElapsedGameTime = spTimeZero
}

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

func (s spanTime) GreatThan(r spanTime) bool {
	return s.ticks > r.ticks
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
