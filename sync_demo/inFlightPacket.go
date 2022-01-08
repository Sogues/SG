package sync_demo

import (
	"container/list"
)

type (
	packetSeqNum     uint16
	transmissionData interface {
		HandleDeliveryFailure(manager *DeliveryNotificationManager)
		HandleDeliverySuccess(manager *DeliveryNotificationManager)
	}
	InFlightPacket struct {
		seqNum         packetSeqNum
		timeDispatched float32
		transmissions  map[uint32]transmissionData
	}
	AckRange struct {
		start packetSeqNum
		count uint32
	}
	DeliveryNotificationManager struct {
		mNextOutgoingSequenceNumber packetSeqNum
		mNextExpectedSequenceNumber packetSeqNum

		mInFlightPackets   list.List
		pendingAcks        list.List
		mShouldSendAcks    bool
		mShouldProcessAcks bool

		mDeliveredPacketCount  uint32
		mDroppedPacketCount    uint32
		mDispatchedPacketCount uint32
	}
)

func NewAckRangeWithSeq(seq packetSeqNum) *AckRange {
	return &AckRange{
		start: seq,
		count: 0,
	}
}

func (a *AckRange) getCountSeq() packetSeqNum { return packetSeqNum(a.count) }

func (a *AckRange) ExtendIfShould(seq packetSeqNum) bool {
	if seq == a.start+a.getCountSeq() {
		a.count++
		return true
	}
	return false
}

func (a *AckRange) Write(b *netBuffer) {
	b.writeInteger(a.start)
	hasCount := a.count > 1
	b.writeBool(hasCount)
	if hasCount {
		countMinusOne := a.count - 1
		var countToAck uint8
		if countMinusOne > 255 {
			countToAck = 255
		} else {
			countToAck = uint8(countMinusOne)
		}
		b.writeInteger(countToAck)
	}
}

func (a *AckRange) Read(b *netBuffer) {
	a.start = packetSeqNum(b.readInteger(a.start))
	hasCount := b.readBool()
	if hasCount {
		var countMinusOne uint8
		countMinusOne = uint8(b.readInteger(countMinusOne))
		a.count = uint32(countMinusOne) + 1
	} else {
		a.count = 1
	}
}

func NewDeliveryNotificationManager(
	inShouldSendAcks, inShouldProcessAcks bool,
) *DeliveryNotificationManager {
	return &DeliveryNotificationManager{
		mShouldSendAcks:    inShouldSendAcks,
		mShouldProcessAcks: inShouldProcessAcks,
	}
}

func (d *DeliveryNotificationManager) WriteState(b *netBuffer) *InFlightPacket {
	out := d.WriteSeqNum(b)
	if d.mShouldSendAcks {
		d.WriteAckData(b)
	}
	return out
}

func (d *DeliveryNotificationManager) ReadAndProcessState(b *netBuffer) bool {
	out := d.ProcessSeqNum(b)
	if d.mShouldProcessAcks {
		d.ProcessAcks(b)
	}
	return out
}

func (d *DeliveryNotificationManager) WriteSeqNum(b *netBuffer) *InFlightPacket {
	d.mNextOutgoingSequenceNumber++
	seq := d.mNextOutgoingSequenceNumber
	b.writeInteger(seq)
	d.mDispatchedPacketCount++
	if d.mShouldProcessAcks {
		flight := NewInFlightPacket(seq)
		d.mInFlightPackets.PushBack(flight)
		return flight
	}
	return nil
}

func (d *DeliveryNotificationManager) WriteAckData(b *netBuffer) {
	hasAcks := d.pendingAcks.Len() != 0
	if hasAcks {
		f := d.pendingAcks.Front()
		f.Value.(*AckRange).Write(b)
		d.pendingAcks.Remove(f)
	}
}

func (d *DeliveryNotificationManager) ProcessSeqNum(b *netBuffer) bool {
	var seq packetSeqNum
	seq = packetSeqNum(b.readInteger(seq))

	if seq >= d.mNextExpectedSequenceNumber {
		d.mNextExpectedSequenceNumber = seq + 1
	} else if seq < d.mNextExpectedSequenceNumber {
		return false
	}
	if d.mShouldSendAcks {
		d.AddPendingAck(seq)
	}
	return true
}

func (d *DeliveryNotificationManager) ProcessAcks(b *netBuffer) {
	hasAcks := b.readBool()
	if !hasAcks {
		return
	}
	ackRange := AckRange{}
	ackRange.Read(b)

	nextActSeqNum := ackRange.start
	onePastAckdSeqNum := nextActSeqNum + ackRange.getCountSeq()
	for nextActSeqNum < onePastAckdSeqNum && 0 != d.mInFlightPackets.Len() {
		front := d.mInFlightPackets.Front()
		nextInflightPacket := front.Value.(*InFlightPacket)
		nextInFlightPacketSequenceNumber := nextInflightPacket.seqNum
		if nextInFlightPacketSequenceNumber < nextActSeqNum {
			d.mInFlightPackets.Remove(front)

			d.HandlePacketDeliveryFailure(nextInflightPacket)

		} else if nextInFlightPacketSequenceNumber == nextActSeqNum {
			d.HandlePacketDeliverySuccess(nextInflightPacket)

			d.mInFlightPackets.Remove(front)

			nextActSeqNum++
		} else {
			nextActSeqNum++
		}
	}
}

func (d *DeliveryNotificationManager) ProcessTimedOutPackets() {
	timeout := TimingIst.GetTimeF() - 0.5

	for 0 != d.mInFlightPackets.Len() {
		front := d.mInFlightPackets.Front()
		p := front.Value.(*InFlightPacket)
		if p.timeDispatched < timeout {
			d.HandlePacketDeliveryFailure(p)
			d.mInFlightPackets.Remove(front)
		} else {
			break
		}
	}
}

func (d *DeliveryNotificationManager) HandlePacketDeliveryFailure(ifp *InFlightPacket) {
	d.mDroppedPacketCount++
	ifp.HandleDeliveryFailure(d)
}

func (d *DeliveryNotificationManager) HandlePacketDeliverySuccess(ifp *InFlightPacket) {
	d.mDroppedPacketCount++
	ifp.HandleDeliverySuccess(d)
}

func (d *DeliveryNotificationManager) AddPendingAck(seq packetSeqNum) {
	if 0 == d.pendingAcks.Len() || !d.pendingAcks.Back().Value.(*AckRange).ExtendIfShould(seq) {
		d.pendingAcks.PushBack(NewAckRangeWithSeq(seq))
	}
}

func NewInFlightPacket(seq packetSeqNum) *InFlightPacket {
	return &InFlightPacket{
		seqNum:         seq,
		timeDispatched: TimingIst.GetTimeF(),
		transmissions:  map[uint32]transmissionData{},
	}
}

func (p *InFlightPacket) HandleDeliveryFailure(d *DeliveryNotificationManager) {
	for _, v := range p.transmissions {
		v.HandleDeliveryFailure(d)
	}
}
func (p *InFlightPacket) HandleDeliverySuccess(d *DeliveryNotificationManager) {
	for _, v := range p.transmissions {
		v.HandleDeliverySuccess(d)
	}
}
