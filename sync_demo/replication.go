package sync_demo

const (
	RACreate ReplicationAction = iota
	RAUpdate
	RADestroy
	RARpc
)

type (
	ReplicationAction uint8
	ReplicationCmd    struct {
		dirtyState uint32
		action     ReplicationAction
	}
	ReplicationTransmission struct {
		Id     uint64
		action ReplicationAction
		state  uint32
	}
	ReplicationManagerTransmissionData struct {
		rp            *ReplicateMgr
		transmissions []*ReplicationTransmission
	}
	ReplicateMgr struct {
		id2Replication map[uint64]*ReplicationCmd
	}
)

func (r *ReplicationManagerTransmissionData) AddTransmission(id uint64, action ReplicationAction, state uint32) {
	r.transmissions = append(r.transmissions, &ReplicationTransmission{
		Id:     id,
		action: action,
		state:  state,
	})
}

func (r *ReplicationManagerTransmissionData) HandleDeliveryFailure(d *DeliveryNotificationManager) {
	for _, v := range r.transmissions {
		id := v.Id
		switch v.action {
		case RACreate:
			r.HandleCreateDeliveryFailure(id)
		case RAUpdate:
			r.HandleUpdateStateDeliveryFailure(id, v.state, d)
		case RADestroy:
			r.HandleDestroyDeliveryFailure(id)
		}
	}
}

func (r *ReplicationManagerTransmissionData) HandleDeliverySuccess(d *DeliveryNotificationManager) {
	for _, v := range r.transmissions {
		switch v.action {
		case RACreate:
			r.HandleCreateDeliverySuccess(v.Id)
		case RADestroy:
			r.HandleDestroyDeliverySuccess(v.Id)

		}
	}
}

func (r *ReplicationManagerTransmissionData) HandleCreateDeliverySuccess(id uint64) {
	r.rp.HandleCreateAckd(id)
}

func (r *ReplicationManagerTransmissionData) HandleDestroyDeliverySuccess(id uint64) {
	r.rp.RemoveFromReplication(id)
}

func (r *ReplicationManagerTransmissionData) HandleCreateDeliveryFailure(id uint64) {
	cp := cg.entries[id]
	if nil == cp {
		return
	}
	r.rp.ReplicateCreate(id, dirtyAll)
}

func (r *ReplicationManagerTransmissionData) HandleUpdateStateDeliveryFailure(id uint64, state uint32, d *DeliveryNotificationManager) {
	cp := cg.entries[id]
	if nil == cp {
		return
	}
	for elem := d.mInFlightPackets.Front(); nil != elem; elem = elem.Next() {
		rmtdp := elem.Value.(*InFlightPacket).transmissions[RPLM]
		for _, v := range rmtdp.(*ReplicationManagerTransmissionData).transmissions {
			state &= ^v.state
		}
	}
	if 0 != state {
		r.rp.SetStateDirty(id, state)
	}
}

func (r *ReplicationManagerTransmissionData) HandleDestroyDeliveryFailure(id uint64) {
	r.rp.ReplicateDestroy(id)
}

func (r *ReplicateMgr) ReplicateCreate(id uint64, state uint32) {
	if nil == r.id2Replication {
		r.id2Replication = map[uint64]*ReplicationCmd{}
	}
	r.id2Replication[id] = &ReplicationCmd{
		dirtyState: state,
		action:     RACreate,
	}
}

func (r *ReplicateMgr) ReplicateDestroy(id uint64) {
	c := r.id2Replication[id]
	if nil == c {
		return
	}
	c.action = RADestroy
}

func (r *ReplicateMgr) RemoveFromReplication(id uint64) {
	delete(r.id2Replication, id)
}

func (r *ReplicateMgr) SetStateDirty(id uint64, state uint32) {
	c := r.id2Replication[id]
	if nil == c {
		return
	}
	c.dirtyState |= state
}

func (r *ReplicateMgr) HandleCreateAckd(id uint64) {
	c := r.id2Replication[id]
	if nil == c {
		return
	}
	if c.action == RACreate {
		c.action = RAUpdate
	}
}

func (r *ReplicateMgr) Write(b *netBuffer, data *ReplicationManagerTransmissionData) {
	for k, v := range r.id2Replication {
		if v.action != RADestroy && 0 == v.dirtyState {
			continue
		}
		cp := cg.entries[k]
		b.writeUint64(k)
		b.writeSpBits([]byte{byte(v.action)}, 2)
		switch v.action {
		case RACreate:
			b.writeUint32(RCAT)
			cp.Write(b, v.dirtyState)
		case RAUpdate:
			cp.Write(b, v.dirtyState)
		case RADestroy:
		}
		data.AddTransmission(k, v.action, v.dirtyState)
		v.dirtyState = 0
	}
}
