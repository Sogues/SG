package ecs

import (
	"fmt"
	"reflect"
	"testing"
)

func TestRegistry_Create(t *testing.T) {
	for i := 0; i < 10; i++ {
		fmt.Println(i, Create())
	}
	type sa struct {
		a int
	}
	RegComponent[sa]()
	a := Emplace[sa](Create(), func(out *sa) {
		out.a = 10
	})
	fmt.Println(a)
}

func BenchmarkReflect(b *testing.B) {
	b.Run("1", func(b *testing.B) {
		type sa struct {
			a int
		}
		for i := 0; i < b.N; i++ {
			reflect.TypeOf(sa{}).Name()
		}
		b.ReportAllocs()
	})
	b.Run("2", func(b *testing.B) {
		type sa struct {
			a int
		}
		for i := 0; i < b.N; i++ {
			reflect.TypeOf(sa{}).String()
		}
		b.ReportAllocs()
	})
	b.Run("3", func(b *testing.B) {
		type sa struct {
			a int
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			Emplace[sa](0, func(out *sa) {

			})
		}

		b.ReportAllocs()
	})
}
