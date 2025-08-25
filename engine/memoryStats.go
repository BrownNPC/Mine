package engine

import (
	"runtime"
	"time"
)


func MemoryGarbageCollectorRunner() {
	for {
		time.Sleep(time.Second)
		// run it every second
		runtime.GC()
	}
}

