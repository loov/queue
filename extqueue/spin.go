package extqueue

import "runtime"

func spin(v *int) {
	*v++
	if *v > 256 {
		runtime.Gosched()
		*v = 0
	}
}
