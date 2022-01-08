package sync_demo

import (
	"fmt"
	"math"
	"net"
)

const (
	dirtyPos = 1 << iota
	dirtyColor
	dirtyPlayerId
	dirtyHealth

	dirtyAll = dirtyPos | dirtyColor | dirtyPlayerId | dirtyHealth
)

type (
	vec3        [3]float32
	ClientProxy struct {
		updAddr *net.UDPAddr
		addr    string
		name    string

		id uint64

		color vec3

		moveList *MoveList

		mIsLastMoveTimestampDirty bool

		location vec3
		velocity vec3

		rotation  float32
		thrustDir float32

		rpl *ReplicateMgr

		dn *DeliveryNotificationManager
	}
)

func (v vec3) Write(b *netBuffer) {
	b.writeFloat(v[0])
	b.writeFloat(v[1])
	b.writeFloat(v[2])
}

func (v vec3) Read(b *netBuffer) vec3 {
	return vec3{
		b.readFloat(),
		b.readFloat(),
		b.readFloat(),
	}
}

func (v vec3) EQ(t vec3) bool {
	return v[0] == t[0] && v[1] == t[1] && v[2] == t[2]
}

func (v vec3) Add(t vec3) vec3 {
	v[0] += t[0]
	v[1] += t[1]
	v[2] += t[2]
	return v
}

func (v vec3) Mul(f float32) vec3 {
	v[0] *= f
	v[1] *= f
	v[2] *= f
	return v
}

func (c *ClientProxy) Update() {
	oldLoc := c.location
	oldVel := c.velocity
	oldRt := c.rotation
	moves := c.moveList.moves
	for it := moves.Front(); nil != it; it = it.Next() {
		unMove := it.Value.(*Move)
		currentState := unMove.inputState
		delta := unMove.deltaTime
		c.ProcessInput(delta, currentState)
		c.SimulateMovement(delta)
	}
	c.moveList.Clear()

	if !oldLoc.EQ(c.location) ||
		!oldVel.EQ(c.velocity) ||
		oldRt != c.rotation {

		for _, v := range cg.entries {
			v.rpl.SetStateDirty(c.id, dirtyPos)
		}
	}
}

func (c *ClientProxy) ProcessInput(dt float32, state InputState) {
	newRt := c.rotation + state.GetHF()*5*dt
	c.rotation = newRt
	c.thrustDir = state.GetVF()
}

func (c *ClientProxy) SimulateMovement(dt float32) {
	c.AdjustVelocityByThrust(dt)
	c.location = c.location.Add(c.velocity.Mul(dt))
}

func (c *ClientProxy) GetForwardVector() vec3 {
	return vec3{
		float32(math.Sin(float64(c.rotation))),
		float32(-math.Cos(float64(c.rotation))),
		0,
	}
}

func (c *ClientProxy) AdjustVelocityByThrust(dt float32) {
	fv := c.GetForwardVector()
	c.velocity = fv.Mul(c.thrustDir * dt * 50)
}

func (c *ClientProxy) Write(b *netBuffer, state uint32) {
	if 0 != state&dirtyPlayerId {
		b.writeBool(true)
		b.writeUint32(uint32(c.id))
	} else {
		b.writeBool(false)
	}

	if 0 != state&dirtyPos {
		b.writeBool(true)
		b.writeFloat(c.velocity[0])
		b.writeFloat(c.velocity[1])

		b.writeFloat(c.location[0])
		b.writeFloat(c.location[1])

		b.writeFloat(c.rotation)
	} else {
		b.writeBool(false)
	}

	if 0 != c.thrustDir {
		b.writeBool(true)
		b.writeBool(c.thrustDir > 0)
	} else {
		b.writeBool(false)
	}

	if 0 != state&dirtyColor {
		b.writeBool(true)
		c.color.Write(b)
	} else {
		b.writeBool(false)
	}

	if 0 != state&dirtyHealth {
		b.writeBool(true)
		b.writeSpBits([]byte{10}, 4)
	} else {
		b.writeBool(false)
	}
}

var (
	cg = &clientGroup{
		entries: map[uint64]*ClientProxy{},
	}

	colours = [4]vec3{
		{1, 1, 0.88},
		{0.68, 0.85, 0.9},
		{0.56, 0.93, 0.56},
		{1, 0.71, 0.76},
	}
)

type clientGroup struct {
	entries map[uint64]*ClientProxy
}

func (c *clientGroup) AddCP(cp *ClientProxy) {
	if nil == c.entries {
		c.entries = map[uint64]*ClientProxy{}
	}
	cp.moveList = NewMoveList()
	c.entries[cp.id] = cp
	cp.color = colours[int(cp.id)%len(colours)]
	cp.location = vec3{3.4, 6.4, 0}
	cp.dn = NewDeliveryNotificationManager(false, true)
	cp.rpl = &ReplicateMgr{}
	for _, v := range c.entries {
		v.rpl.ReplicateCreate(cp.id, dirtyAll)
	}
}

func (c *clientGroup) Update() {
	for _, v := range c.entries {
		v.Update()
	}
}

func (c *clientGroup) SendOutgoingPackets() {
	for _, v := range c.entries {
		v.dn.ProcessTimedOutPackets()
		if v.mIsLastMoveTimestampDirty {
			b := &netBuffer{
				data: make([]byte, 1500),
			}
			b.writeUint32(STAT)

			///////// packet
			ifp := v.dn.WriteState(b)

			///////// move
			b.writeBool(v.mIsLastMoveTimestampDirty)
			if v.mIsLastMoveTimestampDirty {
				b.writeFloat(v.moveList.lastMoveTimestamp)
				v.mIsLastMoveTimestampDirty = false
			}

			//// state
			size := len(c.entries)
			b.writeUint32(uint32(size))
			for _, vv := range c.entries {
				vv.color.Write(b)
				b.writeUint64(vv.id)
				b.writeString(vv.name)
				b.writeUint64(0)
			}

			//// trans
			rmtd := &ReplicationManagerTransmissionData{
				rp: v.rpl,
			}
			v.rpl.Write(b, rmtd)
			ifp.transmissions[RPLM] = rmtd
			_, err := ng.ln.WriteToUDP(b.data[:(b.bitHead+7)>>3], v.updAddr)
			if nil != err {
				fmt.Println(err)
			}
		}
	}
}
