package types

type EntityType uint64

const (
	EntityTypeCat1 uint64 = 48 // 保留8位1类型
	EntityTypeCat2 uint64 = 32 // 保留12位2类型
	EntityTypeCat3 uint64 = 24 // 保留12位3类型
)

const (
	EntityTypeTest EntityType = 1 << (EntityTypeCat1 + iota)
	EntityTypeComponent
	EntityTypeSystem

	EntityTypeNone EntityType = 0
)

const (
	EntityTypeTestNone EntityType = EntityTypeTest + iota
	EntityTypeTest1
	EntityTypeTest2
	EntityTypeTest3
	EntityTypeTest4
	EntityTypeTest5
	EntityTypeTest6
	EntityTypeTest7
	EntityTypeTest8
	EntityTypeTest9
	EntityTypeTest10
)

const (
	EntityTypeSystemNone EntityType = EntityTypeSystem | 1<<(EntityTypeCat2+iota)
	EntityTypeSystemAwake
)
