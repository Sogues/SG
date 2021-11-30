package obj

import (
	"fmt"
	"testing"
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

func (e *entityTestHuman) EntityTypeId() uint64 {
	return 1
}
func (e *entityTestFace) EntityTypeId() uint64 {
	return 2
}
func (e *entityTestBody) EntityTypeId() uint64 {
	return 3
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
