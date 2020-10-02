package utils

import (
	"fmt"
	"runtime"
	"time"
)


func DoEvery(done <-chan struct{}, d time.Duration, f func()) {
	ticker := time.NewTicker(d)	
	for {
		select {
		case <-done:
			ticker.Stop()
			return
		case <-ticker.C:
			f()
		}
	}

}

//prints memory usage for testing purposes
func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
        fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
        fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}