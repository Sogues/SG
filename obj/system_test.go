package obj

import (
	"sync"
	"testing"

	"github.com/Sogues/SG/types"
)

type (
	compTest struct {
		BaseEntity
		val uint32
	}
	compTestSys   struct{ BaseAwakeSystem }
	compTestParam struct {
		Val uint32
	}
)

func (compTest) EntityTypeId() types.EntityType {
	return types.EntityTypeTest1
}

func (compTest) New() Entity {
	return &compTest{}
}

// 可由工具生成
func (compTestSys) EntityTypeId() types.EntityType {
	return compTest{}.EntityTypeId()
}

// 可由工具生成
func (s *compTestSys) Awake(component, param interface{}) {
	c, ok := component.(*compTest)
	if !ok {
		return
	}
	p, ok := param.(*compTestParam)
	if !ok {
		return
	}
	s.awake(c, p)
}

func (s *compTestSys) awake(c *compTest, p *compTestParam) {
	c.val = p.Val << 1
}

func TestSystem(t *testing.T) {
	e := &compTest{}
	e.SetDomain(e, e)
	s := &compTestSys{}
	s.Awake(e, compTestParam{Val: 10})
	t.Logf("%p %+v\n", e, e)
}

var once = sync.Once{}

func TestES_Awake(t *testing.T) {
	once.Do(func() {
		RegSystem(&compTestSys{})
	})
	e := &compTest{}
	e.SetDomain(e, e)
	SystemProcessor.Awake(e, compTestParam{Val: 10})
	t.Logf("%p %+v\n", e, e)
}

//BenchmarkAwake
//BenchmarkAwake/1
//BenchmarkAwake/1-12  	34347932	        37.2 ns/op	       4 B/op	       1 allocs/op
//BenchmarkAwake/2
//BenchmarkAwake/2-12  	36378627	        41.8 ns/op	       4 B/op	       1 allocs/op
//PASS
func BenchmarkAwake(b *testing.B) {
	once.Do(func() {
		RegSystem(&compTestSys{})
	})
	b.Run("1", func(b *testing.B) {
		e := &compTest{}
		e.SetDomain(e, e)
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			SystemProcessor.Awake(e, &compTestParam{Val: 10})
		}
		b.ReportAllocs()
	})
	b.Run("2", func(b *testing.B) {
		e := &compTest{}
		e.SetDomain(e, e)
		p := &compTestParam{Val: 10}
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			SystemProcessor.Awake(e, p)
		}
		b.ReportAllocs()
	})
}

func TestBaseEntity_AddToComponent(t *testing.T) {
	once.Do(func() {
		RegSystem(&compTestSys{})
		RegEntity(&compTest{})
	})
}
