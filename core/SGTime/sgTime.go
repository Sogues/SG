package SGTime

import "time"

const (
	ticksPerMill int64 = 10
	ticksPerSec  int64 = 1000 * ticksPerMill
	tickPerMin   int64 = 60 * ticksPerSec
	tickPerHour  int64 = 60 * tickPerMin
)

var (
	spTimeZero = spanTime{}
)

type (
	spanTime struct {
		ticks int64
	}
	SGTime struct {
		accElpTime spanTime
	}
	SGTick struct {
		nano int64

		elp spanTime

		paused bool
	}
	FixedUpdate struct {
		updateTime SGTime
	}
)

func (s *spanTime) ToMill() int64 {
	return s.ticks / ticksPerMill
}

func (s *spanTime) ToSec() int64 {
	return s.ticks / ticksPerSec
}

func (s *spanTime) Add(st *spanTime) {
	s.ticks += st.ticks
}

func (t *SGTick) Tick() {
	if t.paused {
		t.elp = spTimeZero
		return
	}
	rawNano := time.Now().UnixNano()
	_ = rawNano
}
