package SGTime

import (
	"fmt"
	"testing"
	"time"
)

//BenchmarkTimeUnix
//BenchmarkTimeUnix/1
//BenchmarkTimeUnix/1-12    242596314             4.88 ns/op           0 B/op           0 allocs/op
//BenchmarkTimeUnix/2
//BenchmarkTimeUnix/2-12    246260887             4.90 ns/op           0 B/op           0 allocs/op
//PASS
func BenchmarkTimeUnix(b *testing.B) {
	b.Run("1", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			time.Now().Unix()
		}
		b.ReportAllocs()
	})
	b.Run("2", func(b *testing.B) {
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			time.Now().UnixNano()
		}
		b.ReportAllocs()
	})
}

func TestTimeUnix(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		a := time.Now().UnixNano()
		//for i := 0; i < 1; i++ {
		//	time.Now().UnixNano()
		//}
		b := time.Now()
		fmt.Println(a)
		fmt.Println(b.UnixNano())
		fmt.Println(b.UnixNano() - a)
	})
}

func TestNewFixedUpdate(t *testing.T) {
	t.Run("1", func(t *testing.T) {
		f := NewFixedUpdate(30)
		last := time.Now().UnixNano()
		f.UpdateCallback = func() {
			now := time.Now().UnixNano()
			fmt.Println(now-last, time.Duration(now-last))
			last = now
		}
		for {
			f.Tick()
		}
	})
	t.Run("2", func(t *testing.T) {
		tk := time.NewTicker(time.Millisecond * 30)
		last := time.Now().UnixNano()
		for {
			select {
			case <-tk.C:
				now := time.Now().UnixNano()
				fmt.Println(now-last, time.Duration(now-last))
				last = now
			}
		}
	})
}
