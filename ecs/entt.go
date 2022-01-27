package ecs

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
)

func Create() EntityType {
	out := EntityType(generateIdentifier(len(entities)))
	entities = append(entities, out)
	return out
}

func Emplace[component any](fn func(out *component))*component {
	// todo 内存紧凑
	var c component
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
