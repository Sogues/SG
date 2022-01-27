package ecs

import (
	"fmt"
	"testing"
)

func TestRegistry_Create(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, Create())
	}
	type sa struct {
		a int
	}
	a := Emplace[sa](func(out *sa) {
		out.a = 10
	})
	fmt.Println(a)
}
