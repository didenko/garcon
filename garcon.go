/*
Package garcon (GARbage CONtrol) provides a lean way to control
garbage collector pauses for standalone worker processes.
*/
package garcon

import "runtime"
import "runtime/debug"

// GcnRpl is a type for replies from garcon's functions
type GcnRpl uint8

// Specific replies as returned by garcon's functions
const (
	GcnMemOK GcnRpl = iota
	GcnMemLow
	GcnDoneGC
	GcnFailGC
)

var (
	gcnMemLimit uint64
	gcnMemStats runtime.MemStats
)

// The SetMem function set the treshold below which the memory
// consumption will be considered normal and not worthy running
// a garbage collection.
func SetMem(memLimit uint64) GcnRpl {
	gcnMemLimit = memLimit
	debug.SetGCPercent(-1)
	return Status()
}

// The Status function checks the current allocated but not
// freed memory against the treshold set by SetMem and returns
// either GcnMemOK or GcnMemLow appropriately.
func Status() GcnRpl {
	runtime.ReadMemStats(&gcnMemStats)

	if gcnMemStats.Alloc < gcnMemLimit {
		return GcnMemOK
	} else {
		return GcnMemLow
	}
}

// The GC function runs the garbage collection after workflow
// checks as described in workflow in the README.md file.
func GC() GcnRpl {
	if gcnMemLimit == 0 {
		return GcnFailGC
	}
	runtime.ReadMemStats(&gcnMemStats)
	if gcnMemStats.EnableGC {
		return GcnFailGC
	}
	runtime.GC()
	return GcnDoneGC
}
