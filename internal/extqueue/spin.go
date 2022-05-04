package extqueue

import (
	"runtime"
	"time"
)

func wait() {
	runtime.Gosched()
}

func spin(v *int) {
	*v++
	if *v > 256 {
		runtime.Gosched()
		*v = 0
	}
}

func backoff(np *int) {
	n := *np
	*np++

	if n < 3 {
		return
	} else if n < 10 {
		runtime.Gosched()
	} else if n < 12 {
		time.Sleep(0) // osyield
	} else {
		time.Sleep(10 * time.Microsecond)
	}
}
