package sync_demo

import (
	"fmt"
	"testing"
	"time"
)

func TestTiming_Update(t *testing.T) {
	tm := NewTiming()
	tk := time.NewTicker(time.Second / 30)
	for {
		select {
		case <-tk.C:
			tm.Update()
			fmt.Println(*tm)
		}
	}
}
