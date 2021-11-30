package obj

import (
	"sync"
	"testing"

	"github.com/Sogues/ETForGo/types"
)

type (
	componentTest struct {
		val uint32
	}
	componentTestAwakeParam struct {
		val uint32
	}
)

func (componentTest) ComponentTypeId() types.EntityType {
	return types.EntityTypeTest1
}

type awakeBaseSystemImplTest struct {
}

func (awakeBaseSystemImplTest) ComponentTypeId() types.EntityType {
	return componentTest{}.ComponentTypeId()
}

func (i *awakeBaseSystemImplTest) Awake(component, param interface{}) {
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

func (awakeBaseSystemImplTest) awake(cmp *componentTest, param componentTestAwakeParam) {
	cmp.val = param.val * param.val
}

func TestAwakeBaseSystem_Run(t *testing.T) {
	sys := NewAwakeSystemWithImpl(&awakeBaseSystemImplTest{})
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
	sys := NewAwakeSystemWithImpl(&awakeBaseSystemImplTest{})
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

func TestReg(t *testing.T) {
	Reg(NewAwakeSystemWithImpl(&awakeBaseSystemImplTest{}))
	cpt := &componentTest{}
	cpm := componentTestAwakeParam{val: 10}
	for _, v := range global.systems[cpt.ComponentTypeId()] {
		v.Run(cpt, cpm)
	}
	t.Log(cpt)
}

//systems map[ComponentType]map[SystemType]BaseSystem
//BenchmarkReg
//BenchmarkReg-12    	21108661	        70.6 ns/op	       4 B/op	       1 allocs/op
var benchmarkRegOnce = &sync.Once{}

func BenchmarkReg(b *testing.B) {
	benchmarkRegOnce.Do(func() {
		Reg(NewAwakeSystemWithImpl(&awakeBaseSystemImplTest{}))
	})
	cpt := &componentTest{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cpm := componentTestAwakeParam{val: 10}
		for _, v := range global.systems[cpt.ComponentTypeId()] {
			v.Run(cpt, cpm)
		}
	}
	b.ReportAllocs()
}
