package engine

import (
	"fmt"
	"runtime"
	"sync/atomic"
	"time"

	rl "github.com/gen2brain/raylib-go/raylib"
)

var MemStats atomic.Value

func MemoryStatsCollector() {
	var mem runtime.MemStats
	for {
		runtime.ReadMemStats(&mem)
		MemStats.Store(mem)
		time.Sleep(time.Second)
		// more than 70 mb
		if mem.HeapAlloc >= 70*1024*1024 {
			runtime.GC()
		}
	}
}

func DrawMemoryStats(x, y, fontSize int32) {
	mem := MemStats.Load()
	if mem == nil {
		return
	}
	memory := mem.(runtime.MemStats)
	heapUsage := fmt.Sprintf("Heap Usage: %.2f MB", float64(memory.Alloc)/1024/1024)
	rl.DrawText(heapUsage, int32(x), int32(y), fontSize, rl.RayWhite)
	rl.DrawText("F6 to force garbage collection", int32(x), int32(y+20), 20, rl.Red)
}
