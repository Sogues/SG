package obj

import (
	"testing"

	"github.com/Sogues/ETForGo/types"
)

type (
	testRoleEntity struct {
		BaseEntity
	}
	testRoleEntityAwakeParam struct {
	}
)

func (t testRoleEntity) EntityTypeId() types.EntityType {
	return types.EntityTypeTest1
}

type testRoleAwakeSystem struct {
	*AwakeBaseSystem
}

func (testRoleAwakeSystem) ComponentTypeId() types.EntityType {
	return testRoleEntity{}.EntityTypeId()
}

func (testRoleAwakeSystem) Awake(component, param interface{}) {
	cp, ok := component.(*testRoleEntity)
	if !ok {
		return
	}
	p, ok := param.(testRoleEntityAwakeParam)
	if !ok {
		return
	}
	testRoleAwakeSystem{}.awake(cp, p)
}

func (testRoleAwakeSystem) awake(cmp *testRoleEntity, param testRoleEntityAwakeParam) {
}

func TestCombine(t *testing.T) {
	s := &testRoleAwakeSystem{}
	s.impl = s
}

var b1
func BenchmarkTestRoleAwakeSystemAwake(b *testing.B) {

}
