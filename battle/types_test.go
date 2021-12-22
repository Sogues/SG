package battle

import (
	"fmt"
	"testing"
)

func TestObject(t *testing.T) {
	var a []object

	a = append(a, 10)
	a = append(a, "xx")

	i, ok := a[0].(int)
	if !ok {
		panic("a[0].(int)")
	}
	fmt.Println(i)
	s, ok := a[1].(string)
	if !ok {
		panic("a[1].(string)")
	}
	fmt.Println(s)
}

func TestXXX(t *testing.T) {
	t.Run("x", func(t *testing.T) {
		type st struct {
			a int
		}
		var s []st
		s = append(s, struct{ a int }{a: 10})
		s = append(s, struct{ a int }{a: 11})
		s = append(s, struct{ a int }{a: 12})
		s[1].a = 20
		fmt.Println(s)
	})
}
