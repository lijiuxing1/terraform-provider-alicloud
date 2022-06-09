// Code generated by the FlatBuffers compiler. DO NOT EDIT.

package flatbuffer

import (
	flatbuffers "github.com/google/flatbuffers/go"
)

type FlatBufferRowInGroup struct {
	_tab flatbuffers.Table
}

func GetRootAsFlatBufferRowInGroup(buf []byte, offset flatbuffers.UOffsetT) *FlatBufferRowInGroup {
	n := flatbuffers.GetUOffsetT(buf[offset:])
	x := &FlatBufferRowInGroup{}
	x.Init(buf, n+offset)
	return x
}

func GetSizePrefixedRootAsFlatBufferRowInGroup(buf []byte, offset flatbuffers.UOffsetT) *FlatBufferRowInGroup {
	n := flatbuffers.GetUOffsetT(buf[offset+flatbuffers.SizeUint32:])
	x := &FlatBufferRowInGroup{}
	x.Init(buf, n+offset+flatbuffers.SizeUint32)
	return x
}

func (rcv *FlatBufferRowInGroup) Init(buf []byte, i flatbuffers.UOffsetT) {
	rcv._tab.Bytes = buf
	rcv._tab.Pos = i
}

func (rcv *FlatBufferRowInGroup) Table() flatbuffers.Table {
	return rcv._tab
}

func (rcv *FlatBufferRowInGroup) DataSource() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(4))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FlatBufferRowInGroup) Tags() []byte {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(6))
	if o != 0 {
		return rcv._tab.ByteVector(o + rcv._tab.Pos)
	}
	return nil
}

func (rcv *FlatBufferRowInGroup) Time() int64 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(8))
	if o != 0 {
		return rcv._tab.GetInt64(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBufferRowInGroup) MutateTime(n int64) bool {
	return rcv._tab.MutateInt64Slot(8, n)
}

func (rcv *FlatBufferRowInGroup) FieldValues(obj *FieldValues) *FieldValues {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(10))
	if o != 0 {
		x := rcv._tab.Indirect(o + rcv._tab.Pos)
		if obj == nil {
			obj = new(FieldValues)
		}
		obj.Init(rcv._tab.Bytes, x)
		return obj
	}
	return nil
}

func (rcv *FlatBufferRowInGroup) MetaCacheUpdateTime() uint32 {
	o := flatbuffers.UOffsetT(rcv._tab.Offset(12))
	if o != 0 {
		return rcv._tab.GetUint32(o + rcv._tab.Pos)
	}
	return 0
}

func (rcv *FlatBufferRowInGroup) MutateMetaCacheUpdateTime(n uint32) bool {
	return rcv._tab.MutateUint32Slot(12, n)
}

func FlatBufferRowInGroupStart(builder *flatbuffers.Builder) {
	builder.StartObject(5)
}
func FlatBufferRowInGroupAddDataSource(builder *flatbuffers.Builder, dataSource flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(0, flatbuffers.UOffsetT(dataSource), 0)
}
func FlatBufferRowInGroupAddTags(builder *flatbuffers.Builder, tags flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(1, flatbuffers.UOffsetT(tags), 0)
}
func FlatBufferRowInGroupAddTime(builder *flatbuffers.Builder, time int64) {
	builder.PrependInt64Slot(2, time, 0)
}
func FlatBufferRowInGroupAddFieldValues(builder *flatbuffers.Builder, fieldValues flatbuffers.UOffsetT) {
	builder.PrependUOffsetTSlot(3, flatbuffers.UOffsetT(fieldValues), 0)
}
func FlatBufferRowInGroupAddMetaCacheUpdateTime(builder *flatbuffers.Builder, metaCacheUpdateTime uint32) {
	builder.PrependUint32Slot(4, metaCacheUpdateTime, 0)
}
func FlatBufferRowInGroupEnd(builder *flatbuffers.Builder) flatbuffers.UOffsetT {
	return builder.EndObject()
}
