garcon
======

Package `garcon` (GARbage CONtrol) provides a lean way to control
garbage collector pauses for standalone worker processes.

In case where a worker process communicates with an external
job dispatcher in load-balancing gc-time-sensitive scenarios,
this library suggests the following workflow:

0. At start time there are no material effects on the worker process.

0. When asked to set a memory limit with `garcon.SetMem()` it disables automatic garbage collection and stores the desired memory limit.

0. On regular basis between jobs the worker code queries the `garcon.Status()` function. If the repy is `GcnMemOK`, then the process allocated memory in use is lower than set threshold and the worker proceeds with it's queue with no changes.

0. If the `garcon.Status()` call returned `GcnMemLow`, then the worker code engages with the remote job dispatcher. They negotiate a pause in the worker's queue. Details of the pause negotiation, memory sizes and thresholds are beyond `garcon`'s concern. That includes decisions when too many worker queues are in a paused state and the dispatcher does not allow this worker to pause, or what to do with jobs already in this worker's queue.

0. Once allowed to pause, the worker calls the `garcon.GC()` function to perform actual garbage collection. Before running the collection two logical checks happen - that the memory limit was actually set and that GC normal operation was disabled. Failure of either check indicates the program logical error and should be handled by the worker after receiving the `GcnFailGC` reply.

0. If calling `garcon.GC()` function returned `GcnDoneGC` the worker asks it's dispatcher to resume the job feed, proceeding with regular memory checks using `garcon.Status()`.
