package SGTime

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
		accElpTime spanTime
	}
	SGTick struct {
		startNano int64

		elp   spanTime
		start spanTime
		total spanTime

		paused bool
	}
	FixedUpdate struct {
		updateTime SGTime
	}
)

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

func genSpanTimeFromNano(ns int64) spanTime {
	return spanTime{ns / 10}
}

func (t *SGTick) Tick() {
	//if t.paused {
	//	t.elp = spTimeZero
	//	return
	//}
	//rawNano := time.Now().UnixNano()
	//t.total = t.start.Add(genSpanTimeFromNano(rawNano - t.startNano))
}
