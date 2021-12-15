package obj

import (
	"reflect"
	"testing"

	"github.com/Sogues/SG/types"
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

func (entityTestHuman) New() Entity {
	return &entityTestHuman{}
}
func (entityTestFace) New() Entity {
	return &entityTestFace{}
}
func (entityTestBody) New() Entity {
	return &entityTestBody{}
}

func (e *entityTestHuman) EntityTypeId() types.EntityType {
	return types.EntityTypeTest1
}
func (e *entityTestFace) EntityTypeId() types.EntityType {
	return types.EntityTypeTest2
}
func (e *entityTestBody) EntityTypeId() types.EntityType {
	return types.EntityTypeTest3
}

func getPtr(i interface{}) uintptr {
	return reflect.ValueOf(i).Pointer()
}

func TestBaseEntity_SetDomain(t *testing.T) {
	human := &entityTestHuman{}
	human.SetDomain(human, human)

	if getPtr(human) != getPtr(human.GetDomain()) {
		t.Error("human and human domain not same")
	}
	face := &entityTestFace{}
	face.setParent(face, human)

	if getPtr(human) != getPtr(face.GetDomain()) {
		t.Error("human and face domain not same")
	}
	if getPtr(human) != getPtr(face.GetParent()) {
		t.Error("human and face parent not same")
	}
	t.Logf("human %p %+v \n",
		human, human)
	t.Logf("face %p %+v \n",
		face, face)
}

func TestBaseEntity_SetComponentParent(t *testing.T) {
	human := &entityTestHuman{}
	human.SetDomain(human, human)
	face := &entityTestFace{}
	face.setComponentParent(human, face)
	if nil != face.GetParent() {
		t.Error("mistake invoke, face need nil parent")
	}
	face.setComponentParent(face, human)
	if getPtr(human) != getPtr(face.GetParent()) {
		t.Error("human and face parent not same")
	}
	t.Logf("human %p %+v \n",
		human, human)
	t.Logf("face %p %+v \n",
		face, face)
}