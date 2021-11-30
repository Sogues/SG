package obj

import (
	"fmt"
	"testing"

	"github.com/Sogues/ETForGo/types"
)

type (
	entityTestHuman struct {
		BaseEntity
	}
	entityTestFace struct {
		BaseEntity
	}
	entityTestBody struct {
		BaseEntity
	}
)

func (e *entityTestHuman) EntityTypeId() types.EntityType {
	return types.EntityTypeTest1
}
func (e *entityTestFace) EntityTypeId() types.EntityType {
	return types.EntityTypeTest2
}
func (e *entityTestBody) EntityTypeId() types.EntityType {
	return types.EntityTypeTest3
}

func TestBaseEntity_AddToChildren(t *testing.T) {
	human := &entityTestHuman{}
	human.SetDomain(human.GetBaseEntity())
	human.impl = human
	face := &entityTestFace{}
	face.impl = face
	face.setComponentParent(human.GetBaseEntity())
	fmt.Println(human)
	fmt.Println(face)
}
