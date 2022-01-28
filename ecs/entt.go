package ecs

import (
	"reflect"
)

const (
	entityMask  EntityType = 0xfffff
	versionMask EntityType = 0xfff
	entityShift int        = 20
	reserved    EntityType = entityMask | (versionMask << entityShift)
	// todo 分析pageSize
	pageSize int = 4096
)

func toInteger(val valueType) EntityType {
	return EntityType(val)
}

func toEntity(val valueType) EntityType {
	return toInteger(val) & entityMask
}

func toVersion(val valueType) versionType {
	return versionType(toInteger(val) >> entityShift)
}

type (
	EntityType  uint32
	versionType uint16

	valueType uint32
)

var (
	entities []EntityType
	cs       = map[string]int{}
)

func Create() EntityType {
	out := EntityType(generateIdentifier(len(entities)))
	entities = append(entities, out)
	return out
}

func Emplace[component any](entity EntityType, fn func(out *component)) *component {
	// todo 内存紧凑
	var c component

	// todo 如何解决直接根据类型获取类型名的问题
	// 言简意赅就是 根据类型 能够获取唯一id
	r := reflect.TypeOf(c).String()
	var (
		idx int
		ok  bool
	)
	if idx, ok = cs[r]; !ok {
		idx = len(cs) + 1
		cs[r] = idx
	}

	fn(&c)
	return &c
}

func generateIdentifier(pos int) valueType {
	val := combine(EntityType(pos), 0)
	return val
}

func combine(lhs, rhs EntityType) valueType {
	mask := versionMask << entityShift
	return valueType(
		(lhs & entityMask) | (rhs & EntityType(mask)))
}

type sparse struct {
	v1 []uint64
}

type denseHashMap struct {
	sparse sparse
}

func (d *denseHashMap)hashToBucket(hash int) int {
	return d.fastMod(hash, d.bucketCount())
}

func (d *denseHashMap) fastMod(val, mod int) int {
	return val & (mod-1)
}

func (d *denseHashMap) bucketCount() int {
	return len(d.sparse.v1)
}
