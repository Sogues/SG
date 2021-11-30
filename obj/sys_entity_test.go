package obj

import "github.com/Sogues/ETForGo/types"

type testRoleEntity struct {
	BaseEntity
}

func (t testRoleEntity) EntityTypeId() types.EntityType {
	return types.EntityTypeTest1
}

type testRoleAwakeSystem struct {
}

func (testRoleAwakeSystem) ComponentTypeId() types.EntityType {
	return testRoleEntity{}.EntityTypeId()
}
