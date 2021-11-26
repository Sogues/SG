package eventsystem

import (
	"sync"
	"testing"
)

var (
	AwakeBaseSystemImplTestComponentId = GenComponentTypeId()
)

type (
	componentTest struct {
		val uint32
	}
	componentTestAwakeParam struct {
		val uint32
	}
)

func (componentTest) ComponentTypeId() ComponentType {
	return AwakeBaseSystemImplTestComponentId
}

type AwakeBaseSystemImplTest struct {
}

func (AwakeBaseSystemImplTest) ComponentTypeId() ComponentType {
	return componentTest{}.ComponentTypeId()
}

func (i *AwakeBaseSystemImplTest) Awake(component, param interface{}) {
	cp, ok := component.(*componentTest)
	if !ok {
		return
	}
	p, ok := param.(componentTestAwakeParam)
	if !ok {
		return
	}
	i.awake(cp, p)
}

func (AwakeBaseSystemImplTest) awake(cmp *componentTest, param componentTestAwakeParam) {
	cmp.val = param.val * param.val
}

func TestAwakeBaseSystem_Run(t *testing.T) {
	sys := NewAwakeSystemWithImpl(&AwakeBaseSystemImplTest{})
	cmp := &componentTest{}
	cmpParam := componentTestAwakeParam{val: 10}
	sys.Run(cmp, cmpParam)
	t.Log(cmp)
}

//BenchmarkAwakeBaseSystem_Run
//BenchmarkAwakeBaseSystem_Run/1
//BenchmarkAwakeBaseSystem_Run/1-12  	70658478	        19.4 ns/op	       4 B/op	       1 allocs/op
//BenchmarkAwakeBaseSystem_Run/2
//BenchmarkAwakeBaseSystem_Run/2-12  	59316453	        20.3 ns/op	       4 B/op	       1 allocs/op
//BenchmarkAwakeBaseSystem_Run/3
//BenchmarkAwakeBaseSystem_Run/3-12  	57221059	        25.0 ns/op	       0 B/op	       0 allocs/op
func BenchmarkAwakeBaseSystem_Run(b *testing.B) {
	sys := NewAwakeSystemWithImpl(&AwakeBaseSystemImplTest{})
	b.Run("1", func(b *testing.B) {
		cmp := &componentTest{}
		b.ResetTimer()
		cmpParam := componentTestAwakeParam{val: 10}
		for i := 0; i < b.N; i++ {
			sys.Run(cmp, cmpParam)
		}
		b.ReportAllocs()
	})
	b.Run("2", func(b *testing.B) {
		cmp := &componentTest{}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			cmpParam := componentTestAwakeParam{val: 10}
			sys.Run(cmp, cmpParam)
		}
		b.ReportAllocs()
	})
	b.Run("3", func(b *testing.B) {
		cmp := &componentTest{}

		p := sync.Pool{}
		p.New = func() interface{} {
			return &componentTestAwakeParam{}
		}
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			pp := p.Get().(*componentTestAwakeParam)
			pp.val = 10
			sys.Run(cmp, pp)
			p.Put(pp)
		}
		b.ReportAllocs()
	})
}
