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
